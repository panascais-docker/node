# `panascais/node`

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/panascais-docker/node/main.yml?branch=master&style=flat-square)](https://github.com/panascais-docker/node/actions?query=workflow%3Amain)
[![Docker Pulls](https://img.shields.io/docker/pulls/panascais/node.svg?style=flat-square)](https://hub.docker.com/r/panascais/node)
[![Docker Stars](https://img.shields.io/docker/stars/panascais/node.svg?style=flat-square)](https://hub.docker.com/r/panascais/node)
[![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node.svg?style=flat-square)](https://microbadger.com/images/panascais/node)
[![License](https://img.shields.io/github/license/panascais-docker/node.svg?style=flat-square)](https://hub.docker.com/r/panascais/node)

| **Tag:** | **Command:**                     | **Node Version:** | **Supported:** | **Labels:**                                                                                            |
| -------- | -------------------------------- | ----------------- | -------------- | ------------------------------------------------------------------------------------------------------ |
| `lts`    | `docker pull panascais/node:lts` | `v22.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/lts?style=flat-square)    |
| `latest` | `docker pull panascais/node`     | `v24.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/latest?style=flat-square) |
| `24`     | `docker pull panascais/node:24`  | `v24.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/24?style=flat-square)     |
| `23`     | `docker pull panascais/node:23`  | `v23.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/23?style=flat-square)     |
| `22`     | `docker pull panascais/node:22`  | `v22.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/22?style=flat-square)     |
| `21`     | `docker pull panascais/node:21`  | `v21.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/21?style=flat-square)     |
| `20`     | `docker pull panascais/node:20`  | `v20.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/20?style=flat-square)     |
| `19`     | `docker pull panascais/node:19`  | `v19.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/19?style=flat-square)     |
| `18`     | `docker pull panascais/node:18`  | `v18.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/18?style=flat-square)     |
| `17`     | `docker pull panascais/node:17`  | `v17.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/17?style=flat-square)     |
| `16`     | `docker pull panascais/node:16`  | `v16.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/16?style=flat-square)     |
| `15`     | `docker pull panascais/node:15`  | `v15.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/15?style=flat-square)     |
| `14`     | `docker pull panascais/node:14`  | `v14.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/14?style=flat-square)     |
| `13`     | `docker pull panascais/node:13`  | `v13.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/13?style=flat-square)     |
| `12`     | `docker pull panascais/node:12`  | `v12.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/12?style=flat-square)     |
| `11`     | `docker pull panascais/node:11`  | `v11.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/11?style=flat-square)     |
| `10`     | `docker pull panascais/node:10`  | `v10.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/10?style=flat-square)     |
| `9`      | `docker pull panascais/node:9`   | `v9.x.x`          | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/9?style=flat-square)      |
| `8`      | `docker pull panascais/node:8`   | `v8.x.x`          | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/8?style=flat-square)      |
| `7`      | `docker pull panascais/node:7`   | `v7.x.x`          | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/7?style=flat-square)      |
| `6`      | `docker pull panascais/node:6`   | `v6.x.x`          | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/6?style=flat-square)      |
| `4`      | `docker pull panascais/node:4`   | `v4.x.x`          | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/4?style=flat-square)      |

## Build

Every numeric directory is a line named after its node major version. `configuration/` pins the node release, the `mirror.gcr.io` image digests and the pnpm major per line, and `aliases.json` maps `latest` and `lts` onto a line.

**Build, load and smoke test node 22 for the local platform:**

```sh
go run ./scripts build 22
```

The image is tagged `ghcr.io/panascais-docker/node/node:22`, `22.x` and `22.x.x` (plus the Docker Hub and Quay names) in the local image store. `./scripts/build-all-parallel.sh [line...]` does the same for every line at once and logs each build to `$TMPDIR/node-builds`.

**Refresh the pins from nodejs.org and `mirror.gcr.io`:**

```sh
go run ./scripts update
```

On GitHub Actions, `plan` compares each line's fingerprint with the annotation on the published image, `build` pushes the stale lines by digest per architecture and `publish` tags them as multi-platform images.

## Docker cache mounts

BuildKit cache mounts reuse downloaded apk packages between builds. Only the **packager** stage (`ALPINE_IMAGE`) uses them:

```dockerfile
RUN --mount=type=cache,id=node-apk-${TARGETARCH},sharing=locked,target=/apk/cache \
    apk update --cache-dir /apk/cache \
    && apk add --cache-dir /apk/cache --cache-predownload \
    ...
```

The **builder** stage runs on upstream `node:X-alpine`, which may ship apk-tools 2.x (Alpine 3.22 or older). Flags like `--cache-predownload` are apk-tools 3.x only, so the builder installs its two build deps with plain `apk add --no-cache binutils curl` instead.

Alpine 3.23 ships **apk-tools 3.x**, where `--update` / `-U` means `--cache-max-age 0` (always refetch). Do not use it with cache mounts on the packager stage.

Cache IDs include `${TARGETARCH}` so amd64 and arm64 packages do not mix during multi-platform builds. `sharing=locked` avoids races when matrix jobs build in parallel.

## pnpm

Node 18+ images ship pnpm 12 via [get.pnpm.io](https://get.pnpm.io) into `/usr/local/bin/pnpm` (Alpine-compatible). All Node 18+ images set `PNPM_HOME` to `/root/.local/share/pnpm` and a `PATH` that prefers `/usr/local/bin` while including `/root/.local/share/pnpm/bin`, so the standalone CLI is not affected by store mounts in downstream builds.

This image does not cache-mount pnpm during its own build — it only downloads the standalone CLI (~few MB), so a store cache would add complexity for little gain.

Application Dockerfiles should cache installs at a **separate** path such as `/pnpm/store`, not `/root/.local/share/pnpm/store`. See [pnpm Docker docs](https://pnpm.io/docker).

All Node 18+ images ship `enableGlobalVirtualStore: false` in `/root/.config/pnpm/config.yaml` so global installs (`pnpm add -g`) work with BuildKit cache mounts on the store path (see downstream `ci-node`).

## Contributors

- Silas Rech [(silas@panascais.net)](mailto:silas@panascais.net)
- Maximilian Schagginger [(max@panascais.net)](mailto:max@panascais.net)

## Contributing

Interested in contributing to **Node**? Contributions are welcome, and are accepted via pull requests. Please [review these guidelines](contributing.md) before submitting any pull requests.

## License

Code licensed under [MIT](license.md), documentation under [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/).
