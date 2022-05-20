FROM golang:1.18 as builder

WORKDIR /app

COPY . /app/

RUN go mod tidy

RUN  CGO_ENABLED=0 go build -o main main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]