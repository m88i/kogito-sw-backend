SHELL := /bin/bash

KOGITO_VERSION ?= 0.17.0
KUBECTL_VERSION ?= 1.19.0
OKD_VERSION ?= 4.5.0-0.okd-2020-10-15-235428

all: build-image

build:
	go build -o bin/kogito-sw-backend

kubectl:
ifeq (,$(wildcard bin/kubectl))
	curl -LO "https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl"
	chmod +x ./kubectl
	mv ./kubectl bin/kubectl
endif

oc:
ifeq (,$(wildcard bin/oc))
	curl -LO "https://github.com/openshift/okd/releases/download/4.5.0-0.okd-2020-10-15-235428/openshift-client-linux-${OKD_VERSION}.tar.gz"
	tar -xvzf "openshift-client-linux-${OKD_VERSION}.tar.gz"
	mv ./oc bin/oc
	rm -rf "openshift-client-linux-${OKD_VERSION}.tar.gz"
endif

kogito:
ifeq (,$(wildcard bin/kogito))
	curl -LO "https://github.com/kiegroup/kogito-cloud-operator/releases/download/v${KOGITO_VERSION}/kogito-cli-${KOGITO_VERSION}-linux-amd64.tar.gz"
	tar -xvzf "kogito-cli-${KOGITO_VERSION}-linux-amd64.tar.gz"
	mv ./kogito bin/kogito
	rm -rf "kogito-cli-${KOGITO_VERSION}-linux-amd64.tar.gz"
endif

build-image: build oc kogito
	podman build --tag quay.io/m88i/kogito-sw-backend:latest -f image/Dockerfile .

push:
	podman push quay.io/m88i/kogito-sw-backend:latest
