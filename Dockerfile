FROM golang:1.18-alpine AS builder
ENV GO111MODULE=on
RUN apk add git
RUN apk add --no-cache git make build-base
WORKDIR /strava
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN cd main/ && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server .

FROM alpine:3.8
RUN apk add --no-cache --update ca-certificates tzdata curl
RUN apk add --no-cache git make build-base
RUN cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime
RUN echo "Asia/Ho_Chi_Minh" >  /etc/timezone
COPY --from=builder /strava/main/server  /server
COPY ./configs/realAccessKeys.json  /strava/configs/realAccessKeys.json

WORKDIR /
EXPOSE 18502
CMD ["/server"]