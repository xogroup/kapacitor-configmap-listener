setup:
	curl https://glide.sh/get | sh

build:
	glide install
	go build

.PHONY: install build