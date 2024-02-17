.PHONY: \
	bufbreaking \
	bufbuild \
	bufgen \
	buflint \
	build-image \
	upload-image \
	create-cluster \
	set-kubeconfig \
	delete-cluster \
	kube-apply \
	install-jaeger-operator \
	install-jaeger

bufbreaking:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf breaking --against '.git#branch=main' --error-format=json

bufbuild:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf build

bufgen:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf generate proto

buflint:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf lint

build-image:
	cd ./go/service_b; docker build -t service_b:latest -f ./docker/api/Dockerfile .; cd ../../
	cd ./go/service_c; docker build -t service_c:latest -f ./docker/api/Dockerfile .; cd ../../

upload-image:
	kind load docker-image service_b:latest --name kindcluster
	kind load docker-image service_c:latest --name kindcluster

create-cluster:
	kind create cluster --config kubernetes/kind.yaml --name kindcluster

set-kubeconfig:
	kubectl cluster-info --context kind-kind

delete-cluster:
	kind delete cluster

kube-apply:
	kubectl apply -f kubernetes/service_b.yaml
	kubectl apply -f kubernetes/service_c.yaml

install-jaeger-operator:
	kubectl create namespace observability
	kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.54.0/jaeger-operator.yaml -n observability

install-jaeger:
	kubectl apply -f kubernetes/jaeger.yaml
