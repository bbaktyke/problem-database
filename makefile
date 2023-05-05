run:
	docker-compose up --biuld
























































#for installing and lauching docker image of posstgres and adding migration files

# postgres_run:
# 	docker run --name=problem-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 --rm postgres
# exec:
# 	docker exec -it problem-db /bin/bash
# migrate_install:
# 	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# schema:
# 	migrate create -ext sql -dir .\schema -seq init
# migrate_up:
# 	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up
# migrate_down:
# 	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' down

# #for creating docker container of application

# create:
# 	docker build -t beka .
# run:
# 	docker run -p 8080:8080 --rm --name bbaktyke beka 
# stop:
# 	docker stop bbaktyke
# start:
# 	docker start  bbaktyke
# prune:
# 	docker container prune

	

# redis_run:
# 	docker run --name my-redis -p 6379:6379 -v /my/data/dir:/data -d redis redis-server --appendonly yes

# redis_exec:
# 	docker exec -it <container_name> apt-get update
# 	docker exec -it <container_name> apt-get install redis-tools
# 	docker exec -it <container_name> redis-cli


# rabbitMQ:
# 	docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management

# elaasticsearch:
# 	docker run -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" --memory="2g" --cpus="2" docker.elastic.co/elasticsearch/elasticsearch:7.12.1

