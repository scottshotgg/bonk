FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 go build -o bonk

FROM alpine:3.20.3

COPY --from=builder /app/bonk /bin/bonk

CMD [ "/bin/bonk" ]
