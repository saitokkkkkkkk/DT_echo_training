# Dockerfile
FROM golang:1.17.7-alpine

# ログに出力する時間をJSTにするため、タイムゾーンを設定
ENV TZ /usr/share/zoneinfo/Asia/Tokyo

RUN apk update && apk add git
RUN go get github.com/cosmtrek/air@latest

ENV ROOT=/go/src/app
WORKDIR ${ROOT}/src

# コーディングしないファイルはCOPYしローカルと同期しないようにする
# GoLandとコンテナの兼ね合いでgo.modはそれぞれの環境で中身が異なるため、特に同期させたくない
COPY ./app/go.mod ${ROOT}/go.mod

# ModuleモードをON
ENV GO111MODULE=on