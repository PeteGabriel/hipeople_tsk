FROM golang:1.17
ENV GO111MODULE=on

ENV PORT=4002

WORKDIR /app

RUN mkdir -p /app/static/images

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o bin/main cmd/api/main.go
CMD ./bin/main -port ${PORT}