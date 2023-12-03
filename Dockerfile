FROM golang:1.21 AS build-stage

WORKDIR /app

COPY . ./
COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/driver ./cmd/main.go

FROM alpine AS run-stage
COPY --from=build-stage /bin/driver /bin/driver
ENTRYPOINT ["/bin/driver"]
