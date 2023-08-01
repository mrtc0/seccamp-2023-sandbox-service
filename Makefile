build/docker:
	cd api && make build/docker
	cd backend && make build/docker
	cd payments && make build/docker
	cd sidecar && make build/docker
	cd worker && make build/docker

push/docker:
	cd api && make build/push
	cd backend && make build/push
	cd payments && make build/push
	cd sidecar && make build/push
	cd worker && make build/push
