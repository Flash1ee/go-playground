IMAGE_NAME=postgres
CONTAINER=my_container_db
RESTART=on-failure:2
RESTART=always
LOG_FILE=db.log

docker-pull:
	docker pull $(IMAGE_NAME)

docker-run:
	docker run --name $(CONTAINER) -d -e POSTGRES_PASSWORD=qwerty $(IMAGE_NAME)

docker-run-with-restart:
	docker run --name $(CONTAINER) -d --restart $(RESTART) $(IMAGE_NAME)

docker-update-restart:
	docker update --restart no $(CONTAINER)

docker-logs:
	docker logs  -t --details $(CONTAINER)

docker-logs-runtime:
	docker logs  -t --details $(CONTAINER) --follow
	
docker-logs-in-file:
	docker logs -f $(CONTAINER) > $(LOG_FILE) 2>&1 &

docker-stop:
	docker stop $(CONTAINER)

docker-start:
	docker start $(CONTAINER)

docker-rm:
	docker rm $(CONTAINER)