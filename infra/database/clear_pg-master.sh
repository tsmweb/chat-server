#!/bin/bash

docker-compose stop pg-master
docker-compose rm pg-master
docker volume rm db_pg-master-db
