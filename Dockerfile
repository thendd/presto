FROM golang:1.23.5-bookworm AS base

WORKDIR /usr/presto

COPY . .

RUN go mod download

CMD [ "go", "run", "/usr/presto/cmd/presto/main.go" ]
