FROM golang

RUN go env -w GO111MODULE=off
RUN go get github.com/gin-gonic/gin
RUN go get github.com/stretchr/testify/assert
RUN go get github.com/go-redis/redis
RUN go get github.com/chilts/sid

# ENV GIN_MODE=release

COPY ./app /go/src/cata/app
WORKDIR /go/src/cata/app

RUN go get ./
RUN go build -o app

CMD ./app
EXPOSE 8080

