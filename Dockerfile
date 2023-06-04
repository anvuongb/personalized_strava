FROM golang:1.18-alpine AS builder
ENV GO111MODULE=on
RUN apk add git
RUN apk add --no-cache git make build-base
WORKDIR /strava
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
EXPOSE 3000
CMD ["go", "run", "main/main.go"]