FROM golang:latest AS build

WORKDIR  /app
COPY . .

RUN GO111MODULE=on \
  CGO_ENABLED=0 \
  go build -o proxy ./cmd/proxy

FROM alpine

WORKDIR /app

RUN apk upgrade --update-cache --available && apk add \
  openssl \
  ca-certificates \
  && rm -rf /var/cache/apk/*

ADD ssl ssl
ADD scripts/gen_cert.sh ssl/
RUN chmod +x ssl/gen_cert.sh

RUN cp ssl/ca.crt /etc/ssl
RUN update-ca-certificates

COPY --from=build /app/proxy .
RUN chmod +x proxy

ENTRYPOINT ["/app/proxy"]

EXPOSE 8000/tcp
EXPOSE 8080/tcp
