.PHONY: test

install-dependencies:
	go get

build-static: install-dependencies
	go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o exchangeratesgraphql

build: install-dependencies
	go build -o exchangeratesgraphql
	chmod +x exchangeratesgraphql

run:
	./exchangeratesgraphql -serve

scrape:
	./exchangeratesgraphql -scrape

test:
	go test

coverage:
	goveralls -repotoken ${COVERALLS_TOKEN}

docker-build:
	docker build -t jonatasbaldin/exchangeratesgraphql:latest -t jonatasbaldin/exchangeratesgraphql:${TRAVIS_TAG} .

docker-push:
	echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
	docker push jonatasbaldin/exchangeratesgraphql:latest
	docker push jonatasbaldin/exchangeratesgraphql:${TRAVIS_TAG}