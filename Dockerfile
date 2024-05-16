FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o message_sender .

EXPOSE 8080

ENTRYPOINT [ "./message_sender" ]
