FROM golang:1.20.1-bullseye

WORKDIR /workspace

COPY . .

WORKDIR /workspace/plugins

RUN go build -o main

CMD ["./main"]