FROM golang:latest as BUILDER

# build binary
COPY . /go/src/github.com/opensourceways/xihe-server
RUN cd /go/src/github.com/opensourceways/xihe-server/async-server && GO111MODULE=on CGO_ENABLED=0 go build

# copy binary config and utils
FROM alpine:latest

RUN adduser mindspore -u 5000 -D
USER mindspore
WORKDIR /home/mindspore/app

ENV BIGMODEL_WUKONGHF="test"

COPY  --from=BUILDER /go/src/github.com/opensourceways/xihe-server/async-server/async-server /home/mindspore/app

ENTRYPOINT ["/home/mindspore/app/async-server"]
