FROM golang:1.25.0 AS build

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o /out/appserver ./cmd/appserver
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o /out/cyops-pgvector-migrate ./cmd/cyops-pgvector-migrate

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

COPY --from=build /out/appserver /appserver
COPY --from=build /out/cyops-pgvector-migrate /usr/local/bin/cyops-pgvector-migrate

USER 65532:65532
EXPOSE 8443
ENTRYPOINT ["/appserver"]
