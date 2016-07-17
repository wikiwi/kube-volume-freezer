###  Configuration ###
GO_PACKAGE      ?= github.com/wikiwi/kube-volume-freezer
REPOSITORY      ?= registry.wikiwi.io/vinh/kube-volume-freezer

### Build Tools ###
GO ?= go
GLIDE ?= glide
GIT ?= git
DOCKER ?= docker
TAR ?= tar
ZIP ?= zip
SHA256SUM ?= sha256sum

# Glide Options
GLIDE_OPTS ?=
GLIDE_GLOBAL_OPTS ?=

### Artifact settings for Github Release ###
ARTIFACTS_ARCHIVES ?= kvfctl_linux_amd64.tar.bz2 \
                      kvfctl_darwin_amd64.tar.bz2 \
                      kvfctl_freebsd_amd64.tar.bz2 \
                      kvfctl_windows_amd64.zip

ARTIFACTS_TARGETS := $(ARTIFACTS_ARCHIVES:%=artifacts/%) artifacts/SHA256SUMS

### CI Settings ###
# Set branch with most current HEAD of master e.g. master or origin/master.
# E.g. Gitlab doesn't pull the master branch but fetches it to origin/master.
MASTER_BRANCH ?= master

### Environment ###
HAS_GLIDE := $(shell command -v ${GLIDE};)
HAS_GIT := $(shell command -v ${GIT};)
HAS_GO := $(shell command -v ${GO};)
GOOS := $(shell ${GO} env GOOS)
GOARCH := $(shell ${GO} env GOARCH)
BINARIES := $(notdir $(wildcard cmd/*))

# Load versioning logic.
include versioning.mk

# Docker Image info.
IMAGE := ${REPOSITORY}:${BUILD_REF}

# Show build info.
info:
	@echo $(shell echo $0)
	@echo "Version: ${BUILD_VERSION}"
	@echo "Image:   ${IMAGE}"
	@echo "Tags:    ${TAGS}"

# Creating compressed artifacts from binaries.
artifacts/%.tar.bz2:
	$(eval FILE := bin/$(word 2,$(subst _, ,$*))/$(word 3,$(subst _, ,$*))/$(word 1,$(subst _, ,$*)))
	${MAKE} ${FILE}
	mkdir -p artifacts
	${TAR} -jcvf "$@" ${FILE}
artifacts/%.zip:
	$(eval FILE := bin/$(word 2,$(subst _, ,$*))/$(word 3,$(subst _, ,$*))/$(word 1,$(subst _, ,$*)).exe)
	${MAKE} ${FILE}
	mkdir -p artifacts
	${ZIP} "$@" "${FILE}"

artifacts/SHA256SUMS:
	cd artifacts && ${SHA256SUM} ${ARTIFACTS_ARCHIVES} > $(notdir $@)

.PHONY: build
ifneq (${GOOS}, "windows")
build: ${BINARIES:%=bin/${GOOS}/${GOARCH}/%}
else
build: ${BINARIES:%=bin/${GOOS}/${GOARCH}/%.exe}
endif

.PHONY: build-cross
build-cross: ${BINARIES:%=build-cross-%}
build-cross-%: bin/linux/amd64/% bin/freebsd/amd64/% bin/darwin/amd64/% bin/windows/amd64/%.exe
	$(NOOP)

.PHONY: build-for-docker
build-for-docker: ${BINARIES:%= bin/linux/amd64/%}

# docker-build will build the docker image.
.PHONY: docker-build
docker-build: build-for-docker
	${DOCKER} build --pull -t ${IMAGE} -f Dockerfile.alpine .

# docker-push will push all tags to the repository
.PHONY: docker-push
docker-push: ${TAGS:%=docker-push-%}
docker-push-%:
	${DOCKER} tag ${IMAGE} ${REPOSITORY}:$* && docker push ${REPOSITORY}:$*

# artifacts create
.PHONY: artifacts
artifacts: ${ARTIFACTS_TARGETS}

# clean deletes build artifacts from the project.
.PHONY: clean
clean:
	rm -rf bin artifacts

# test will start the project test suites.
.PHONY: test
test:
	echo Running unit tests
	cd pkg && go test ./...
	echo Running integration tests
	cd test && go test ./...

# bootstrap will install project dependencies.
.PHONY: bootstrap
bootstrap:
ifndef HAS_GO
	$(error You must install Go)
endif
ifndef HAS_GIT
	$(error You must install Git)
endif
ifndef HAS_GLIDE
	${GO} get -u github.com/Masterminds/glide
endif
	${GLIDE} ${GLIDE_GLOBAL_OPTS} install ${GLIDE_OPTS}

include build.mk

