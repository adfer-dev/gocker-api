FROM golang:alpine

LABEL maintainer="adferdev"

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o /build

EXPOSE 8080

CMD [ "/build" ]