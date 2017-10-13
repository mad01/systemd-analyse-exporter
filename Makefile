PROJ=systemd-analyse-exporter
ORG_PATH=github.com/mad01
REPO_PATH=$(ORG_PATH)/$(PROJ)

VERSION ?= $(shell ./scripts/git-version)
LD_FLAGS="-X main.Version=$(VERSION) -extldflags \"-static\" "
version.Version=$(VERSION)
$( shell mkdir -p _bin )
$( shell mkdir -p _release )

export GOBIN=$(PWD)/_bin


default: build

clean:
	@rm -r _bin _release

test:
	@go test -v -i $(shell go list ./... | grep -v '/vendor/')
	@go test -v $(shell go list ./... | grep -v '/vendor/')

build: build/dev

build/dev:
	@go install -v -ldflags $(LD_FLAGS) 

build/release:
	@go build -v -o _release/$(PROJ) -ldflags $(LD_FLAGS) 


docker/build:
	@docker build -t quay.io/mad01/$(PROJ):$(VERSION) --file Dockerfile .

docker/push:
	@docker push quay.io/mad01/$(PROJ):$(VERSION)

docker/login:
	@docker login -u $(QUAY_LOGIN) -p="$(QUAY_PASSWORD)" quay.io
