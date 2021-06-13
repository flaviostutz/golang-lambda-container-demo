FROM golang:1.16.5-alpine3.13 AS BUILD

# RUN apk add gcc build-base

ENV LOG_LEVEL 'info'
ENV READONLY 'false'
#endpoint: http, lambda
ENV ENDPOINT 'lambda'

WORKDIR /app

ADD /go.mod /app/
ADD /go.sum /app/

RUN go mod download

ADD /startup.sh /
ADD / /app/

# RUN go test -v -p 1

RUN GOOS=linux GOARCH=amd64 go build -o /bin/golang-lambda-container-demo

CMD [ "/startup.sh" ]

