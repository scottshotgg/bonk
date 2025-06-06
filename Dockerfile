FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 go build -o bonk

FROM alpine:3.22.0

COPY --from=builder /app/bonk /bin/bonk

CMD [ "/bin/bonk" ]
