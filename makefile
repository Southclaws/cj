VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
-include .env

.PHONY: version

fast:
	go build $(LDFLAGS) -o cj

static:
	CGO_ENABLED=0 GOOS=linux go build -a $(LDFLAGS) -o cj .

local: fast
	DEBUG=1 \
	BIND=localhost:8080 \
	MONGO_USER=$(MONGO_USER) \
	MONGO_PASS=$(MONGO_PASS) \
	MONGO_HOST=$(MONGO_HOST) \
	MONGO_PORT=$(MONGO_PORT) \
	MONGO_NAME=$(MONGO_NAME) \
	DISCORD_TOKEN=$(DISCORD_TOKEN) \
	ADMINISTRATIVE_CHANNEL=$(ADMINISTRATIVE_CHANNEL) \
	PRIMARY_CHANNEL=$(PRIMARY_CHANNEL) \
	HEARTBEAT=$(HEARTBEAT) \
	BOT_ID=$(BOT_ID) \
	GUILD_ID=$(GUILD_ID) \
	VERIFIED_ROLE=$(VERIFIED_ROLE) \
	NORMAL_ROLE=$(NORMAL_ROLE) \
	ADMIN=$(ADMIN) \
	LANGUAGE_DATA=$(LANGUAGE_DATA) \
	LANGUAGE=$(LANGUAGE) \
	NO_INIT_SYNC=$(NO_INIT_SYNC) \
	./cj

version:
	git tag $(VERSION)
	git push
	git push origin $(VERSION)

test:
	go test -v -race

# Docker

build:
	docker build --no-cache -t southclaws/cj:$(VERSION) -f Dockerfile.dev .

build-prod:
	docker build --no-cache -t southclaws/cj:$(VERSION) .

build-test:
	docker build --no-cache -t southclaws/cj-test:$(VERSION) -f Dockerfile.testing .

push: build-prod
	docker push southclaws/cj:$(VERSION)
	
run:
	-docker rm cj-test
	docker run \
		--name cj-test \
		--network host \
		-e BIND=0.0.0.0:8080 \
		-e MONGO_USER=$(MONGO_USER) \
		-e MONGO_HOST=$(MONGO_HOST) \
		-e MONGO_PORT=$(MONGO_PORT) \
		-e MONGO_NAME=$(MONGO_NAME) \
		-e MONGO_COLLECTION=$(MONGO_COLLECTION) \
		-e DISCORD_TOKEN=$(DISCORD_TOKEN) \
		-e ADMINISTRATIVE_CHANNEL=$(ADMINISTRATIVE_CHANNEL) \
		-e PRIMARY_CHANNEL=$(PRIMARY_CHANNEL) \
		-e HEARTBEAT=$(HEARTBEAT) \
		-e BOT_ID=$(BOT_ID) \
		-e GUILD_ID=$(GUILD_ID) \
		-e VERIFIED_ROLE=$(VERIFIED_ROLE) \
		-e NORMAL_ROLE=$(NORMAL_ROLE) \
		-e DEBUG_USER=$(DEBUG_USER) \
		-e ADMIN=$(ADMIN) \
		-e LANGUAGE_DATA="/cjlang" \
		-e LANGUAGE=($LANGUAGE) \
		-e DEBUG=1 \
		-e NO_INIT_SYNC=1 \
		southclaws/cj:$(VERSION)

enter:
	docker run -it --entrypoint=bash southclaws/cj:$(VERSION)

enter-mount:
	docker run -v $(shell pwd)/testspace:/samp -it --entrypoint=bash southclaws/cj:$(VERSION)

# Test stuff

test-container: build-test
	docker run --network host southclaws/cj-test:$(VERSION)

mongodb:
	-docker stop mongodb
	-docker rm mongodb
	docker run --name mongodb -p 27017:27017 -d mongo
