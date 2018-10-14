# If the USE_SUDO_FOR_DOCKER env var is set, prefix docker commands with 'sudo'
ifdef USE_SUDO_FOR_DOCKER
	SUDO_CMD = sudo
endif

IMAGE ?= asia.gcr.io/k8s-project-199813/osbapi-baas
TAG ?= $(shell git describe --tags --always)
PULL ?= IfNotPresent

build: ## Builds the app
	go build -i github.com/cclin81922/osbapi-baas/cmd/osbapibaas

linux: ## Builds a Linux executable
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -o osbapibaas-linux --ldflags="-s" github.com/cclin81922/osbapi-baas/cmd/osbapibaas

image: linux ## Builds a Linux based image
	cp osbapibaas-linux image/osbapibaas
	$(SUDO_CMD) docker build image/ -t "$(IMAGE):$(TAG)"

clean: ## Cleans up build artifacts
	rm -f osbapibaas
	rm -f osbapibaas-linux
	rm -f image/osbapibaas

push: image ## Pushes the image to dockerhub, REQUIRES SPECIAL PERMISSION
	$(SUDO_CMD) docker push "$(IMAGE):$(TAG)"

deploy-baas: image ## Deploys baas with helm
	helm upgrade --install baas-skeleton --namespace baas-skeleton \
	charts/osbapibaas \
	--set image.repository="$(IMAGE)",image.tag="$(TAG)",image.pullPolicy="$(PULL)"

remove-baas: ## Removes baas with helm
	helm delete --purge baas-skeleton

create-ns: ## Creates a namespace
	kubectl create ns baas-skeleton

remove-ns: ## Removes a namespace
	kubectl delete ns baas-skeleton

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: build linux image clean push deploy-baas remove-baas create-ns remove-ns help
