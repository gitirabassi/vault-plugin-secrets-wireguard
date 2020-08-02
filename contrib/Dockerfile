FROM rust

WORKDIR /workspace
RUN cargo install prometheus_wireguard_exporter

FROM alpine
COPY --from=0 /workspace/prometheus_wireguard_exporter /prometheus_wireguard_exporter
