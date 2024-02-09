SHELL = /bin/bash

run.local:
	@go build cmd/main.go && ./main

run.dev:
	@docker compose up -d

docker.build:
	@docker build -t mgenteluci/rinha2024q1 .

docker.push:
	@docker push mgenteluci/rinha2024q1
