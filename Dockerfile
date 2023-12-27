FROM openeuler/openeuler:23.03 as BUILDER
RUN sed -i "s|repo.openeuler.org|mirrors.pku.edu.cn/openeuler|g" /etc/yum.repos.d/openEuler.repo && \ 
    dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

ARG USER
ARG PASS
RUN echo "machine github.com login $USER password $PASS" > /root/.netrc
RUN go env -w GOPRIVATE=github.com/opensourceways/xihe-extra-services,github.com/opensourceways/xihe-server,github.com/opensourceways/xihe-training-center,github.com/opensourceways/xihe-aicc-finetune,github.com/opensourceways/xihe-finetune,github.com/opensourceways/xihe-sync-repo

# build binary
COPY . /go/src/github.com/opensourceways/xihe-server
RUN cd /go/src/github.com/opensourceways/xihe-server && GO111MODULE=on CGO_ENABLED=0 go build -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'"
# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN sed -i "s|repo.openeuler.org|mirrors.pku.edu.cn/openeuler|g" /etc/yum.repos.d/openEuler.repo && \ 
    dnf -y update && \
    dnf in -y shadow && \
    dnf remove -y gdb-gdbserver && \
    groupadd -g 5000 mindspore && \
    useradd -u 5000 -g mindspore -s /sbin/nologin -m mindspore

RUN echo > /etc/issue && echo > /etc/issue.net && echo > /etc/motd
RUN mkdir /opt/app -p
RUN chmod 700 /opt/app
RUN chown mindspore:mindspore /opt/app

RUN echo 'set +o history' >> /root/.bashrc
RUN sed -i 's/^PASS_MAX_DAYS.*/PASS_MAX_DAYS   90/' /etc/login.defs
RUN rm -rf /tmp/*

USER mindspore
WORKDIR /opt/app/

COPY  --chown=mindspore --from=BUILDER /go/src/github.com/opensourceways/xihe-server/xihe-server /opt/app
COPY  --chown=mindspore ./points/infrastructure/taskdocimpl/doc_chinese.tmpl  /opt/app/points/task-docs-templates/doc_chinese.tmpl
COPY  --chown=mindspore ./points/infrastructure/taskdocimpl/doc_english.tmpl  /opt/app/points/task-docs-templates/doc_english.tmpl

RUN chmod 550 /opt/app/xihe-server
RUN chmod 640 /opt/app/points/task-docs-templates/doc_chinese.tmpl
RUN chmod 640 /opt/app/points/task-docs-templates/doc_english.tmpl
RUN chmod 750 /opt/app/points/task-docs-templates
RUN chmod 750 /opt/app/points

RUN echo "umask 027" >> /home/mindspore/.bashrc
RUN echo 'set +o history' >> /home/mindspore/.bashrc

ENTRYPOINT ["/opt/app/xihe-server"]
