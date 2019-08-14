name := neosapience/cata
tag := v0.0.1

.PHONY: build
build:
	docker build -t ${name}:dev .

.PHONY: dist
dist:
	docker build -t ${name}:${tag} . -f Dockerfile.dist

.PHONY: up
up:
	docker-compose up -d

.PHONY: up-dist
up-dist:
	docker-compose -f docker-compose.yml -f docker-compose.dist.yml up -d

.PHONY: up-test
up-test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up -d


.PHONY: logs
logs:
	docker-compose logs -f app 

.PHONY: test
test: build up-test logs

.PHONY: sh
sh:
	docker run --rm -it ${name}:dev bash

.PHONY: ls
ls:
	@docker images ${name}

.PHONY: ps
ps:
	@docker-compose ps

.PHONY: ps
down:
	@docker-compose down

.PHONY: push
push:
	@docker push ${name}:${tag}



