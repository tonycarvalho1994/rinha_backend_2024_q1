FROM golang:latest as builder

WORKDIR /srv

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY certs.crt /etc/ssl/certs

RUN update-ca-certificates

COPY . .

#ENV GOPROXY=direct
ENV CGO_ENABLED=0
RUN go mod tidy
RUN go build -o /srv/server -a -ldflags '-extldflags "-static" -s -w'

EXPOSE 8080

ENTRYPOINT [ "/srv/server" ]

FROM scratch
#
## Copiando o binário compilado da primeira etapa
COPY --from=builder /srv/server /server
#
EXPOSE 8080
#
## Definindo o comando padrão a ser executado
CMD ["/server"]