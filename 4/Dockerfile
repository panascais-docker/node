FROM alpine:3.12

ARG BUILD_DATE
ARG VCS_REF

ENV NODE_VERSION 4.9.1

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
    binutils-gold \
    curl \
    g++ \
    gcc \
    gnupg \
    libgcc \
    linux-headers \
    make \
    python2 \
    # gpg keys listed at https://github.com/nodejs/node#release-keys
    && for key in \
    108F52B48DB57BB0CC439B2997B01419BD92F80A \
    114F43EE0176B71C7BC219DD50A3051F888C628D \
    1C050899334244A8AF75E53792EF661D867B9DFA \
    4ED778F539E3634C779C87C6D7062848A1AB005C \
    56730D5401028683275BD23C23EFEFE93C4CFFFE \
    71DCFD284A79C3B38668286BC97EC7A07EDE3FC1 \
    77984A986EBC2AA786BC0F66B01FBB92821C587A \
    7937DFD2AB06298B2293C3187D33FF9D0246406D \
    8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600 \
    93C7E9E91B49E432C2F75674B0A78B0A6C481CF6 \
    94AE36675C464D64BAFA68DD7434390BDBE9B9C5 \
    9554F04D7259F04124DE6B476D5A82AC7E37093B \
    A48C2BEE680E841632CD4E44F07496B3EB3C1762 \
    B9AE9905FFD7803F25714661B63B535A4C206CA9 \
    B9E2F5981AA6E0CD28160D9FF13993A75599653C \
    C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8 \
    C82FA3AE1CBEDC6BE46B9360C43CEC45C17AB93C \
    DD8F2338BAE7501E3DD5AC78C273792F7D83545D \
    FD3A5288F042B6850C66B31F09FE44734EB7990E \
    ; do \
    gpg --batch --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys "$key" || \
    gpg --batch --keyserver hkp://ipv4.pool.sks-keyservers.net --recv-keys "$key" || \
    gpg --batch --keyserver hkp://pgp.mit.edu:80 --recv-keys "$key" ; \
    done \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION.tar.xz" \
    && curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" \
    && gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc \
    && grep " node-v$NODE_VERSION.tar.xz\$" SHASUMS256.txt | sha256sum -c - \
    && tar -xf "node-v$NODE_VERSION.tar.xz" \
    && cd "node-v$NODE_VERSION" \
    && ./configure \
    && make -j$(getconf _NPROCESSORS_ONLN) V= \
    && make install \
    && curl -L https://unpkg.com/@pnpm/self-installer | node \
    && apk del .build-deps \
    && cd .. \
    && rm -Rf "node-v$NODE_VERSION" \
    && rm "node-v$NODE_VERSION.tar.xz" SHASUMS256.txt.asc SHASUMS256.txt

COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]

CMD [ "node" ]
