# docker build：make docker-build TAG=v1.0.0
# docker push：make docker-push TAG=v1.0.0
# docker rm：make docker-rm TAG=v1.0.0
# docker rmi：make docker-rmi TAG=v1.0.0

# 应用名称
APP_NAME=sync_eth
# get version from tag
TAG ?= $(shell git describe --tags --abbrev=0 2>/dev/null)

DOCKER_ACC ?= chain5j
DOCKER_REPO ?= $(APP_NAME)

####################################
docker-build:
	# show progress bar
	rm -rf sync_eth
	GOOS=linux CGO_ENABLED=0 go build -o sync_eth
	export DOCKER_BUILDKIT=1 && docker build  --no-cache -t $(DOCKER_ACC)/$(DOCKER_REPO):$(TAG) ./
	rm -rf sync_eth
docker-rm:
	docker rm -f $(APP_NAME)
docker-rmi:
	docker rmi -f $(DOCKER_ACC)/$(DOCKER_REPO):$(TAG)
docker-rm-err:
	docker rm `docker ps -a | grep Exited | awk '{print $$1}'`
docker-rmi-err:
	docker rmi $$(docker images -q -f dangling=true)
docker-run:
	docker run -it --name $(APP_NAME) $(DOCKER_ACC)/$(DOCKER_REPO):$(TAG)
docker-start:
	docker exec -it $(APP_NAME) /bin/sh start
docker-logs:
	docker logs -f $(APP_NAME)
docker-push:
	docker push $(DOCKER_ACC)/$(DOCKER_REPO):$(TAG)
