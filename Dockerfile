FROM golang:1.23.5-alpine3.21 AS builder

WORKDIR /root/build

COPY . .

RUN go mod download \
    && go mod verify

RUN go test ./... \
    && go build -o app ./cmd/presto/main.go

FROM alpine:3.21

RUN adduser -D -s /bin/sh presto

USER presto

WORKDIR /home/presto

COPY --chown=presto:presto --from=builder /root/build/app .

CMD [ "/home/presto/app", "--debug_flags=$DEBUG_FLAGS" ]