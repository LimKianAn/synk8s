
# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/metal-stack/sync-crd:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: sync-crd

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

REPO_URL ?= github.com/metal-stack/firewall-controller
REPO_VERSION ?= latest
SUB_PATH ?= api/v1
CRD_KIND ?= ClusterwideNetworkPolicy

sync-crd: manager back
	mv bin/manager bin/sync-crd

# Build manager binary
manager: edit download fmt vet
	go build -o bin/manager main.go

edit:
	sed 's#=> .*#=> ${REPO_URL} ${REPO_VERSION}#' -i go.mod && \
	sed 's#repo-url/.*#repo-url/${SUB_PATH}"#' -i main.go && \
	sed 's#crd = api..*#crd = api.${CRD_KIND}#' -i main.go && \
	go mod tidy

download:
	go mod download

back:
	sed 's#=> .*#=> REPO_URL REPO_VERSION#' -i go.mod && \
	sed 's#repo-url/.*#repo-url/SUB_PATH"#' -i main.go && \
	sed 's#crd = api..*#crd = api.CRD_KIND#' -i main.go

# stall CRDs into a cluster
install:
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall:
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy:
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build the docker image
docker-build: edit test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}
