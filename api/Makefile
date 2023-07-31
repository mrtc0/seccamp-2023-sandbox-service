build/docker:
	@docker build -t mrtc0/seccamp-2023-api:latest .

build/push: build/docker
	@docker push mrtc0/seccamp-2023-api:latest

run/local:
	go run main.go
