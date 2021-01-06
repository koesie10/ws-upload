FROM golang:1.15.6 as builder

RUN mkdir /app
WORKDIR /app
COPY . /app/
RUN CGO_ENABLED=0 go build -o ws-upload .

FROM alpine
RUN apk --update add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/ws-upload /bin/ws-upload

EXPOSE 9108
ENTRYPOINT ["/bin/ws-upload"]
