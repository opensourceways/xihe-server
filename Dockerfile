FROM golang:latest as BUILDER

# build binary
COPY . /go/src/github.com/opensourceways/xihe-server
RUN cd /go/src/github.com/opensourceways/xihe-server && GO111MODULE=on CGO_ENABLED=0 go build

# copy binary config and utils
FROM alpine:latest

RUN adduser mindspore -u 5000 -D
USER mindspore
WORKDIR /opt/app/

COPY  --from=BUILDER /go/src/github.com/opensourceways/xihe-server/xihe-server /opt/app
COPY  ./points/infrastructure/taskdocimpl/doc_chinese.tmpl  /opt/app/points/task-docs-templates/doc_chinese.tmpl
COPY  ./points/infrastructure/taskdocimpl/doc_english.tmpl  /opt/app/points/task-docs-templates/doc_english.tmpl

ENTRYPOINT ["/opt/app/xihe-server"]
