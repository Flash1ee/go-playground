build:
	go build -o server.out -v ./cmd/server

build-docker-server:
	docker build --network host -f ./docker/Dockerfile . --tag clickstream-app

run-compose:
	docker-compose --env-file ./config/credintials.env up
pause:
	docker-compose down

stop: # с удалением volumes!
	docket-compose stop