FROM alpine:3.18 as builder

ARG BUILD_DATE
ARG VCS_REF

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
    python3 \
    xz \
    # gpg keys listed at https://github.com/nodejs/node#release-keys
    && for key in \
    108F52B48DB57BB0CC439B2997B01419BD92F80A \
    114F43EE0176B71C7BC219DD50A3051F888C628D \
    141F07595B7B3FFE74309A937405533BE57C7D57 \
    1C050899334244A8AF75E53792EF661D867B9DFA \
    4ED778F539E3634C779C87C6D7062848A1AB005C \
    56730D5401028683275BD23C23EFEFE93C4CFFFE \
    71DCFD284A79C3B38668286BC97EC7A07EDE3FC1 \
    74F12602B6F1C4E913FAA37AD3A89613643B6201 \
    77984A986EBC2AA786BC0F66B01FBB92821C587A \
    7937DFD2AB06298B2293C3187D33FF9D0246406D \
    890C08DB8579162FEE0DF9DB8BEAB4DFCF555EF4 \
    8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600 \
    93C7E9E91B49E432C2F75674B0A78B0A6C481CF6 \
    94AE36675C464D64BAFA68DD7434390BDBE9B9C5 \
    9554F04D7259F04124DE6B476D5A82AC7E37093B \
    A363A499291CBBC940DD62E41F10027AF002F8B0 \
    A48C2BEE680E841632CD4E44F07496B3EB3C1762 \
    B9AE9905FFD7803F25714661B63B535A4C206CA9 \
    B9E2F5981AA6E0CD28160D9FF13993A75599653C \
    C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8 \
    C82FA3AE1CBEDC6BE46B9360C43CEC45C17AB93C \
    CC68F5A3106FF448322E48ED27F5E38D5B0A215F \
    DD792F5973C6DE52C432CBDAC77ABFA00DDBF2B7 \
    DD8F2338BAE7501E3DD5AC78C273792F7D83545D \
    FD3A5288F042B6850C66B31F09FE44734EB7990E \
    ; do \
    gpg --batch --keyserver hkps://keys.openpgp.org --recv-keys "$key" || \
    gpg --batch --keyserver keyserver.ubuntu.com --recv-keys "$key" ; \
    done \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v22.0.0/node-v22.0.0.tar.xz" \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v22.0.0/SHASUMS256.txt.asc" \
    && gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc \
    && grep " node-v22.0.0.tar.xz\$" SHASUMS256.txt | sha256sum -c - \
    && tar -xf "node-v22.0.0.tar.xz" \
    && cd "node-v22.0.0" \
    && mkdir -p /usr/local/sbin \
    && ln -s /usr/bin/python3 /usr/local/sbin/python3.9 \
    && ./configure \
    && make install -j$(getconf _NPROCESSORS_ONLN) V= DISTTYPE="release" VARIATION="musl" \
    && cd .. \
    && yes | strip /usr/local/bin/node || true \
    && yes | ln -s /usr/local/bin/node /usr/local/bin/nodejs || true \
    && yes | curl -fsSL https://get.pnpm.io/v6.js | node - add --global pnpm \
    # smoke tests
    && node --version \
    && pnpm --version \
    && which pnpm \
    && npm --version

FROM alpine:3.18 as packager

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
