FROM golang:1.25.0 AS build

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o /out/manager ./cmd/manager

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

COPY --from=build /out/manager /manager

USER 65532:65532
ENTRYPOINT ["/manager"]
