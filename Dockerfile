FROM golang:alpine AS builder

WORKDIR /src
COPY . .

RUN apk add --no-cache curl make
RUN go version
RUN make clean
RUN make

FROM smallstep/step-cli:latest

COPY --from=builder /src/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY --from=builder /src/bin/server-acme /usr/local/bin/server-acme
COPY --from=builder /src/bin/server-cert /usr/local/bin/server-cert
COPY --from=builder /src/bin/client /usr/local/bin/client

USER root
RUN apk add --no-cache libcap \
    && setcap CAP_NET_BIND_SERVICE=+eip /usr/local/bin/server-acme \
    && setcap CAP_NET_BIND_SERVICE=+eip /usr/local/bin/server-cert
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.5 \
    && curl -s -L -o /usr/local/bin/grpc-health-probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
    && chmod +x /usr/local/bin/grpc-health-probe
RUN chmod +x /usr/local/bin/entrypoint.sh
USER step

VOLUME ["/home/step"]
STOPSIGNAL SIGTERM

ENV DOMAIN=""
ENV STEP_CA_URL=""
ENV STEP_CA_FINGERPRINT=""

HEALTHCHECK CMD /usr/local/bin/grpc-health-probe -tls \
    -tls-ca-cert /home/step/certs/root_ca.crt \
    -tls-server-name $DOMAIN \
    -addr localhost:443

ENTRYPOINT ["/bin/bash", "/usr/local/bin/entrypoint.sh"]
CMD /usr/local/bin/server-acme \
    --domain $DOMAIN \
    --cacert /home/step/certs/root_ca.crt \
    --directory $STEP_CA_URL/acme/acme/directory