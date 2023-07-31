build/docker:
	@docker build -t mrtc0/seccamp-2023-payments:latest .

build/push: build/docker
	@docker push mrtc0/seccamp-2023-payments:latest

run/local:
	go run main.go
