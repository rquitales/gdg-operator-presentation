# gdg-operator-presentation
Presentation for K8s operator talk. This is a basic introduction to what Kubernetes Operators are, and how we can build a simple serverless platform with it.

## Build and Runn Locally
To run locally, you will need Docker and Node/Yarn installed. Node and yarn is required locally currently, as building the production files through Docker on Mac is super slow. As a hack, this will build the production static files outside of Docker first, for performance - but is bad for dependency management. This needs to be fixed.


You can then start using the commited Makefile:

```sh
make run
```
This will expose the presentation on your localhost on port 8088 which you can visit on your browser (eg: http://localhost:8088)

## Running Pre-Built Container

A pre-built container image is already available on docker.io, which you can use immediately.

```sh
docker run --rm -p 8088:8082 rquitales/presentation:gdg-operators
```
This will expose the presentation on your localhost on port 8088 which you can visit on your browser (eg: http://localhost:8088)

## View slides online
The slides are currently hosted on http://slides.rquitales.com (Currently not on https as I haven't configure SSL for the websocket connections, needs to be fixed!).

I recommend you creating a temporary DigitalOcean token for testing purposes, and destroying the token after using on the public slides.
