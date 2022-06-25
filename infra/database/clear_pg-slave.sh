#!/bin/bash

docker-compose stop pg-slave
docker-compose rm pg-slave
docker volume rm db_pg-slave-db
