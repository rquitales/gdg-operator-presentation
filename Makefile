run: build
	docker run --rm -p 8088:8082 rquitales/presentation:gdg-operators

run-only:
	docker run --rm -p 8088:8082 rquitales/presentation:gdg-operators

build:
	yarn --cwd ./slides-ts/ build
	docker build -t rquitales/presentation:gdg-operators .
