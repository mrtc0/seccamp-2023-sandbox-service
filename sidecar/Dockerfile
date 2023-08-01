# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /sidecar-proxy

FROM --platform=linux/amd64 gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /sidecar-proxy /sidecar-proxy

EXPOSE 8080

ENTRYPOINT ["/sidecar-proxy"]
