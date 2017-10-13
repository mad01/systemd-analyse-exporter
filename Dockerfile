FROM golang:1.8.4-jessie as builder
ENV buildpath=/usr/local/go/src/build/systemd-analyse-exporter
ARG build=notSet
RUN mkdir -p $buildpath
ADD . $buildpath
WORKDIR $buildpath

RUN make build/release

FROM debian:8
COPY --from=builder /usr/local/go/src/build/systemd-analyse-exporter/_release/systemd-analyse-exporter /systemd-analyse-exporter

ENTRYPOINT ["/systemd-analyse-exporter"]
