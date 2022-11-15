# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /workspaceone-exporter

EXPOSE 9740

CMD [ "/workspaceone-exporter" ]
