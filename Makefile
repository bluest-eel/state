VERSION = $(shell cat VERSION)
DOCKER_ORG = bluesteelabm

default: up

up:
	@VERSION=$(VERSION) RELAYD_VERSION=$(WS_RELAY_VERSION) \
	docker-compose -f ./docker/dev/compose.yml up

rebuild:
	@VERSION=$(VERSION) RELAYD_VERSION=$(WS_RELAY_VERSION) \
	docker-compose -f ./docker/dev/compose.yml build

rebuild-up: rebuild up

down:
	@VERSION=$(VERSION) RELAYD_VERSION=$(WS_RELAY_VERSION) \
	docker-compose -f ./docker/dev/compose.yml down

sqlsh: NODE ?= db1
sqlsh:
	@echo '>> Connecting to db $(NODE) ...'
	@docker exec -it $(NODE) ./cockroach sql --insecure

bash: NODE ?= db1
bash:
	@docker exec -it $(NODE) bash

$(WS_RELAY_CODE_DIR):
	@cd $(WS_RELAY_DIR) && \
	git clone $(WS_RELAY_REPO) $(WS_RELAY_CODE_NAME) && \
	git checkout v$(WS_RELAY_VERSION)

$(WS_RELAY_RENAME): $(WS_RELAY_CODE_DIR)
	@docker build -t $(DOCKER_ORG)/$(WS_RELAY_RENAME):$(WS_RELAY_VERSION) $(WS_RELAY_DIR)

images: $(WS_RELAY_RENAME)

tags:
	@docker tag $(DOCKER_ORG)/$(WS_RELAY_RENAME):$(WS_RELAY_VERSION) \
	$(DOCKER_ORG)/$(WS_RELAY_RENAME):latest

dockerhub: tags
	@docker push $(DOCKER_ORG)/$(WS_RELAY_RENAME):$(WS_RELAY_VERSION)
	@docker push $(DOCKER_ORG)/$(WS_RELAY_RENAME):latest

clean-docker:
	@docker system prune -f