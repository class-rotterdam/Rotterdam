##########################################
FROM golang:alpine as builder

ARG VERSION
ARG DATE

RUN apk add --no-cache git curl

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/SLALite

COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -o SLALite .

##########################################
FROM alpine:3.6

# /opt/slalite
WORKDIR /opt/slalite
COPY --from=builder /go/src/SLALite/SLALite .
COPY docker/run_slalite_mongo.sh /opt/slalite/run_slalite.sh
RUN chmod 775 /opt/slalite/run_slalite.sh

# install curl
RUN apk add --no-cache curl

# /etc/slalite
RUN mkdir /etc/slalite
COPY docker/slalite_mongo.yml /etc/slalite/slalite.yml
RUN addgroup -S slalite && adduser -D -G slalite slalite
RUN chown -R slalite:slalite /etc/slalite && chmod 775 /etc/slalite
RUN chmod 775 /etc/slalite/slalite.yml

EXPOSE 8090
USER slalite
ENTRYPOINT ["/opt/slalite/SLALite"]

