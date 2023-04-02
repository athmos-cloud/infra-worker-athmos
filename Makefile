

PLUGINS_REPO=git@github.com:athmos-cloud/infra-crossplane-plugin.git
DOCKER_IMAGE_TEST=infra-worker-test

_plugins:
	@git clone $(PLUGINS_REPO) plugins
	@cd plugins && ./plugin.sh _plugins && cp -r _plugins ../
	@rm -rf plugins
.PHONY: _plugins

_clear-plugins:
	@rm -rf .plugins
.PHONY: _clear-plugins

test: _plugins
	$(MAKE) _build-docker DOCKERFILE=Dockerfile_test DOCKER_IMAGE=infra-worker-test
	@docker run --env CONFIG_FILE_LOCATION="/go/src/app/config.yaml" infra-worker-test
	$(MAKE)  _clear-plugins
.PHONY: test

_build-docker:
	@docker build -t $(DOCKER_IMAGE) -f $(DOCKERFILE) .
.PHONY: _build-docker

