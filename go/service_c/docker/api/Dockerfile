FROM golang:1.22.0-bookworm as builder
ENV CGO_ENABLED 0
WORKDIR /workspace
COPY . .
RUN apt-get update && apt-get install --no-install-recommends -y make
RUN make build

FROM gcr.io/distroless/static-debian12:latest as api
COPY --from=builder --chown=nonroot:nonroot /workspace/dist/api /api
USER nonroot
CMD ["/api"]
