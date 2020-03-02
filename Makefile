DOCKER_REPO ?= uzxmx/cert-manager-webhook-alidns
DOCKER_TAG  := $(if $(DOCKER_TAG),$(DOCKER_TAG),latest)

.PHONY: docker-build
docker-build:
	docker build -t "$(DOCKER_REPO):$(DOCKER_TAG)" .

.PHONY: docker-push
docker-push:
	docker push "$(DOCKER_REPO):$(DOCKER_TAG)"

.PHONY: docker-login
docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
