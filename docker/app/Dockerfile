FROM golang:1.16-alpine

WORKDIR /go/src
COPY ./src .

RUN apk upgrade --update && \
    apk --no-cache add git

RUN go get -u github.com/cosmtrek/air && \
    go build -o /go/bin/air github.com/cosmtrek/air

RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/joho/godotenv

CMD ["air", "-c", ".air.toml"]