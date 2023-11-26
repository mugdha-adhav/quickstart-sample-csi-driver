# syntax=docker/dockerfile:1
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY . ./
COPY go.mod go.sum ./
RUN go mod download
RUN go build -o /bin/driver ./cmd/main.go
CMD ["/bin/driver"]

# FROM scratch AS run-stage
# COPY --from=build-stage /bin/driver /bin/driver
# CMD ["/bin/driver"]
