

PLUGINS_REPO=git@github.com:athmos-cloud/infra-crossplane-plugin.git
DOCKER_IMAGE_TEST=infra-worker-test
TMP_INFRA_PLUGINS_DIR=/tmp/athmos/plugins


_plugins:
	@git clone $(PLUGINS_REPO) plugins
	@cd plugins && ./plugin.sh _plugins && cp -r _plugins ../
	@rm -rf plugins
.PHONY: _plugins

_clear-plugins:
	@rm -rf .plugins
.PHONY: _clear-plugins

test-docker: _plugins
	$(MAKE) _build-docker DOCKERFILE=Dockerfile_test DOCKER_IMAGE=infra-worker-test
	@docker run --env CONFIG_FILE_LOCATION="/go/src/app/config.yaml" infra-worker-test
	$(MAKE)  _clear-plugins
.PHONY: test

_plugin-local:
	@rm -rf /tmp/plugins
	@git clone git@github.com:athmos-cloud/infra-crossplane-plugin.git /tmp/plugins
	@mkdir -p $(TMP_INFRA_PLUGINS_DIR)
	@cd /tmp/plugins && ./plugin.sh _plugins && cp -r _plugins/* $(TMP_INFRA_PLUGINS_DIR)
.PHONY: _plugin-local

test: _plugin-local
	@CONFIG_FILE_LOCATION="$(shell pwd)/config.yaml" PLUGINS_LOCATION=$(TMP_INFRA_PLUGINS_DIR) go test -v ./...
.PHONY: test

_build-docker:
	@docker build -t $(DOCKER_IMAGE) -f $(DOCKERFILE) .
.PHONY: _build-docker

restart:
	@docker-compose restart infra-worker
.PHONY: restart