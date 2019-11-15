
#GOPATH=$(PWD)/.go

VERSION=0.0.01

echo:
	@echo $(GOPATH)
	@echo "type make version to do release $(VERSION)"

version:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r sam3 -t v$(VERSION) -d "version $(VERSION)"


gx:
	go get github.com/whyrusleeping/gx
	go get github.com/whyrusleeping/gx-go

deps: gx
	gx --verbose install --global
	gx-go rewrite

publish:
	gx-go rewrite --undo

fmt: echo
	find . -path ./vendor -prune -o -name "*.go" -exec gofmt -w {} \;
	find . -path ./vendor -prune -o -name "*.i2pkeys" -exec rm {} \;

echobot:
	go build -o echo/echo ./echo/main.go

echorun:
	cd echo && ./echo

lint:
	golint *.go | less

vet:
	go vet *.go

test:
	go test

get:
	go get -u github.com/rtradeltd/go-garlic-tcp-transport/codec
	go get -u github.com/rtradeltd/go-garlic-tcp-transport/conn
	go get -u github.com/rtradeltd/go-garlic-tcp-transport/common
	go get -u github.com/rtradeltd/go-garlic-tcp-transport
