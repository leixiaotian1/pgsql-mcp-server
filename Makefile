APP_NAME=sql-mcp-server
APP_PORT=8088
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_CONTAINER=$(APP_NAME)-container
DOCKER_NETWORK=my-app-network
ENV_FILE=.env

.PHONY: build run stop logs clean

## 构建 Docker 镜像
build:
	docker build -t $(DOCKER_IMAGE) .

## 运行容器（加载 .env 文件，映射端口）
run: stop
	docker run -d \
	    --network=$(DOCKER_NETWORK) \
		--name $(DOCKER_CONTAINER) \
		--env-file $(ENV_FILE) \
		-p $(APP_PORT):$(APP_PORT) \
		$(DOCKER_IMAGE)

## 停止并删除容器
stop:
	-@docker stop $(DOCKER_CONTAINER) 2>null || exit 0
	-@docker rm $(DOCKER_CONTAINER) 2>null || exit 0


## 查看容器日志
logs:
	docker logs -f $(DOCKER_CONTAINER)

## 清理镜像和容器
clean: stop
	-@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
