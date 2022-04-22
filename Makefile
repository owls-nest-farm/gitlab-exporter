CC=go
BIN_NAME=gitlab-exporter

clean:
	rm -rf ${BIN_NAME} migration migration.tar.gz

debug:
	${CC} build
	dlv exec ./${BIN_NAME} -- --file examples/gitlab.json

