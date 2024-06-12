# syntax=docker/dockerfile:1

FROM golang:1.22
WORKDIR /go-server
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go install -v ./... && go build

CMD ["./go-server"]
EXPOSE 8000