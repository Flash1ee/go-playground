build:
	go build -o server.out -v ./cmd/server

build-docker-server:
	docker build --network host -f ./docker/Dockerfile . --tag clickstream-app

run-compose:
	#docker-compose --env-file ./config/credintials.env up
	./run.sh
pause:
	docker-compose stop

stop: # с удалением volumes!
	docket-compose down

mongo-init:
	docker exec mongodb /migrations/mongo-init.sh

mongo-rs:
	docker exec mongodb /scripts/rs-init.sh
