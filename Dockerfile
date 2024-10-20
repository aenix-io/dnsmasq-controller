# Build the dnsmasq-controller binary
FROM golang:1.22.7 AS builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o dnsmasq-controller

# Load distroless base to get initial passwd/group files
FROM gcr.io/distroless/static-debian12:latest AS app

# Install dnsmasq 
FROM alpine:3.19.1 AS dnsmasq
# Use distroless passwd/group
COPY --from=app /etc/passwd /etc/passwd
COPY --from=app /etc/group /etc/group
RUN apk add --no-cache dnsmasq

# Copy dnsmasq and altered passwd/group to distroless image
FROM app

COPY --from=dnsmasq /usr/sbin/dnsmasq /usr/sbin/dnsmasq
COPY --from=dnsmasq /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=dnsmasq /lib/libc.musl-x86_64.so.1 /lib/
COPY --from=dnsmasq /etc/passwd /etc/passwd
COPY --from=dnsmasq /etc/group /etc/group
COPY --from=builder /workspace/dnsmasq-controller /dnsmasq-controller

ENTRYPOINT ["/dnsmasq-controller"]
