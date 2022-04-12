
# Image URL to use all building/pushing image targets
IMG ?= docker.io/kvaps/dnsmasq-controller:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

CONTROLLER_GEN = go run sigs.k8s.io/controller-tools/cmd/controller-gen@v0.8.0

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/dnsmasq-controller main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd/bases

# Uninstall CRDs from a cluster
uninstall: manifests
	kubectl delete -f config/crd/bases

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	# RBAC
	kubectl apply -n default \
	  -f config/rbac/service_account.yaml \
	  -f config/rbac/role.yaml \
	  -f config/rbac/role_binding.yaml \
	  -f config/rbac/leader_election_role.yaml \
	  -f config/rbac/leader_election_role_binding.yaml
	# DNS-server (for infra.example.org)
	kubectl apply -n default -f config/controller/dns-server.yaml
	# DHCP-server
	kubectl apply -n default -f config/controller/dhcp-server.yaml

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	$(CONTROLLER_GEN) crd rbac:roleName=dnsmasq-controller paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate:
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}
