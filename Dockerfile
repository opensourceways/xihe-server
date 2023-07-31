FROM golang:latest as BUILDER

# build binary
COPY . /go/src/github.com/opensourceways/xihe-server
RUN cd /go/src/github.com/opensourceways/xihe-server && GO111MODULE=on CGO_ENABLED=0 go build

RUN groupadd --gid 5000 mindspore \
  && useradd --home-dir /home/mindspore --create-home --uid 5000 --gid 5000 --shell /bin/sh --skel /dev/null mindspore

USER mindspore

# copy binary config and utils
FROM alpine:latest
WORKDIR /usr/src/app

COPY  --from=BUILDER /go/src/github.com/opensourceways/xihe-server/xihe-server /usr/src/app

ENTRYPOINT ["/usr/src/app/xihe-server"]
