FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL maintainer="ricardozanini@gmail.com"
LABEL io.openshift.expose-services="9000:http"

ENV HOME=/home/kogito

RUN microdnf update &&\
    microdnf install shadow-utils &&\
    microdnf clean all

RUN adduser -ms /sbin/nologin -d /home/kogito -g 0 -r -u 1001 kogito && mkdir /home/kogito/.kube

COPY bin/kogito-sw-backend /usr/local/bin/kogito-sw-backend
COPY bin/kogito /usr/local/bin
COPY bin/oc /usr/local/bin
COPY image/data /home/kogito/data

RUN chgrp -R 0 /home/kogito && chmod -R g+rwX /home/kogito

USER 1001
WORKDIR /home/kogito

EXPOSE 9000

ENTRYPOINT ["/usr/local/bin/kogito-sw-backend"]