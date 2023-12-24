FROM golang:1.21-alpine AS build-stage

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=linux go build -o /bin/driver ./cmd/main.go

FROM alpine AS run-stage
COPY --from=build-stage /bin/driver /bin/driver
ENTRYPOINT ["/bin/driver"]
