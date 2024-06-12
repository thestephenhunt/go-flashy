# syntax=docker/dockerfile:1

FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go-server
COPY . .
RUN go mod download && go mod verify
RUN go mod tidy
RUN GOOS=linux go build -o ./bin/go-server ./main.go

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go-server /go/bin
EXPOSE 80
ENTRYPOINT /go/bin/go-server --port 80