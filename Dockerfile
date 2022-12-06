FROM golang:1.19.3 AS builder

COPY . /app
WORKDIR /app
RUN go build

FROM alpine:latest

COPY --from=builder /app/ferp /app/ferp

CMD [ "/app/ferp" ]
