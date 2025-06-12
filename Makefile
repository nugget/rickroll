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

.PHONY: mod go-telnet-local localdev productiondev rickrolld run container runcontainer clean buildx release

mod:
	go get -u github.com/nugget/go-telnet
	go mod tidy
	go mod verify
	git diff go.mod go.sum && git commit go.mod go.sum -m "make mod"

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

buildx:
	docker buildx create --name $(builder)
	docker buildx use $(builder)
	docker buildx install
	@echo

release: buildx
	@echo "# making: prod"
	docker buildx use $(builder)
	docker buildx build $(oci-build-labels) -t $(image):$(prodtag) $(FROM_IMAGE_TAGARGS) --platform=$(platforms) --push . 
	docker buildx rm $(builder)
	docker pull $(image):$(prodtag)
	docker inspect $(image):$(prodtag) | jq '.[0].Config.Labels' 
	@echo 

