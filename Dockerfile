# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY /pkg/*.go ./pkg/
COPY /pkg/api/*.go ./pkg/api/
COPY /pkg/internal/*.go ./pkg/internal/
COPY *.go ./

RUN go build -o /tdiff
RUN go test -timeout 30s ./...


##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /tdiff /tdiff
COPY ./config.docker.json /config.json

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/tdiff"]