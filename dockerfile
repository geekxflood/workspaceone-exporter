# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /ws1-exporter

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /ws1-exporter /ws1-exporter

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/ws1-exporter"]