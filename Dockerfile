FROM golang:1.15.6-buster AS builder
COPY ./pkg envii_exporter/pkg
COPY ./cmd envii_exporter/cmd
COPY ./go.mod envii_exporter/go.mod
COPY ./go.sum envii_exporter/go.sum
RUN cd envii_exporter \
    && go mod download \
    && CGO_ENABLED=0 go build -o /go/envii_exporter ./cmd/envii_exporter

# TODO: minimize footprint
FROM ubuntu:bionic
COPY --from=builder /go/envii_exporter /
ENTRYPOINT ["/envii_exporter"]
