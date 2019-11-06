FROM golang:1.13 as golang
ADD . /src/
RUN set -x && \
  cd /src && \
  CGO_ENABLED=0 GOOS=linux go build -o ngfaas_server server_one/main.go

FROM alpine:3.10
RUN apk add --no-cache ca-certificates

RUN addgroup -g 2000 apprunner
RUN adduser -u 2000 -G apprunner -S apprunner

COPY --chown=apprunner:apprunner --from=golang /src/ngfaas_server /app/ngfaas_server

WORKDIR /app

USER apprunner

CMD ["/app/ngfaas_server"]
