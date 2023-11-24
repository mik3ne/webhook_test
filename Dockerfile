# syntax = docker/dockerfile:1.6

FROM golang:1.21-alpine AS builder

COPY . /var/go/

WORKDIR /var/go/cmd/server/

RUN go build -mod=vendor -x -v -o /webhook_server


FROM scratch AS runner

COPY --from=builder /webhook_server /
COPY --from=builder /var/go/config/config.yml /config/config.yml

ENTRYPOINT [ "/webhook_server" ]