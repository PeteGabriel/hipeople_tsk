FROM golang:1.17
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/main cmd/api/main.go
CMD [ "./bin/main" ]