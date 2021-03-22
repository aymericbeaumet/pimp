FROM alpine:3.13.2
COPY pimp_*.apk /tmp/
RUN apk add --allow-untrusted /tmp/pimp_*.apk
ENTRYPOINT ["/usr/local/bin/pimp"]
