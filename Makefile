BIN_NAME=gitlab-exporter

clean:
	rm -rf ${BIN_NAME} migration

debug:
	go build
	dlv exec ./${BIN_NAME} -- --file examples/gitlab.json

