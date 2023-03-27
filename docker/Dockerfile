FROM golang:alpine3.17

RUN apk add build-base

WORKDIR /worker

ENV GIN_MODE=release

COPY . .

RUN go mod download && go build -tags musl -o /main

CMD [ "/main"]

