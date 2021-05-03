[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/panascais-docker/node/main?style=flat-square)](https://github.com/panascais-docker/node/actions?query=workflow%3Amain)
[![Docker Pulls](https://img.shields.io/docker/pulls/panascais/node.svg?style=flat-square)](https://hub.docker.com/r/panascais/node)
[![Docker Stars](https://img.shields.io/docker/stars/panascais/node.svg?style=flat-square)](https://hub.docker.com/r/panascais/node)
[![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node.svg?style=flat-square)](https://microbadger.com/images/panascais/node)
[![License](https://img.shields.io/github/license/panascais-docker/node.svg?style=flat-square)](https://hub.docker.com/r/panascais/node)

| **Tag:** | **Command:**                       | **Node Version:** | **Supported:** | **Labels:**                                                                                                                                                                                                                                                                                                                |
|----------|------------------------------------|-------------------|----------------|----------------------------------------------------------------------------------------|
| `lts`    | `docker pull panascais/node:lts`   | `v14.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/lts)      |
| `latest` | `docker pull panascais/node`       | `v16.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/latest)   |
| `16`     | `docker pull panascais/node:16`    | `v16.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/16)       |
| `15`     | `docker pull panascais/node:15`    | `v15.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/15)       |
| `14`     | `docker pull panascais/node:14`    | `v14.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/14)       |
| `13`     | `docker pull panascais/node:13`    | `v13.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/13)       |
| `12`     | `docker pull panascais/node:12`    | `v12.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/12)       |
| `11`     | `docker pull panascais/node:11`    | `v11.x.x`         | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/11)       |
| `10`     | `docker pull panascais/node:10`    | `v10.x.x`         | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/10)       |
| `9`      | `docker pull panascais/node:9`     | `v9.x.x`          | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/9)        |
| `8`      | `docker pull panascais/node:8`     | `v8.x.x`          | ✓              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/8)        |
| `7`      | `docker pull panascais/node:7`     | `v7.x.x`          | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/7)        |
| `6`      | `docker pull panascais/node:6`     | `v6.x.x`          | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/6)        |
| `4`      | `docker pull panascais/node:4`     | `v4.x.x`          | ×              | ![Docker Image Size](https://img.shields.io/docker/image-size/panascais/node/4)        |

## Build

**Example for bash and node version 14:**
```sh
docker build \
    --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    --build-arg VCS_REF=`git rev-parse --short HEAD` \
    --build-arg NODE_VERSION=`cat ../configuration.json | jq '.versions."14"' -r` \
    --build-arg CHECKSUM=`cat ./configuration.json | jq '.checksums."14"' -r` \
    -t panascais/node:14 \
    ./14
```

**Example for fish and node version 14:**
```fish
docker build \
    --progress=plain \
    --build-arg BUILD_DATE=(date -u +"%Y-%m-%dT%H:%M:%SZ") \
    --build-arg VCS_REF=(git rev-parse --short HEAD) \
    --build-arg NODE_VERSION=(cat ./configuration.json | jq '.versions."14"' -r) \
    --build-arg CHECKSUM=(cat ./configuration.json | jq '.checksums."14"' -r) \
    -t panascais/node:14 \
    ./14
```

## Contributors

 - Silas Rech [(silas@panascais.net)](mailto:silas@panascais.net)

## Contributing:

Interested in contributing to **Node**? Contributions are welcome, and are accepted via pull requests. Please [review these guidelines](contributing.md) before submitting any pull requests.

## License:
Code licensed under [MIT](license.md), documentation under [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/).
