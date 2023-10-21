FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

# build binary
COPY . /go/src/github.com/opensourceways/xihe-server
RUN cd /go/src/github.com/opensourceways/xihe-server && GO111MODULE=on CGO_ENABLED=0 go build -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'"
# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow && \
    groupadd -g 5000 mindspore && \
    useradd -u 5000 -g mindspore -s /bin/bash -m mindspore

USER mindspore
WORKDIR /opt/app/

COPY  --chown=mindspore --from=BUILDER /go/src/github.com/opensourceways/xihe-server/xihe-server /opt/app
COPY  --chown=mindspore ./points/infrastructure/taskdocimpl/doc_chinese.tmpl  /opt/app/points/task-docs-templates/doc_chinese.tmpl
COPY  --chown=mindspore ./points/infrastructure/taskdocimpl/doc_english.tmpl  /opt/app/points/task-docs-templates/doc_english.tmpl

RUN chmod 550 /opt/app/xihe-server
RUN chmod 640 /opt/app/points/task-docs-templates/doc_chinese.tmpl
RUN chmod 750 /opt/app/points/task-docs-templates
RUN chmod 640 /opt/app/points/task-docs-templates/doc_english.tmpl
RUN chmod 750 /opt/app/points/task-docs-templates

ENTRYPOINT ["/opt/app/xihe-server"]
