##########################################
FROM golang:alpine as builder

RUN apk add --no-cache git curl

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/atos/rotterdam

COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -o Rotterdam .

##########################################
FROM alpine:3.6
WORKDIR /opt/rotterdam

COPY --from=builder /go/src/atos/rotterdam/Rotterdam .
COPY ./config/config.json ./config/config.json
COPY ./config/qos_templates.json ./config/qos_templates.json
COPY ./rest-api ./rest-api

# install curl
RUN apk add --no-cache curl


EXPOSE 8333
ENTRYPOINT ["/opt/rotterdam/Rotterdam"]