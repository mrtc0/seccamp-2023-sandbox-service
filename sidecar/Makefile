build/docker:
	@docker build -t mrtc0/seccamp-2023-sidecar-proxy:latest .
	@docker build -t mrtc0/seccamp-2023-sidecar-init:latest -f Dockerfile.init-container

build/push: build/docker
	@docker push mrtc0/seccamp-2023-sidecar-proxy:latest
	@docker push mrtc0/seccamp-2023-sidecar-init:latest

run/local:
	go run main.go

