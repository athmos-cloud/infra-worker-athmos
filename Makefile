

PLUGINS_REPO=git@github.com:athmos-cloud/infra-crossplane-plugin.git
DOCKER_IMAGE_TEST=infra-worker-test
CLUSTER_NAME=infra-worker-test

cluster-test:
	@kind create cluster --name $(CLUSTER_NAME)
	$(MAKE) _crossplane-operator
.PHONY: cluster-test

delete-cluster-test:
	@kind delete cluster --name $(CLUSTER_NAME)
.PHONY: delete-cluster-test

_crossplane-operator:
	@kubectl delete namespace --ignore-not-found=true crossplane-system
	@kubectl create namespace crossplane-system
	@helm repo add crossplane-stable https://charts.crossplane.io/stable &&\
   	helm repo update &&\
   	helm install crossplane --namespace crossplane-system crossplane-stable/crossplane
.PHONY: _crossplane-operator

crossplane-providers:
	@kubectl apply -f config/
.PHONY: crossplane-providers

_build-docker:
	@docker build -t $(DOCKER_IMAGE) -f $(DOCKERFILE) .
.PHONY: _build-docker

test-configs:
	@kubectl apply -f config
.PHONY: test-configs

restart:
	@docker-compose restart infra-worker
.PHONY: restart

del-test-ns:
	@kubectl get namespaces -o jsonpath='{.items[*].metadata.name}' | tr " " "\n" | grep "^test" | xargs -n 1 kubectl delete namespace
.PHONY: del-test-ns
