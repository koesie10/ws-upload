FROM golang:1.23.6 AS builder

RUN mkdir /app
WORKDIR /app
COPY . /app/
RUN CGO_ENABLED=0 go build -o ws-upload ./cmd/ws-upload

FROM alpine:3.21.3
RUN apk --update add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/ws-upload /bin/ws-upload

EXPOSE 9108
ENTRYPOINT ["/bin/ws-upload"]
