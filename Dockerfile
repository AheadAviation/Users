FROM golang:1.10.3-alpine3.8 AS builder
LABEL maintainer=<tim.curless@thinkahead.com>

COPY ./ /go/src/github.com/aheadaviation/users/
WORKDIR /go/src/github.com/aheadaviation/users/

RUN apk update && apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build -o /go/bin/users

FROM scratch

ENV MONGO_HOST mytestdb:27017
ENV HATEAOS users
ENV USER_DATABASE mongodb

COPY --from=builder /go/bin/users /go/bin/users
EXPOSE 8084
ENTRYPOINT ["/go/bin/users"]
