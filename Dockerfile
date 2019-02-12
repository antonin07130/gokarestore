FROM golang:1.11.0-alpine3.8 as build
RUN apk update && apk add gcc musl-dev librdkafka-dev git

ENV GO111MODULE on
ENV p github.com/antonin07130/gokarestore
ADD . /go/src/${p}
RUN cd /go/src/${p} \
 && go test ./... \
 && go install ./...

FROM alpine:3.8

RUN apk add -U --no-cache librdkafka
COPY --from=build /go/bin/gokarestore /bin/
ENTRYPOINT ["bin/gokarestore"]
