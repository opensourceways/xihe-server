FROM golang:latest as BUILDER

# build binary
COPY . /go/src/github.com/opensourceways/xihe-server
RUN cd /go/src/github.com/opensourceways/xihe-server/internal-server && GO111MODULE=on CGO_ENABLED=0 go build

# copy binary config and utils
FROM alpine:latest

RUN adduser mindspore -u 5000 -D
USER mindspore
WORKDIR /opt/app/

COPY  --from=BUILDER /go/src/github.com/opensourceways/xihe-server/internal-server/internal-server /opt/app

ENTRYPOINT ["/opt/app/internal-server"]
