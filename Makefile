SHELL = /bin/bash

run.dev:
	@docker compose up -d --build --force-recreate

docker.build:
	@docker build -t mgenteluci/rinha2024q1 .

docker.push:
	@docker push mgenteluci/rinha2024q1

run.prod:
	@docker compose -f docker-compose-prod.yml up -d --force-recreate
