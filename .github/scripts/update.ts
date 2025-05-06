// mapping to lookup the latest checksum & version by major version
const distributions: Record<string, string> = {
    '12': 'https://nodejs.org/dist/latest-v12.x/',
    '14': 'https://nodejs.org/dist/latest-v14.x/',
    '15': 'https://nodejs.org/dist/latest-v15.x/',
    '16': 'https://nodejs.org/dist/latest-v16.x/',
    '17': 'https://nodejs.org/dist/latest-v17.x/',
    '18': 'https://nodejs.org/dist/latest-v18.x/',
    '19': 'https://nodejs.org/dist/latest-v19.x/',
    '20': 'https://nodejs.org/dist/latest-v20.x/',
    '21': 'https://nodejs.org/dist/latest-v21.x/',
    '22': 'https://nodejs.org/dist/latest-v22.x/',
    '23': 'https://nodejs.org/dist/latest-v23.x/',
    '24': 'https://nodejs.org/dist/latest-v24.x/',
};

// mapping to lookup the compatible pnpm version by major version
const pnpm: Record<string, string> = {
    '12': '6',
    '14': '7',
    '15': '7',
    '16': '8',
    '17': '8',
    '18': '10',
    '19': '10',
    '20': '10',
    '21': '10',
    '22': '10',
    '23': '10',
    '24': '10',
}

const architectures = ['amd64', 'arm64'];

const badVersionBlacklist = ['22.7.0', '22.8.0'];

// utilities

const fetchJson = async <T>(url: string) => {
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`Failed to fetch ${url}: ${response.statusText}`);
    }

    const json = await response.json();
    return json as T
};

const fetchText = async (url: string) => {
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`Failed to fetch ${url}: ${response.statusText}`);
    }

    const text = await response.text();
    return text;
};

const extractMajor = (version: string) => {
    const [major] = version.slice(1).split('.')
    if (major?.trim().length === 0) {
        throw new Error(`Unable to extract major version from '${version}'`)
    }

    return major
}

const writeJson = async  <T>(file: string, json: T) => {
    return Bun.write(file, `${JSON.stringify(json, undefined, 4)}\n`, { createPath: true })
}

const hasArchitecture = (images: { architecture: string }[], architecture: string) => {
    for (const image of images) {
        if (image.architecture === architecture) {
            return true;
        }
    }
    return false;
}

const hasArchitectures = (images: { architecture: string }[]) => {
    return architectures.every(architecture => hasArchitecture(images, architecture))
}

const checksumsFile = './configuration/checksums.json';
const distributionsFile = './configuration/distributions.json';
const pnpmFile = './configuration/pnpm.json';
const tagsFile = './configuration/tags.json';
const versionsFile = './configuration/versions.json';

const [checksumsBefore, distributionsBefore, pnpmBefore, tagsBefore, versionsBefore] = await Promise.all([Bun.file(checksumsFile).json(), Bun.file(distributionsFile).json(), Bun.file(pnpmFile).json(), Bun.file(tagsFile).json(), Bun.file(versionsFile).json()])

// we extend the lookup table by dynamically detecting what is currently the "latest" & "lts" release
for (const distribution of await fetchJson<{ version: string, lts: boolean }[]>('https://nodejs.org/dist/index.json')) {
    const major = extractMajor(distribution.version)

    if (!distributions.latest && !distribution.lts) {
        const distributionChecked = distributions[major] ?? distributionsBefore.latest;
        if (distributionChecked) {
            distributions.latest = distributionChecked;
        }
        const pnpmChecked = pnpm[major] ?? pnpmBefore.latest;
        if (pnpmChecked) {
            pnpm.latest = pnpmChecked;
        }
        continue
    }

    if (!distributions.lts && distribution.lts) {
        const distributionChecked = distributions[major] ?? distributionsBefore.lts;
        if (distributionChecked) {
            distributions.lts = distributionChecked;
        }
        const pnpmChecked = pnpm[major] ?? pnpmBefore.lts;
        if (pnpmChecked) {
            pnpm.lts = pnpmChecked;
        }
        break;
    }
}

// for each major mapping we will figure out the checksum & version
const latestByMajor = await Promise.all(Object.entries(distributions).map(async ([distribution, base]) => {
    const url = `${base}SHASUMS256.txt`
    const pattern = '^(?<checksum>[0-9a-f]{64})\\s+node-v(?<major>\\d+)\\.(?<minor>\\d+)\\.(?<patch>\\d+)\\.tar\\.xz$'

    const { groups } = new RegExp(pattern, 'gm').exec(await fetchText(url)) ?? {}
    if (!groups) {
        throw new Error(`Unable to parse content at '${url}'`)
    }

    const { checksum, major, minor, patch } = groups
    const version = [major, minor, patch].join('.');

    if (badVersionBlacklist.includes(version)) {
        return [distribution, { checksum: checksumsBefore[distribution], distribution: distributionsBefore[distribution], version: versionsBefore[distribution], tag: tagsBefore[distribution] }] as const
    }

    const dockerHubUrl = `https://hub.docker.com/v2/repositories/library/node/tags/?page_size=1&name=${version}-alpine`;

    const { results = [] } = await fetchJson<{ results: { name: string, images: { architecture: string }[] }[] }>(dockerHubUrl)
    const [first = { name: '', images: [] }] = results

    if (!first) {
        console.log({ results, distribution, base, url, dockerHubUrl })
    }

    if (!hasArchitectures(first.images)) {
        return [distribution, { checksum: checksumsBefore[distribution], distribution: distributionsBefore[distribution], version: versionsBefore[distribution], tag: tagsBefore[distribution] }] as const
    }

    const { name: tag } = first ?? { name: tagsBefore[distribution] as string | undefined }
    if (!tag?.length) {
        throw new Error(`Unable to find latest docker tag for '${distribution}'`)
    }

    return [distribution, { checksum, distribution: distributions[distribution], version, tag }] as const
}))

// we combine checksum & version results into a record
const { checksums: checksumsAfter, distributions: distributionsAfter, pnpm: pnpmAfter, tags: tagsAfter, versions: versionsAfter } = latestByMajor.reduce((document, [major, { checksum, distribution, version, tag }]) => {
    document.checksums[major] = checksum;
    document.distributions[major] = distribution;
    document.pnpm[major] = pnpm[major];
    document.tags[major] = tag;
    document.versions[major] = version;
    return document;
}, { checksums: {} as Record<string, string>, distributions: {} as Record<string, string>, pnpm: {} as Record<string, string>, tags: {} as Record<string, string>, versions: {} as Record<string, string> });

// compare checksums & versions from before & after, if they've changed print "continue", otherwise print "exit" to instruct actions to proceed or skip
if (Bun.deepEquals(checksumsBefore, checksumsAfter, true) && Bun.deepEquals(distributionsBefore, distributionsAfter, true) && Bun.deepEquals(pnpmBefore, pnpmAfter, true) && Bun.deepEquals(tagsBefore, tagsAfter, true) && Bun.deepEquals(versionsBefore, versionsAfter, true)) {
    console.log('exit');
} else {
    await Promise.all([writeJson(checksumsFile, checksumsAfter), writeJson(distributionsFile, distributionsAfter), writeJson(pnpmFile, pnpmAfter), writeJson(tagsFile, tagsAfter), writeJson(versionsFile, versionsAfter)])
    console.log('continue');
}
