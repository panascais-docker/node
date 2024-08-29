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
};

// mapping to lookup the compatible pnpm version by major version
const pnpm: Record<string, string> = {
    '12': '6',
    '14': '7',
    '15': '7',
    '16': '8',
    '17': '8',
    '18': '9',
    '19': '9',
    '20': '9',
    '21': '9',
    '22': '9',
}

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

// we extend the lookup table by dynamically detecting what is currently the "latest" & "lts" release
for (const distribution of await fetchJson<{ version: string, lts: boolean }[]>('https://nodejs.org/dist/index.json')) {
    const major = extractMajor(distribution.version)

    if (!distributions.latest && !distribution.lts) {
        distributions.latest = distributions[major]
        pnpm.latest = pnpm[major]
        continue
    }

    if (!distributions.lts && distribution.lts) {
        distributions.lts = distributions[major]
        pnpm.lts = pnpm[major]
        break;
    }
}

const checksumsFile = './configuration/checksums.json';
const pnpmFile = './configuration/pnpm.json';
const tagsFile = './configuration/tags.json';
const versionsFile = './configuration/versions.json';

const [checksumsBefore, pnpmBefore, tagsBefore, versionsBefore] = await Promise.all([Bun.file(checksumsFile).json(), Bun.file(pnpmFile).json(), Bun.file(tagsFile).json(), Bun.file(versionsFile).json()])

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

    const { results } = await fetchJson<{ results: { name: string }[] }>(`https://hub.docker.com/v2/repositories/library/node/tags/?page_size=1&name=${version}-alpine`)
    const [first] = results ?? [];
    const { name: tag } = first ?? { name: tagsBefore[distribution] as string | undefined }
    if (!tag?.length) {
        throw new Error(`Unable to find latest docker tag for '${distribution}'`)
    }

    return [distribution, { checksum, version, tag }] as const
}))

// we combine checksum & version results into a record
const { checksums: checksumsAfter, pnpm: pnpmAfter, tags: tagsAfter, versions: versionsAfter } = latestByMajor.reduce((document, [major, { checksum, version, tag }]) => {
    document.checksums[major] = checksum;
    document.pnpm[major] = pnpm[major];
    document.tags[major] = tag;
    document.versions[major] = version;
    return document;
}, { checksums: {} as Record<string, string>, pnpm: {} as Record<string, string>, tags: {} as Record<string, string>, versions: {} as Record<string, string> });

// compare checksums & versions from before & after, if they've changed print "continue", otherwise print "exit" to instruct actions to proceed or skip
if (Bun.deepEquals(checksumsBefore, checksumsAfter, true) && Bun.deepEquals(pnpmBefore, pnpmAfter, true) && Bun.deepEquals(tagsBefore, tagsAfter, true) && Bun.deepEquals(versionsBefore, versionsAfter, true)) {
    console.log('exit');
} else {
    await Promise.all([writeJson(checksumsFile, checksumsAfter), writeJson(pnpmFile, pnpmAfter), writeJson(tagsFile, tagsAfter), writeJson(versionsFile, versionsAfter)])
    console.log('continue');
}
