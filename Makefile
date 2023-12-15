# database name
DB_NAME ?= postgres

# database type
DB_TYPE ?= postgres

# database username
DB_USER ?= postgres

# database password
DB_PWD ?= bersen

# psql URL
IP=127.0.0.1

PSQLURL ?= $(DB_TYPE)://$(DB_USER):$(DB_PWD)@$(IP):5432/$(DB_NAME)

# sqlc yaml file
SQLC_YAML ?= ./sqlc.yaml

.PHONY : postgresup postgresdown psql createdb teardown_recreate generate

docker-start:
	docker start testpgdb

docker-stop: 
	docker stop testpgdb

postgresup:
	docker run --name testpgdb -v $(PWD):/usr/share/infra -e POSTGRES_PASSWORD=$(DB_PWD) -p 5432:5432 -d $(DB_NAME)

postgresdown:
	docker stop testpgdb  || true && 	docker rm testpgdb || true

migrate-create:
	migrate create -ext sql -dir backend/infrastructure/repositories/pgrepo/migrations -seq create_schema

migrate-up:
	migrate -path backend/infrastructure/repositories/pgrepo/migrations -database "postgresql://postgres:bersen@localhost/postgres?sslmode=disable" -verbose up

migrate-down:
	migrate -path backend/infrastructure/repositories/pgrepo/migrations -database "postgresql://postgres:bersen@localhost/postgres?sslmode=disable" -verbose down

docker-compose-up:
	docker compose --project-directory . -f backend/infrastructure/docker/docker-compose.yml up -d

docker-compose-down:
	docker compose --project-directory . -f backend/infrastructure/docker/docker-compose.yml down

prom-run:
	docker run --name prom \
		-v ./backend/infrastructure/prom/config.yml:/etc/prometheus/prometheus.yml \
		-p 9090:9090 --add-host=host.docker.internal:host-gateway prom/prometheus:v2.29.2

jaeger-run: 
	docker run --name jaeger \
	-p 5775:5775/udp \
	-p 6831:6831/udp \
	-p 6832:6832/udp \
	-p 5778:5778 \
	-p 16686:16686 \
	-p 14268:14268 \
	-p 14250:14250 \
	-p 9411:9411 \
	jaegertracing/all-in-one:1.29.0

psql:
	docker exec -it testpgdb psql $(PSQLURL)

# task to create database without typing it manually
createdb:
	docker exec -it testpgdb psql 

teardown_recreate: postgresdown postgresup
	sleep 5
	$(MAKE) createdb

generate:
	@echo "Generating Go models with sqlc "
	sqlc generate -f $(SQLC_YAML)
