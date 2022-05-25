CC=go
BIN_NAME=gitlab-exporter

build:
	${CC} build

clean:
	rm -rf ${BIN_NAME} /tmp/migration migration.tar.gz

debug: build
	dlv exec ./${BIN_NAME} -- --file examples/simple.json

