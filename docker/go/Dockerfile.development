FROM golang:1.18

RUN apt-get update

WORKDIR /go/src/server

COPY ./server /go/src/server

# ホットリロード
RUN go install github.com/pilu/fresh@latest

# マイグレーションツールの導入
RUN go install bitbucket.org/liamstask/goose/cmd/goose@latest
