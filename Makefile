DOCKER_REPO=nicholasjackson/smi-controller-istio
DOCKER_VERSION=0.1.0
SHELL := /bin/bash

build_docker_setup:
	docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	docker buildx create --name multi
	docker buildx use multi
	docker buildx inspect --bootstrap

build_docker_local: build_docker_setup
	docker buildx build --platform linux/amd64 \
		-t ${DOCKER_REPO}:${DOCKER_VERSION} \
		-f ./Dockerfile \
		. \
		--load
	docker buildx rm multi || true

build_docker_push: build_docker_setup
	docker buildx build --platform linux/arm64,linux/amd64 \
		-t ${DOCKER_REPO}:${DOCKER_VERSION} \
		-f ./Dockerfile \
		. \
		--push
	docker buildx rm multi || true

fetch_certs:
	mkdir -p /tmp/k8s-webhook-server/serving-certs/
	
	kubectl get secret smi-controller-webhook-certificate -n shipyard -o json | \
		jq -r '.data."tls.crt"' | \
		base64 -d > /tmp/k8s-webhook-server/serving-certs/tls.crt
	
	kubectl get secret smi-controller-webhook-certificate -n shipyard -o json | \
		jq -r '.data."tls.key"' | \
		base64 -d > /tmp/k8s-webhook-server/serving-certs/tls.key

run_local: fetch_certs
	go run .

functional_test: fetch_certs
	cd test && go run .