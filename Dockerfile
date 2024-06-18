# syntax=docker/dockerfile:1

FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go mod verify
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -v -o ./go-server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /go
COPY --from=build /go/src/app /go
EXPOSE 8000
ENTRYPOINT /go/go-server --port 80