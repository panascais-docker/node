FROM node:22.13.0-alpine3.21 AS builder

RUN apk add --no-cache --update binutils curl
RUN yes | strip /usr/local/bin/node || true \
    && yes | ln -s /usr/local/bin/node /usr/local/bin/nodejs || true \
    && yes | curl -fsSL https://get.pnpm.io/v6.js | node - add --global pnpm@9 \
    # smoke tests
    && node --version \
    && pnpm --version \
    && which node \
    && which pnpm

FROM alpine:3.20 AS packager

ARG BUILD_DATE
ARG NODE_VERSION
ARG VCS_REF

ENV PNPM_HOME=~/.local/share/pnpm
ENV PATH=${PNPM_HOME}:${PATH}

LABEL org.label-schema.vcs-url="https://github.com/panascais-docker/node.git" \
    org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.name="Node Image" \
    org.label-schema.description="Panascais Node Image" \
    org.label-schema.vendor="Panascais ehf." \
    org.label-schema.schema-version="1.0.0" \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.version=$NODE_VERSION \
    org.opencontainers.image.title="Node Image" \
    org.opencontainers.image.description="Panascais Node Image" \
    org.opencontainers.image.url=https://nodejs.org \
    org.opencontainers.image.documentation="https://github.com/panascais-docker/node" \
    org.opencontainers.image.vendor="Panascais ehf." \
    org.opencontainers.image.licenses=MIT \
    org.opencontainers.image.source="https://github.com/panascais-docker/node" \
    maintainer="Panascais Open Source <oss@panascais.net>"

COPY --from=builder /usr/local/bin/node         /usr/local/bin/node
COPY --from=builder /usr/local/bin/pnpm         /usr/local/bin/pnpm
COPY --from=builder /usr/local/lib/node_modules /usr/local/lib/node_modules
COPY --from=builder /usr/local/pnpm-global      /usr/local/pnpm-global
COPY --from=builder /usr/local/share            /usr/local/share

COPY docker-entrypoint.sh /usr/local/bin/

RUN addgroup -g 1000 node \
    && adduser -u 1000 -G node -s /bin/sh -D node \
    && apk add --no-cache --update libc6-compat libstdc++ ca-certificates \
    && update-ca-certificates \
    && ln -s /usr/local/bin/node /usr/local/bin/nodejs \
    && ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm \
    && ln -s /usr/local/lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx \
    # smoke tests
    && node --version \
    && pnpm --version \
    && which node \
    && which pnpm

ENTRYPOINT ["docker-entrypoint.sh"]
CMD [ "node" ]
