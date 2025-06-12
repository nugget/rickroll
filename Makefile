PROJECT?=rickrolld
REGISTRY?=docker.io
LIBRARY=nugget

image=$(REGISTRY)/$(LIBRARY)/$(PROJECT)

platforms=linux/amd64,linux/arm64

BINARYNAME?=rickrolld

prodtag=latest
devtag=dev
builder=builder-$(PROJECT)

GIT ?= $(shell which git)
PWD ?= $(shell pwd)

OCI_IMAGE_CREATED="$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")"
OCI_IMAGE_REVISION?="$(shell $(GIT) rev-parse HEAD)"
OCI_IMAGE_VERSION?="$(shell $(GIT) describe --always --long --tags --dirty)"

oci-build-labels?=\
	--build-arg OCI_IMAGE_CREATED=$(OCI_IMAGE_CREATED) \
	--build-arg OCI_IMAGE_VERSION=$(OCI_IMAGE_VERSION) \
	--build-arg OCI_IMAGE_REVISION=$(OCI_IMAGE_REVISION) 

oci-version-point=$(shell echo $(OCI_IMAGE_VERSION) | cut -f 2 -d v | cut -f 1 -d '-')
oci-version-minor=$(shell echo $(oci-version-point) | cut -f 1-2 -d .)
oci-version-major=$(shell echo $(oci-version-point) | cut -f 1 -d .)

.PHONY: version mod go-telnet-local localdev productiondev rickrolld run container runcontainer clean buildx release

version:
	@echo "OCI Version: $(OCI_IMAGE_VERSION)"

mod:
	go get -u github.com/nugget/go-telnet
	go mod tidy
	go mod verify
	git diff go.mod go.sum || git commit go.mod go.sum -m "make mod"

go-telnet-local:
	-git clone git@github.com:nugget/go-telnet.git

localdev: go-telnet-local
	go mod edit -replace=github.com/nugget/go-telnet="$(PWD)/go-telnet"
	go mod tidy

productiondev:
	go mod edit -dropreplace=github.com/nugget/go-telnet
	go mod tidy

rickrolld: 
	mkdir -p dist
	go mod tidy
	cd rickrolld && CGO_ENABLED=0 go build -o ../dist/$(BINARYNAME) .

run: rickrolld
	./dist/$(BINARYNAME) -v

container:
	docker build $(oci-build-labels) . -t nugget/rickrolld:dev --load

runcontainer: container
	docker run -p 23:23 nugget/rickrolld:dev

clean:
	@echo "# making: clean"
	-docker buildx rm $(builder)
	@echo

buildx: mod
	docker buildx create --name $(builder)
	docker buildx use $(builder)
	docker buildx install
	@echo

release: version buildx
	@echo "# making: prod"
	docker buildx use $(builder)
	docker buildx build $(oci-build-labels) \
		-t $(image):$(prodtag) \
		-t $(image):$(oci-version-major) \
		-t $(image):$(oci-version-minor) \
		-t $(image):$(oci-version-point) \
		--platform=$(platforms) --push . 
	docker buildx rm $(builder)
	docker pull $(image):$(prodtag)
	docker inspect $(image):$(prodtag) | jq '.[0].Config.Labels' 
	@echo 

