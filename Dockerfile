FROM golang:alpine

COPY . .

RUN go build -o main ./app.go

CMD ["./main"]