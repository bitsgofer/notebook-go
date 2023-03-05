DOCKER=docker
BUILD_DIR=_build
RELEASE_DIR=_release
BIN_NAME=notebook
UID=`id -u`
GID=`id -g`
USER=`whoami`

clean:
	rm -rf ${BUILD_DIR} ${RELEASE_DIR}
	docker rm -f ${DEVTOOL_CONTAINER} || true
.PHONY: clean

# == RUN THIS TO SETUP THE REPO AFTER A FRESH CHECKOUT =========================
setup:
	git config --local include.path ../.gitconfig
	git config --file=.gitconfig core.hooksPath .githooks
.PHONY: setup
# ==============================================================================

# == dev env docker ============================================================
DEVTOOL_IMAGE=bitsgofer.com/notebook-go/devenv:latest
DEVTOOL_CONTAINER=${BIN_NAME}-devenv
DEVTOOL_PACKAGE_PATH=/workspace/src/github.com/bitsgofer/notebook-go

devenv-docker:
	$(DOCKER) build \
		--progress=auto \
		--rm \
		--file dockerfiles/devenv.Dockerfile \
		--build-arg UID=${UID} \
		--build-arg GID=${GID} \
		--build-arg USER=${USER} \
		--tag ${DEVTOOL_IMAGE} \
		.
.PHONY: devenv-docker
# ==============================================================================

# == build (assumed work/automation env is linux/amd64) ========================
${BUILD_DIR}:
	mkdir -p ${BUILD_DIR}
	chown -R ${UID}:${GID} ${BUILD_DIR}

build: ${BUILD_DIR}
	CGO_ENABLED=0 \
	GOOS="linux" \
	GOARCH="amd64" \
	go build \
		-o ${BUILD_DIR}/notebook \
		./cmd/${BIN_NAME}
.PHONY: build

docker-build: ${BUILD_DIR} devenv-docker
	$(DOCKER) create \
		--name=${DEVTOOL_CONTAINER} \
		--mount type=bind,src=$(PWD),dst=${DEVTOOL_PACKAGE_PATH} \
		${DEVTOOL_IMAGE} \
			/bin/bash -c "make build"
	$(DOCKER) start ${DEVTOOL_CONTAINER} > /dev/null
	$(DOCKER) wait ${DEVTOOL_CONTAINER}
	$(DOCKER) cp ${DEVTOOL_CONTAINER}:${DEVTOOL_PACKAGE_PATH}/${BUILD_DIR} ./
	$(DOCKER) rm -f ${DEVTOOL_CONTAINER} > /dev/null
.PHONY: docker-build
# ==============================================================================

# == ensure good commit history (useful after complex rebase work) =============

# 8446449 is the empty initial commit. We could periodically update this
# to the latest tagged release so there's less to check.
TEST_FROM_COMMIT=v0.0.1

verify-commits-can-be-built:
	git rebase -i --exec "make build" ${TEST_FROM_COMMIT}
.PHONY: verify-commits-can-be-built

verify-commits-can-be-tested:
	git rebase -i --exec "make test" ${TEST_FROM_COMMIT}
.PHONY: verify-commits-can-be-tested
# ==============================================================================

# == cross-platform releases ===================================================
${RELEASE_DIR}:
	mkdir -p ${RELEASE_DIR}
	chown -R ${UID}:${GID} ${RELEASE_DIR}

RELEASE_PLATFORMS := linux-amd64 linux-arm64 darwin-amd64 windows-amd64
_PLATFORM_SPLITTED = $(subst /, ,$(subst -, ,$@))
RELEASE_GOOS = $(word 3, $(_PLATFORM_SPLITTED))
RELEASE_GOARCH = $(word 4, $(_PLATFORM_SPLITTED))

# list all binary releases with foreach, e.g:
# ${RELEASE_DIR}/notebook-linux-amd64, ${RELEASE_DIR}/notebook-linux-arm64
all-binaries: ${RELEASE_DIR}
all-binaries: $(foreach platform,$(RELEASE_PLATFORMS),$(RELEASE_DIR)/notebook-$(platform))
.PHONY: all-binaries

# per-platform release
$(RELEASE_DIR)/notebook-%:
	CGO_ENABLED=0 \
	GOOS=${RELEASE_GOOS} \
	GOARCH=${RELEASE_GOARCH} \
	go build \
		-o ${@} \
		./cmd/${BIN_NAME}
.PHONY: ${RELEASE_DIR}/notebook-%

docker-release: ${RELEASE_DIR} devenv-docker
	$(DOCKER) create \
		--name=${DEVTOOL_CONTAINER} \
		--mount type=bind,src=$(PWD),dst=${DEVTOOL_PACKAGE_PATH} \
		${DEVTOOL_IMAGE} \
			/bin/bash -c "go clean -cache && go mod tidy && make all-binaries -j 4"
	$(DOCKER) start ${DEVTOOL_CONTAINER} > /dev/null
	$(DOCKER) wait ${DEVTOOL_CONTAINER}
	$(DOCKER) cp ${DEVTOOL_CONTAINER}:${DEVTOOL_PACKAGE_PATH}/${RELEASE_DIR} ./
	$(DOCKER) rm -f ${DEVTOOL_CONTAINER} > /dev/null
.PHONY: docker-build
# ==============================================================================

# == test ======================================================================
TEST_PKG=

test:
	go test \
		-cover \
		-race \
		-v \
		./${TEST_PKG}/...
.PHONY: test
# ==============================================================================
