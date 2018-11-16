FROM golang:1.10.3-alpine3.8 AS builder
LABEL maintainer=<tim.curless@thinkahead.com>

COPY ./ /go/src/github.com/aheadaviation/users/
WORKDIR /go/src/github.com/aheadaviation/users/

RUN apk update && apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure -v
RUN CGO_ENABLED=0 go build -o /go/bin/users

FROM alpine:3.8

ENV MONGO_HOST db-users \
    HATEAOS users \
    USERS_DATABASE mongodb

HEALTHCHECK --interval=10s CMD wget -q0- localhost:8084/health

COPY --from=builder /go/bin/users /usr/local/bin/users
RUN chmod +x /usr/local/bin/users
EXPOSE 8084
ENTRYPOINT ["users"]
