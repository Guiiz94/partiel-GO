FROM golang:1.21.1-alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /go/app
COPY . .

RUN go install
RUN go build .

ENTRYPOINT [ "/go/app/partiel-GO" ]