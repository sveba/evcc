FROM --platform=$BUILDPLATFORM golang:1.18-alpine as builder

RUN apk update && \
    apk add --no-cache ca-certificates tzdata && \
	update-ca-certificates

FROM alpine:3.15

# Create user
RUN adduser -D evcc && \
    mkdir /evcc && \
	chown evcc /evcc && \
	chmod 777 /evcc

# Import from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY evcc /usr/local/bin/evcc

COPY docker/bin/* /bin/

# UI and /api
EXPOSE 7070/tcp
# KEBA charger
EXPOSE 7090/udp
# OCPP charger
EXPOSE 8887/tcp
# SMA Energy Manager
EXPOSE 9522/udp

HEALTHCHECK --interval=60s --start-period=60s --timeout=30s --retries=3 CMD [ "evcc", "health" ]

USER evcc

WORKDIR /evcc

ENTRYPOINT [ "/bin/entrypoint.sh" ]
CMD [ "evcc", "--sqlite", "/evcc/evcc.db" ]
