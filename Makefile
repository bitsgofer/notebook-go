BUILD_DIR=_build
BIN_NAME=notebook
GO_BUILD_ENV=CGO_ENABLED=0 GOOS="linux" GOARCH="amd64"

clean:
	rm -rf ${BUILD_DIR}
.PHONY: clean

${BUILD_DIR}:
	mkdir -p ${BUILD_DIR}

build: ${BUILD_DIR}
	${GO_BUILD_ENV} \
	go build \
		-o ${BUILD_DIR}/notebook \
		./cmd/${BIN_NAME}
.PHONY: build
