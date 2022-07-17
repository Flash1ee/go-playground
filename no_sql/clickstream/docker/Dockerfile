#FROM golang:1.17.1 as builder
#
#EXPOSE 80
#
#WORKDIR /app
#
#COPY . .
#
## Если что-то не собирается из-за CGO, может быть, при проверке сертификатов из гошки.
## Убрать CGO_ENABLED
## Итоговый image взять с gcc, например ubuntu
#
#RUN ls && go build -o server.out -v ./cmd/server
#
#FROM alpine
#
#COPY --from=builder /app /app
#
#WORKDIR /app
#
#CMD ls -la && /app/server.out

FROM golang:1.17.1 as builder

WORKDIR /app

EXPOSE 80

COPY . .

# Если что-то не собирается из-за CGO, может быть, при проверке сертификатов из гошки.
# Убрать CGO_ENABLED
# Итоговый image взять с gcc, например ubuntu
RUN CGO_ENABLED=0 make build

FROM alpine

COPY --from=builder /app /app

WORKDIR /app

CMD ls && ./server.out
