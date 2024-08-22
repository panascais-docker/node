// mapping to lookup the latest checksum & version by major version
const distributions: Record<string, string> = {
    12: 'https://nodejs.org/dist/latest-v12.x/',
    14: 'https://nodejs.org/dist/latest-v14.x/',
    15: 'https://nodejs.org/dist/latest-v15.x/',
    16: 'https://nodejs.org/dist/latest-v16.x/',
    17: 'https://nodejs.org/dist/latest-v17.x/',
    18: 'https://nodejs.org/dist/latest-v18.x/',
    19: 'https://nodejs.org/dist/latest-v19.x/',
    20: 'https://nodejs.org/dist/latest-v20.x/',
    21: 'https://nodejs.org/dist/latest-v21.x/',
    22: 'https://nodejs.org/dist/latest-v22.x/',
};

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
        continue
    }

    if (!distributions.lts && distribution.lts) {
        distributions.lts = distributions[major]
        break;
    }
}

// for each major mapping we will figure out the checksum & version
const latestByMajor = await Promise.all(Object.entries(distributions).map(async ([distribution, base]) => {
    const url = `${base}SHASUMS256.txt`
    const pattern = '(?<checksum>[0-9a-f]{64})\\s+node-v(?<major>\\d+)\\.(?<minor>\\d+)\\.(?<patch>\\d+)[^\\s]+'

    const { groups: results } = new RegExp(pattern, 'gm').exec(await fetchText(url)) ?? {}
    if (!results) {
        throw new Error(`Unable to parse content at '${url}'`)
    }

    const { checksum, major, minor, patch } = results
    return [distribution, { checksum, version: [major, minor, patch].join('.') }] as const
}))

// we combine checksum & version results into a record
const { checksums: checksumsAfter, versions: versionsAfter } = latestByMajor.reduce((document, [major, { checksum, version }]) => {
    document.checksums[major] = checksum;
    document.versions[major] = version;
    return document;
}, { checksums: {} as Record<string, string>, versions: {} as Record<string, string> });

const checksumsFile = './configuration/checksums.json';
const versionsFile = './configuration/versions.json';

const [checksumsBefore, versionsBefore] = await Promise.all([Bun.file(checksumsFile).json(), Bun.file(versionsFile).json()])

// compare checksums & versions from before & after, if they've changed print "continue", otherwise print "exit" to instruct actions to proceed or skip
if (Bun.deepEquals(checksumsBefore, checksumsAfter, true) && Bun.deepEquals(versionsBefore, versionsAfter, true)) {
    console.log('exit');
} else {
    await Promise.all([writeJson(checksumsFile, checksumsAfter), writeJson(versionsFile, versionsAfter)])
    console.log('continue');
}
