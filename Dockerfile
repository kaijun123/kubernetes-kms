FROM golang:1.20.1-bullseye

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

# Download all dependencies into the go mod cache
RUN go mod download

ENV IP_ADDRESS "192.168.10.42"

COPY . .

WORKDIR /workspace/pkg/plugins

RUN go build -o main

CMD ["./main"]