
# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/metal-stack/synk8s:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: synk8s

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

# Build manager binary
synk8s: edit download fmt vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o bin/synk8s main.go

GROUPVERSION ?= corev1
KIND ?= Secret
FIELD ?= Data
edit:
	sed 's#type Resource = .*#type Resource = ${GROUPVERSION}.${KIND}#' -i controllers/resource_controller.go && \
	sed 's#func set(source.*#func set(source, dest \*Resource) { source.${FIELD} = dest.${FIELD} }#' -i controllers/resource_controller.go && \
	go mod tidy

download:
	go mod download

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
docker-build: edit synk8s
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}
