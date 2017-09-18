setup:
	go get github.com/Masterminds/glide
	cd $$GOPATH/src/github.com/Masterminds/glide
	go install
	cd -

build:
	glide install
	go build

.PHONY: install