FROM golang:1.8.3-alpine as builder
ENV buildpath=/usr/local/go/src/build/dpl
ARG build=notSet
RUN mkdir -p $buildpath
ADD . $buildpath
WORKDIR $buildpath

# install deps
RUN apk add --update bash make \
    && apk --update add --no-cache
RUN make bin/systemd-analyse-exporter/release

FROM alpine:3.6
COPY --from=builder /usr/local/go/src/build/systemd-analyse-exporter/_release/systemd-analyse-exporter /systemd-analyse-exporter

ENTRYPOINT ["/systemd-analyse-exporter"]