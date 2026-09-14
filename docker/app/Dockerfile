FROM golang:1.26-alpine

WORKDIR /go/src

RUN apk upgrade --update && \
    apk --no-cache add git

RUN go install github.com/air-verse/air@latest

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src .

CMD ["air", "-c", ".air.toml"]
