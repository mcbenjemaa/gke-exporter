# ARG GOLANG_VERSION=1.16
# ARG GOLANG_OPTIONS="CGO_ENABLED=0 GOOS=linux GOARCH=amd64"

# FROM golang:${GOLANG_VERSION} as build

# ARG VERSION=""
# ARG COMMIT=""

# WORKDIR /gcp-exporter

# COPY go.* ./
# COPY main.go .
# COPY collector ./collector

# RUN env ${GOLANG_OPTIONS} \
#     go build \
#     -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
#     -a -installsuffix cgo \
#     -o /go/bin/gcp-exporter \
#     ./main.go

# FROM gcr.io/distroless/base-debian10

# COPY --from=build /go/bin/gke-info-exporter /

# EXPOSE 9505

# ENTRYPOINT ["/gke-info-exporter"]


# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY collector ./collector
COPY pkg ./pkg

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o gke-info-exporter main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/gke-info-exporter .
USER 65532:65532
EXPOSE 9505

ENTRYPOINT ["/gke-info-exporter"]