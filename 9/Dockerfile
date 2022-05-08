FROM node:9-alpine as builder

RUN apk add --no-cache \
    autoconf \
    bash \
    binutils \
    ccache \
    curl \
    g++ \
    gcc \
    gdb \
    gnupg \
    libc6-compat \
    libgcc \
    libstdc++ \
    linux-headers \
    make \
    python2 \
    python3 \
    xz
RUN strip /usr/local/bin/node \
    && npx pnpm@4 i -g pnpm@4 \
    # smoke tests
    && node -v \
    && pnpm -v

FROM alpine:3.15 as packager

ARG BUILD_DATE
ARG VCS_REF

LABEL org.label-schema.vcs-url="https://github.com/panascais-docker/node.git" \
    org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.name="Node Image" \
    org.label-schema.description="Panascais Node Image" \
    org.label-schema.vendor="Panascais ehf." \
    org.label-schema.schema-version="1.0.0" \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.title="Node Image" \
    org.opencontainers.image.description="Panascais Node Image" \
    org.opencontainers.image.url=https://nodejs.org \
    org.opencontainers.image.documentation="https://github.com/panascais-docker/node" \
    org.opencontainers.image.vendor="Panascais ehf." \
    org.opencontainers.image.licenses=MIT \
    org.opencontainers.image.source="https://github.com/panascais-docker/node" \
    maintainer="Panascais Open Source <oss@panascais.net>"

RUN addgroup -g 1000 node \
    && adduser -u 1000 -G node -s /bin/sh -D node \
    && apk add --no-cache \
    libc6-compat \
    libstdc++;

COPY --from=builder /usr/local/bin/node /usr/local/bin/node
COPY --from=builder /usr/local/bin/pnpm /usr/local/bin/pnpm
COPY --from=builder /usr/local/lib/node_modules /usr/local/lib/node_modules
COPY --from=builder /usr/local/share /usr/local/share
COPY --from=builder /usr/local/pnpm-global /usr/local/pnpm-global

COPY docker-entrypoint.sh /usr/local/bin/

RUN ln -s /usr/local/bin/node /usr/local/bin/nodejs \
    && ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm \
    && ln -s /usr/local/lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx;

ENTRYPOINT ["docker-entrypoint.sh"]

CMD [ "node" ]