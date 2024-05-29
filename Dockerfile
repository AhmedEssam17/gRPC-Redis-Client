FROM golang:latest

WORKDIR /app

RUN git clone https://github.com/AhmedEssam17/gRPC-Redis-Client.git

WORKDIR /app/gRPC-Redis-Client/client

RUN go mod download \
    && go build -o client .

EXPOSE 8888

CMD ["./client"]
