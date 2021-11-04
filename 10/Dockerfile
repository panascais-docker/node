FROM alpine:3.14

ARG BUILD_DATE
ARG VCS_REF
ARG NODE_VERSION
ARG CHECKSUM

LABEL org.label-schema.vcs-url="https://github.com/panascais-docker/node.git" \
    org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.name="Node Image" \
    org.label-schema.description="Panascais Node Image" \
    org.label-schema.vendor="Panascais Inc." \
    org.label-schema.schema-version="1.0.0-rc1" \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.title="Node Image" \
    org.opencontainers.image.description="Panascais Node Image" \
    org.opencontainers.image.url=https://nodejs.org \
    org.opencontainers.image.documentation="https://github.com/panascais-docker/node" \
    org.opencontainers.image.vendor="Panascais Inc." \
    org.opencontainers.image.licenses=MIT \
    org.opencontainers.image.source="https://github.com/panascais-docker/node" \
    maintainer="Panascais Open Source <oss@panascais.net>"

RUN addgroup -g 1000 node \
    && adduser -u 1000 -G node -s /bin/sh -D node \
    && apk add --no-cache \
    libstdc++ \
    && apk add --no-cache --virtual .build-deps \
    curl \
    && if [ -n "${CHECKSUM}" ]; then \
    set -eu; \
    curl -fsSLO --compressed "https://unofficial-builds.nodejs.org/download/release/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64-musl.tar.xz"; \
    echo "${CHECKSUM}  node-v${NODE_VERSION}-linux-x64-musl.tar.xz" | sha256sum -c - \
    && tar -xJf "node-v${NODE_VERSION}-linux-x64-musl.tar.xz" -C /usr/local --strip-components=1 --no-same-owner \
    && ln -s /usr/local/bin/node /usr/local/bin/nodejs; \
    else \
    apk add --no-cache --virtual .build-deps \
    bash \
    binutils-gold \
    ccache \
    g++ \
    gcc \
    gnupg \
    libgcc \
    linux-headers \
    make \
    python2 \
    python3 \
    xz \
    # gpg keys listed at https://github.com/nodejs/node#release-keys
    && for key in \
    # "Ruy Adorno <ruyadorno@hotmail.com>"
    108F52B48DB57BB0CC439B2997B01419BD92F80A \
    # "Julien Gilli <jgilli@fastmail.fm>"
    114F43EE0176B71C7BC219DD50A3051F888C628D \
    # "Danielle Adams <adamzdanielle@gmail.com>"
    1C050899334244A8AF75E53792EF661D867B9DFA \
    # "Beth Griggs <bgriggs@redhat.com>"
    4ED778F539E3634C779C87C6D7062848A1AB005C \
    # "Italo A. Casas <me@italoacasas.com>"
    56730D5401028683275BD23C23EFEFE93C4CFFFE \
    # "keybase.io/jasnell <jasnell@keybase.io>"
    71DCFD284A79C3B38668286BC97EC7A07EDE3FC1 \
    # "Danielle Adams <adamzdanielle@gmail.com>"
    74F12602B6F1C4E913FAA37AD3A89613643B6201 \
    # "Gibson Fahnestock <gibfahn@gmail.com>"
    77984A986EBC2AA786BC0F66B01FBB92821C587A \
    # "Timothy J Fontaine (Personal) <tjfontaine@gmail.com>"
    7937DFD2AB06298B2293C3187D33FF9D0246406D \
    # "Michaël Zasso (Targos) <targos@protonmail.com>"
    8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600 \
    # "isaacs (http://blog.izs.me/) <i@izs.me>"
    93C7E9E91B49E432C2F75674B0A78B0A6C481CF6 \
    # "Colin Ihrig <cjihrig@gmail.com>"
    94AE36675C464D64BAFA68DD7434390BDBE9B9C5 \
    # "Christopher Dickinson <christopher.s.dickinson@gmail.com>"
    9554F04D7259F04124DE6B476D5A82AC7E37093B \
    # "Ruben Bridgewater <ruben@bridgewater.de>"
    A48C2BEE680E841632CD4E44F07496B3EB3C1762 \
    # "Evan Lucas <evanlucas@me.com>"
    B9AE9905FFD7803F25714661B63B535A4C206CA9 \
    # "Shelley Vohr (security is major key) <shelley.vohr@gmail.com>"
    B9E2F5981AA6E0CD28160D9FF13993A75599653C \
    # "Myles Borins <myles.borins@gmail.com>"
    C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8 \
    # "Richard Lau <rlau@redhat.com>"
    C82FA3AE1CBEDC6BE46B9360C43CEC45C17AB93C \
    # "Rod Vagg <rod@vagg.org>"
    DD8F2338BAE7501E3DD5AC78C273792F7D83545D \
    # "Jeremiah Senkpiel <fishrock123@rocketmail.com>"
    FD3A5288F042B6850C66B31F09FE44734EB7990E \
    ; do \
    gpg --batch --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys "$key" || \
    gpg --batch --keyserver hkp://ipv4.pool.sks-keyservers.net --recv-keys "$key" || \
    gpg --batch --keyserver hkp://pgp.mit.edu:80 --recv-keys "$key" || \
    gpg --batch --keyserver hkp://keys.openpgp.org --recv-keys "$key"; \
    done \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION.tar.xz" \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" \
    && gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc \
    && grep " node-v$NODE_VERSION.tar.xz\$" SHASUMS256.txt | sha256sum -c - \
    && tar -xf "node-v$NODE_VERSION.tar.xz" \
    && cd "node-v$NODE_VERSION" \
    && ./configure \
    && make -j$(getconf _NPROCESSORS_ONLN) V= DESTCPU="x64" ARCH="x64" DISTTYPE="release" VARIATION="musl" \
    && make install \
    && rm -rf "node-v$NODE_VERSION" \
    && rm -rf  /usr/local/*.md \
    && rm -rf  /usr/local/LICENSE \
    && rm -rf  /usr/local/include/node/ \
    && rm "node-v$NODE_VERSION.tar.xz" SHASUMS256.txt.asc SHASUMS256.txt; \
    fi \
    && npx pnpm@4 i -g pnpm@4 \
    && apk del .build-deps;

COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]

CMD [ "node" ]
