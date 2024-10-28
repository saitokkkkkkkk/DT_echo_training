# Dockerfile

# Goのバージョンを1.23に更新
FROM golang:1.23-alpine

# タイムゾーン設定
ENV TZ /usr/share/zoneinfo/Asia/Tokyo

RUN apk update && apk add git

# Airのインストール
RUN go install github.com/air-verse/air@latest

# 作業ディレクトリの設定
ENV ROOT=/go/src/app
WORKDIR ${ROOT}/src

# go.modのコピー
COPY ./app/go.mod ${ROOT}/go.mod

# ModuleモードをON
ENV GO111MODULE=on