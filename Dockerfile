# Build the manager binary
FROM golang:1.19 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager main.go

RUN curl -Lo asynqmon.tar.gz https://github.com/hibiken/asynqmon/releases/download/v0.7.1/asynqmon_v0.7.1_linux_amd64.tar.gz && \
    tar -xvzf asynqmon.tar.gz && \
    chmod +x asynqmon

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY ./assets ./assets
COPY --from=builder /workspace/manager .
COPY --from=builder /workspace/asynqmon .
USER 65532:65532

COPY --from=builder /go/bin/dlv /dlv

ENTRYPOINT ["/manager"]
