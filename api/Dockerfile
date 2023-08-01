# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /api

FROM --platform=linux/amd64 gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /api /api

EXPOSE 9000

ENTRYPOINT ["/api"]
