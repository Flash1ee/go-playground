# Simple golang app
IMAGE_SIMPLE_NAME=hello-simple-image
CONTAINER_SIMPLE_NAME=hello-simple-app

EXEC_NAME=main
GOOS_LINUX=linux 
GOARCH_LINUX=amd64

GOOS_MAC=darwin 
GOARCH_MAC=amd64
########################################
# HTTP golang app
IMAGE_SERVER_NAME=counter-server-image
CONTAINER_SERVER_NAME=counter-server-app

CONTAINER_SERVER_PORT=8088
SYSTEM_SERVER_PORT=8088
########################################
# Light golang app
IMAGE_LIGHT_NAME=hello-light-image
CONTAINER_LIGHT_NAME=hello-light-app
###################################################################################################
prepare:
	go env -w GOOS=$(GOOS_LINUX)
	go env -w GOARCH=$(GOARCH_LINUX)

build: clean prepare
	go build -o $(EXEC_NAME) hello.go 
	go env -w GOOS=$(GOOS_MAC)
	go env -w GOARCH=$(GOARCH_MAC)

clean:
	rm -rf $(EXEC_NAME)

build-simple: build
	docker build  -t $(IMAGE_SIMPLE_NAME) -f Dockerfile.simple .

run-simple:
	docker run --name $(CONTAINER_SIMPLE_NAME) $(IMAGE_SIMPLE_NAME) 

logs-simple:
	docker logs -t --details $(CONTAINER_SIMPLE_NAME)

stop-simple:
	docker stop $(CONTAINER_SIMPLE_NAME)

rm-simple:
	docker rm $(CONTAINER_SIMPLE_NAME)
###################################################################################################
build-server:
	docker build  -t $(IMAGE_SERVER_NAME) -f Dockerfile.server .

run-server:
	docker run -p $(SYSTEM_SERVER_PORT):$(CONTAINER_SERVER_PORT) --name $(CONTAINER_SERVER_NAME) -d $(IMAGE_SERVER_NAME) 

stop-server:
	docker stop $(CONTAINER_SERVER_NAME)

rm-server:
	docker rm $(CONTAINER_SERVER_NAME)

###################################################################################################
build-light:
	docker build  -t $(IMAGE_LIGHT_NAME) -f Dockerfile.light .

run-light:
	docker run --name $(CONTAINER_LIGHT_NAME) $(IMAGE_LIGHT_NAME) 

stop-light:
	docker stop $(CONTAINER_LIGHT_NAME)

rm-light:
	docker rm -f $(CONTAINER_LIGHT_NAME)