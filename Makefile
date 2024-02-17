.PHONY: \
	bufbreaking \
	bufbuild \
	bufgen \
	buflint \
	create-cluster \
	set-kubeconfig \
	delete-cluster

bufbreaking:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf breaking --against '.git#branch=main' --error-format=json

bufbuild:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf build

bufgen:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf generate proto

buflint:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf lint

create-cluster:
	kind create cluster --config kubernetes/kind.yaml --name kindcluster

set-kubeconfig:
	kubectl cluster-info --context kind-kind

delete-cluster:
	kind delete cluster
