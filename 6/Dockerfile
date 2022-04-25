FROM alpine:3.15 as builder

ARG BUILD_DATE
ARG VCS_REF
ARG NODE_VERSION

RUN addgroup -g 1000 node \
    && adduser -u 1000 -G node -s /bin/sh -D node \
    && apk add --no-cache \
    libstdc++ \
    && apk add --no-cache --virtual .build-deps \
    curl \
    binutils \
    && ARCH= && alpineArch="$(apk --print-arch)" \
    && case "${alpineArch##*-}" in \
    amd64|x86_64|i386:x86-64) \
    ARCH='x64' \
    ;; \
    aarch64|arm64) \
    ARCH='arm64' \
    ;; \
    *) ;; \
    esac \
    && apk add --no-cache --virtual .build-deps-full \
    bash \
    ccache \
    g++ \
    gcc \
    gnupg \
    libgcc \
    libc6-compat \
    linux-headers \
    make \
    python2 \
    python3 \
    xz \
    # gpg keys listed at https://github.com/nodejs/node#release-keys
    && for key in \
    # "Beth Griggs <bgriggs@redhat.com>"
    4ED778F539E3634C779C87C6D7062848A1AB005C \
    # "Bryan English <bryan@bryanenglish.com>"
    141F07595B7B3FFE74309A937405533BE57C7D57 \
    # "Colin Ihrig <cjihrig@gmail.com>"
    94AE36675C464D64BAFA68DD7434390BDBE9B9C5 \
    # "Danielle Adams <adamzdanielle@gmail.com>"
    74F12602B6F1C4E913FAA37AD3A89613643B6201 \
    # "James M Snell <jasnell@keybase.io>"
    71DCFD284A79C3B38668286BC97EC7A07EDE3FC1 \
    # "Michaël Zasso <targos@protonmail.com>"
    8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600 \
    # "Myles Borins <myles.borins@gmail.com>"
    C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8 \
    # "Richard Lau <rlau@redhat.com>"
    C82FA3AE1CBEDC6BE46B9360C43CEC45C17AB93C \
    # "Rod Vagg <rod@vagg.org>"
    DD8F2338BAE7501E3DD5AC78C273792F7D83545D \
    # "Ruben Bridgewater <ruben@bridgewater.de>"
    A48C2BEE680E841632CD4E44F07496B3EB3C1762 \
    # "Ruy Adorno <ruyadorno@hotmail.com>"
    108F52B48DB57BB0CC439B2997B01419BD92F80A \
    # "Shelley Vohr <shelley.vohr@gmail.com>"
    B9E2F5981AA6E0CD28160D9FF13993A75599653C \
    ; do \
    gpg --batch --keyserver hkps://keys.openpgp.org --recv-keys "$key" || \
    gpg --batch --keyserver keyserver.ubuntu.com --recv-keys "$key" ; \
    done \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION.tar.xz" \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" \
    && gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc \
    && grep " node-v$NODE_VERSION.tar.xz\$" SHASUMS256.txt | sha256sum -c - \
    && tar -xf "node-v$NODE_VERSION.tar.xz" \
    && cd "node-v$NODE_VERSION" \
    && ./configure \
    && make install -j$(getconf _NPROCESSORS_ONLN) V= DISTTYPE="release" VARIATION="musl" \
    && cd .. \
    && yes | strip /usr/local/bin/node || true \
    && yes | ln -s /usr/local/bin/node /usr/local/bin/nodejs || true \
    && yes | npx pnpm@1 i -g pnpm@1 \
    # smoke tests
    && node --version \
    && pnpm --version \
    && which pnpm \
    && npm --version

FROM alpine:3.15 as packager

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