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

**Example for bash and node version 22:**

```sh
docker build \
    --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    --build-arg VCS_REF=`git rev-parse --short HEAD` \
    --build-arg NODE_VERSION=`cat ./configuration/versions.json | jq '."22"' -r` \
    --build-arg CHECKSUM=`cat ./configuration/checksums.json | jq '."22"' -r` \
    -t panascais/node:22 \
    ./22
```

**Example for fish and node version 22:**

```fish
docker build \
    --build-arg BUILD_DATE=(date -u +"%Y-%m-%dT%H:%M:%SZ") \
    --build-arg VCS_REF=(git rev-parse --short HEAD) \
    --build-arg NODE_VERSION=(cat ./configuration/versions.json | jq '."22"' -r) \
    --build-arg CHECKSUM=(cat ./configuration/checksums.json | jq '."22"' -r) \
    -t panascais/node:22 \
    ./22
```

## Contributors

- Silas Rech [(silas@panascais.net)](mailto:silas@panascais.net)
- Maximilian Schagginger [(max@panascais.net)](mailto:max@panascais.net)

## Contributing

Interested in contributing to **Node**? Contributions are welcome, and are accepted via pull requests. Please [review these guidelines](contributing.md) before submitting any pull requests.

## License

Code licensed under [MIT](license.md), documentation under [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/).
