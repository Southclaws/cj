VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
MONGO_PASS := $(shell cat MONGO_PASS.private)
DISCORD_TOKEN := $(shell cat DISCORD_TOKEN.private)

.PHONY: version

fast:
	go build $(LDFLAGS) -o cj

static:
	CGO_ENABLED=0 GOOS=linux go build -a $(LDFLAGS) -o cj .

local: fast
	DEBUG=1 \
	BIND=localhost:8080 \
	MONGO_USER=root \
	MONGO_PASS=$(MONGO_PASS) \
	MONGO_HOST=localhost \
	MONGO_PORT=27017 \
	MONGO_NAME=cj \
	DISCORD_TOKEN=$(DISCORD_TOKEN) \
	ADMINISTRATIVE_CHANNEL="282581078643048448" \
	PRIMARY_CHANNEL="231799104731217931" \
	HEARTBEAT="10" \
	BOT_ID="285421343594512384" \
	GUILD_ID="231799104731217931" \
	VERIFIED_ROLE="285459413882634241" \
	NORMAL_ROLE="285459500029444096" \
	ADMIN="86435690711093248" \
	LANGUAGE_DATA="./lang" \
	NO_INIT_SYNC="1" \
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
		-e BIND=localhost:8080 \
		-e MONGO_USER=cj \
		-e MONGO_HOST=localhost \
		-e MONGO_PORT=27017 \
		-e MONGO_NAME=cj \
		-e MONGO_COLLECTION=servers \
		-e DISCORD_TOKEN=$(DISCORD_TOKEN) \
		-e ADMINISTRATIVE_CHANNEL="282581078643048448" \
		-e PRIMARY_CHANNEL="231799104731217931" \
		-e HEARTBEAT="10" \
		-e BOT_ID="285421343594512384" \
		-e GUILD_ID="231799104731217931" \
		-e VERIFIED_ROLE="285459413882634241" \
		-e NORMAL_ROLE="285459500029444096" \
		-e DEBUG_USER="86435690711093248" \
		-e ADMIN="86435690711093248" \
		-e LANGUAGE_DATA="/cjlang" \
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
