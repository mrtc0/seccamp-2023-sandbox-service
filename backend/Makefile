build/docker:
	@docker build -t mrtc0/seccamp-2023-backend:latest .

build/push: build/docker
	@docker push mrtc0/seccamp-2023-backend:latest

run/local:
	export PAYMENTS_API_ADDR=http://localhost:7000/ ITEMS_API_ADDR=http://localhost:9000/ && go run main.go
