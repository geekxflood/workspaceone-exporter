IMAGE_NAME = $(shell basename $(shell pwd))
TAG_LATEST = $(IMAGE_NAME):latest
TAG_DATE = $(IMAGE_NAME):$(shell date +"%Y-%m-%dT%H_%M")

.PHONY: all image binary timestamp clean

all: image build

image:
	docker build . \
		-t $(TAG_LATEST) \
		-t $(TAG_DATE)

build:
	mkdir -p binary
	go build -a  \
		-gcflags=all="-l -B" \
		-ldflags="-w -s" \
		-o binary/$(IMAGE_NAME) \
		./...

timestamp:
	mkdir -p binary
	go build -a \
		-gcflags=all="-l -B" \
		-ldflags="-w -s" \
		-o binary/$(IMAGE_NAME)_$(shell date +"%Y-%m-%dT%H_%M") \
		./...

clean:
	rm -rf binary
