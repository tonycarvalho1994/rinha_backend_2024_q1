FROM golang:latest as builder

WORKDIR /srv

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY certs.crt /etc/ssl/certs

RUN update-ca-certificates

COPY . .

ENV CGO_ENABLED=0
RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o /srv/server -a -ldflags '-extldflags "-static" -s -w'

FROM scratch
WORKDIR /
COPY --from=builder /srv/server /server
EXPOSE 8080
CMD ["/server"]