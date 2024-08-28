import { $ } from 'bun'

// constants

const pnpmFile = './configuration/pnpm.json';
const tagsFile = './configuration/tags.json';
const [, , distribution] = Bun.argv;

// utilities

const getPnpmVersion = async () => {
    const json = await Bun.file(pnpmFile).json();

    const version = json[distribution] as string;
    if (typeof version !== 'string' || !version?.length) {
        throw new Error(`Invalid pnpm version found: '${version}' in '${JSON.stringify(json)}\n'`);
    }

    return version;
}

const getImageTags = async () => {
    const json = await Bun.file(tagsFile).json();

    const node = json[distribution] as string;
    if (typeof node !== 'string' || !node?.length || !node.includes('alpine')) {
        throw new Error(`Invalid node tag found: '${node}' in '${JSON.stringify(json)}\n'`);
    }

    const alpine = node.split('-').pop()?.split('alpine')?.pop();
    if (typeof alpine !== 'string' || !alpine?.length || !alpine.includes('.')) {
        throw new Error(`Invalid alpine tag found: '${alpine}' in '${JSON.stringify(json)}\n'`);
    }

    return { node, alpine };
}

const getDockerfilePath = () => {
    return `./${distribution}/Dockerfile`
}

const getDockerfile = async () => {
    const dockerfile = await Bun.file(getDockerfilePath()).text();
    return dockerfile;
}

const setDockerfile = async (dockerfile: string) => {
    return Bun.write(getDockerfilePath(), dockerfile);
}

// business logic

const [pnpm, tag, dockerfile] = await Promise.all([getPnpmVersion(), getImageTags(), getDockerfile()]);
const lines = new Array<string>()

let fromCounter = 0
for (const line of dockerfile.split('\n')) {
    if (line.startsWith("FROM")) {
        switch (fromCounter++) {
            case 0: {
                lines.push(`FROM node:${tag.node} AS builder`);
                break;
            }
            case 1: {
                lines.push(`FROM alpine:${tag.alpine} AS packager`);
                break;
            }
            default: {
                throw new Error(`Dockerfile contains more than one FROM statement`)
            }
        }
        continue;
    }
    if (line.includes('https://get.pnpm.io/v6.js')) {
        lines.push(`    && yes | curl -fsSL https://get.pnpm.io/v6.js | node - add --global pnpm@${pnpm} \\`);
        continue;
    }
    lines.push(line);
}

await setDockerfile(lines.join('\n'));

const [patch] = tag.node.split('-');
const [major, min] = patch.split('.')
const minor = `${major}.${min}`

if (Bun.env.GITHUB_ACTIONS === 'true') {
    const {
        CONTAINER_REGISTRY_TOKEN,
        CONTAINER_REGISTRY_USERNAME,
        DOCKER_REGISTRY_TOKEN,
        DOCKER_REGISTRY_USERNAME,
        QUAY_REGISTRY_TOKEN,
        QUAY_REGISTRY_USERNAME,
    } = Bun.env;

    // hack to pull first without authentication to avoid even more rate limits
    await $`docker buildx create --use && docker buildx build \
        --build-arg BUILD_DATE=${await $`date -u +"%Y-%m-%dT%H:%M:%SZ"`.text()} \
        --build-arg NODE_VERSION=${patch} \
        --build-arg VCS_REF=${await $`git rev-parse --short HEAD`.text()} \
        --platform linux/amd64,linux/arm64 \
        --progress=plain \
        ./${distribution}`;

    await $`echo ${CONTAINER_REGISTRY_TOKEN} | docker login ghcr.io -u ${CONTAINER_REGISTRY_USERNAME} --password-stdin`;
    await $`echo ${DOCKER_REGISTRY_TOKEN} | docker login -u ${DOCKER_REGISTRY_USERNAME} --password-stdin`;
    await $`echo ${QUAY_REGISTRY_TOKEN} | docker login quay.io -u ${QUAY_REGISTRY_USERNAME} --password-stdin`;

    await $`docker buildx build \
        --build-arg BUILD_DATE=${await $`date -u +"%Y-%m-%dT%H:%M:%SZ"`.text()} \
        --build-arg NODE_VERSION=${patch} \
        --build-arg VCS_REF=${await $`git rev-parse --short HEAD`.text()} \
        --push \
        --platform linux/amd64,linux/arm64 \
        --progress=plain \
        -t ghcr.io/panascais-docker/node/node:${distribution} \
        -t panascais/node:${distribution} \
        -t quay.io/panascais/node:${distribution} \
        ./${distribution}`;
} else {
    await $`docker buildx build \
        --build-arg BUILD_DATE=${await $`date -u +"%Y-%m-%dT%H:%M:%SZ"`.text()} \
        --build-arg NODE_VERSION=${patch} \
        --build-arg VCS_REF=${await $`git rev-parse --short HEAD`.text()} \
        --load \
        --progress=plain \
        -t ghcr.io/panascais-docker/node/node:${distribution} \
        -t ghcr.io/panascais-docker/node/node:${major} \
        -t ghcr.io/panascais-docker/node/node:${minor} \
        -t ghcr.io/panascais-docker/node/node:${patch} \
        -t panascais/node:${distribution} \
        -t panascais/node:${major} \
        -t panascais/node:${minor} \
        -t panascais/node:${patch} \
        -t quay.io/panascais/node:${distribution} \
        -t quay.io/panascais/node:${major} \
        -t quay.io/panascais/node:${minor} \
        -t quay.io/panascais/node:${patch} \
        ./${distribution}`;
}
