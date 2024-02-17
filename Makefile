.PHONY: \
	bufbreaking \
	bufbuild \
	bufgen \
	buflint

bufbreaking:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf breaking --against '.git#branch=main' --error-format=json

bufbuild:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf build

bufgen:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf generate proto

buflint:
	docker run --volume "$$(pwd):/workspace" --workdir /workspace bufbuild/buf lint
