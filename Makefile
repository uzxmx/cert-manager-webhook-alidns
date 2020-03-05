DOCKER_REPO ?= uzxmx/cert-manager-webhook-alidns
DOCKER_TAG  := $(if $(DOCKER_TAG),$(DOCKER_TAG),latest)

ifneq (,$(shell which gsed))
SED := gsed
else
SED := sed
endif

ifeq (,$(shell $(SED) --version 2>/dev/null | sed -n 1p | grep 'GNU sed'))
$(error You must install gnu sed)
endif

.PHONY: docker-build
docker-build:
	docker build -t "$(DOCKER_REPO):$(DOCKER_TAG)" .

.PHONY: docker-push
docker-push:
	docker push "$(DOCKER_REPO):$(DOCKER_TAG)"

.PHONY: docker-login
docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

.PHONY: dist
dist:
	mkdir -p dist
	cp -R deploy/chart dist/alidns
	[ -n "$(VERSION)" ] && $(SED) -i -Ee 's/^(version:).*$$/\1 $(VERSION)/' dist/alidns/Chart.yaml && \
		$(SED) -i -Ee 's/^(  tag:).*$$/\1 $(VERSION)/' dist/alidns/values.yaml
	(cd dist && tar zcf alidns-$(VERSION).tgz alidns)

.PHONY: clean
clean:
	rm -rf dist
