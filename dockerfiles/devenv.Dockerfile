FROM golang:1.20 AS devenv

ENV GOPATH=/workspace

# Tell git to ignore ownership of the ./.git directory, since it will be mangled
# when we bind-mount the repo to into a container.
# Without this, git complains, which in turn make the go compiler complains
# (as go 1.18+ start reading VCS/git info to put into the binaries).
# REF: https://github.com/golang/go/issues/53532
#
# => This is turned on right now.
# If we manage to run and bind-mount the whole repo using the same uid:gid,
# we might not need to set the safe.directory
RUN git config --global --add safe.directory ${GOPATH}/src/github.com/bitsgofer/notebook-go
#===============================================================================

FROM devenv
ARG PACKAGE_PATH=${GOPATH}/src/github.com/bitsgofer/notebook-go
RUN mkdir -p ${PACKAGE_PATH}
COPY . ${PACKAGE_PATH}
WORKDIR ${PACKAGE_PATH}
