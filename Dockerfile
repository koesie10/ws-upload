FROM golang:1.20.4 as builder

RUN mkdir /app
WORKDIR /app
COPY . /app/
RUN CGO_ENABLED=0 go build -o ws-upload ./cmd/ws-upload

FROM alpine:3.18.0
RUN apk --update add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/ws-upload /bin/ws-upload

EXPOSE 9108
ENTRYPOINT ["/bin/ws-upload"]
