FROM golang:latest AS build

WORKDIR  /go/src
COPY . .

RUN GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build -o forum cmd/forum-api/main.go

FROM alpine
WORKDIR /app

ADD configs/docker.yml .

COPY --from=build /go/src/forum .
RUN chmod +x forum

ENTRYPOINT ["/app/forum"]

EXPOSE 8000/tcp