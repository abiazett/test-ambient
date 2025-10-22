FROM golang:1.21 as builder

WORKDIR /workspace

# Copy Go module files
COPY go.mod go.mod
COPY go.sum go.sum

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o odh .

# Use distroless as minimal base image
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/odh /odh
USER 65532:65532

ENTRYPOINT ["/odh"]