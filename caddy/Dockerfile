FROM caddy:2.7.6-builder-alpine AS builder

RUN xcaddy build \
    --with github.com/c1pca/single-cert

FROM caddy:2.7.6-alpine

RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    rm -rf /var/cache/apk/* && \
    apk add --no-cache curl

COPY --from=builder /usr/bin/caddy /usr/bin/caddy

COPY ./config.json /etc/caddy/config.json

COPY ./certs/fullchain.pem /secrets/fullchain.pem

CMD caddy run --config=/etc/caddy/config.json