# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.19-alpine AS build

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY *.go ./

RUN go mod tidy

# compile application
RUN go build -o /workspaceone-exporter

##
## STEP 2 - DEPLOY
##
FROM scratch

WORKDIR /

COPY --from=build /workspaceone-exporter /workspaceone-exporter

EXPOSE 9740

ENTRYPOINT ["/workspaceone-exporter"]