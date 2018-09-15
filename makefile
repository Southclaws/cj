VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
-include .env

.PHONY: version

fast:
	go build $(LDFLAGS) -o cj

static:
	CGO_ENABLED=0 GOOS=linux go build -a $(LDFLAGS) -o cj .

local: fast
	./cj

version:
	git tag $(VERSION)
	git push
	git push origin $(VERSION)

test:
	go test -v -race


# -
# Docker
#-


build:
	docker build --no-cache -t southclaws/cj:$(VERSION) .

push:
	docker push southclaws/cj:$(VERSION)
	docker tag southclaws/cj:$(VERSION) southclaws/cj:latest
	docker push southclaws/cj:latest
	
run:
	-docker rm cj
	docker run \
		--name cj \
		--network host \
		--env-file .env \
		southclaws/cj:$(VERSION)

run-prod:
	-docker stop cj
	-docker rm cj
	docker run \
		--name cj \
		--detach \
		--restart always \
		--env-file .env \
		southclaws/cj:$(VERSION)
	docker network connect mongodb cj

enter:
	docker run -it --entrypoint=bash southclaws/cj:$(VERSION)

enter-mount:
	docker run -v $(shell pwd)/testspace:/samp -it --entrypoint=bash southclaws/cj:$(VERSION)

# Test stuff

test-container: build-test
	docker run --network host southclaws/cj:$(VERSION)

mongodb:
	-docker stop mongodb
	-docker rm mongodb
	docker run --name mongodb -p 27017:27017 -d mongo
