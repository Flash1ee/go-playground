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

# После первого старта контейнера
init-mongo-sh-config:
	docker-compose exec mongo-cfg-server01 sh -c "mongo < ./scripts/init-configserver.js"
	sleep 5
	docker-compose exec mongodb-shard01 sh -c "mongo < ./scripts/init-shard01.js"

# После запуска init-mongo-sh-config
init-router:
	docker-compose exec mongo-router01 sh -c "mongo < ./scripts/init-router.js"

enable-sharding:
	docker-compose exec mongo-router01 sh ./scripts/enableSharding.sh

migrate-data:
	docker-compose exec mongo-router01 sh -c "mongo < ./scripts/db-data-init.js"

verify-router:
	docker-compose exec mongo-router01 sh -c "echo 'sh.status()' | mongo"

verify-sh-rs:
	docker-compose exec mongodb-shard01 sh -c "echo 'rs.status()' | mongo"

check-db-status:
	docker-compose exec mongo-router01 sh -c "echo 'use clickstream\ndb.stats()\ndb.events.getShardDistribution()' | mongo"

