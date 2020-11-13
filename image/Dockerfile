FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL maintainer="ricardozanini@gmail.com"
LABEL io.openshift.expose-services="9000:http"

RUN microdnf update &&\
    microdnf install shadow-utils &&\
    microdnf clean all

RUN adduser -ms /bin/bash -d /home/kogito -u 1001 kogito

USER kogito
WORKDIR /home/kogito

COPY kogito-sw-backend /usr/local/bin/kogito-sw-backend
#COPY bin/kogito /usr/local/bin
#COPY bin/kubectl /usr/local/bin
COPY --chown=kogito:kogito image/data/* /home/kogito/data

EXPOSE 9000

ENTRYPOINT ["/usr/local/bin/kogito-sw-backend"]