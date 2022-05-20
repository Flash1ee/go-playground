IMAGE_SERVER_NAME=counter-serv-volume-image
CONTAINER_SERVER_NAME=counter-serv-volume-app

CONTAINER_SERVER_PORT=8088
SYSTEM_SERVER_PORT=81

VOLUME_NAME=counter-volume

volume-create:
	docker volume create $(VOLUME_NAME)

volume-inspect:
	docker volume inspect $(VOLUME_NAME)

volume-rm:
	docker volume rm $(VOLUME_NAME)

build-server:
	docker build  -t $(IMAGE_SERVER_NAME) -f Dockerfile .

run-server:
	docker run -p $(SYSTEM_SERVER_PORT):$(CONTAINER_SERVER_PORT) -v $(CURDIR)/configs/:/app/configs/ --name $(CONTAINER_SERVER_NAME) $(IMAGE_SERVER_NAME) 

stop-server:
	docker stop $(CONTAINER_SERVER_NAME)

rm-server:
	docker rm $(CONTAINER_SERVER_NAME)