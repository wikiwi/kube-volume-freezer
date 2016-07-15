##  Configuration ##
GO_PACKAGE      ?= github.com/wikiwi/kube-volume-freezer
REGISTRY        ?= registry.wikiwi.io
IMAGE_PREFIX    ?= vinh
SHORT_NAME      ?= kube-volume-freezer

# Glide Options
GLIDE_OPTS ?=

# Set branch with most current HEAD of master e.g. master or origin/master.
MASTER_BRANCH ?= master

# List of binaries to be build.
BINARIES ?= kvf-master kvf-minion kvfctl

# VERSION contains the project verison e.g. 0.1.0-alpha.1
VERSION := $(shell grep -E -o "[0-9]+\.[0-9]+\.[0-9]+[^\"]*" pkg/version/version.go)
# VERSION_MINOR contains the project version up to the minor value e.g. v0.1
VERSION_MINOR := $(shell echo ${VERSION} | grep -E -o "[0-9]+\.[0-9]+")
# VERSION_STAGE contains the project version stage e.g. alpha
VERSION_STAGE := $(shell echo ${VERSION} | grep -E -o "(pre-alpha|alpha|beta|rc)")

# Extract git Information of current commit.
GIT_SHA := $(shell git rev-parse HEAD)
GIT_SHA_SHORT := $(shell git rev-parse --short HEAD)
GIT_SHA_MASTER := $(shell git rev-parse ${MASTER_BRANCH})
GIT_TAG := $(shell git tag -l --contains HEAD | head -n1)
GIT_BRANCH := $(shell git branch | grep -E '^* ' | cut -c3- )
IS_DIRTY := $(shell git status --porcelain)

ifndef IS_DIRTY
  ifeq ($(GIT_SHA),$(GIT_SHA_MASTER))
    IS_CANARY       := true
    ifeq ($(GIT_TAG),$(VERSION))
      IS_RELEASE      := true
      ifeq ($(LATEST),$(VERSION_MINOR))
        IS_LATEST := true
      endif
    endif
  endif
endif

# Set build referece.
ifdef IS_DIRTY
  BUILD_REF      := $(GIT_SHA_SHORT)-dev
else
  BUILD_REF      := $(GIT_SHA_SHORT)
endif

# BUILD_VERSION will be compiled into the projects binaries.
ifdef IS_RELEASE
  BUILD_VERSION    ?= ${VERSION}
else
  BUILD_VERSION    ?= ${VERSION}+${BUILD_REF}
endif

# Docker Image settings.
REPOSITORY := ${REGISTRY}/${IMAGE_PREFIX}/${SHORT_NAME}
IMAGE := ${REPOSITORY}:${BUILD_REF}
IMAGE_FILE := ${SHORT_NAME}.tar.gz

# Set image tags.
TAGS :=
ifeq ($(IS_CANARY),true)
  TAGS := canary ${TAGS}
endif
ifdef IS_RELEASE
  TAGS := ${VERSION} ${TAGS}
  ifdef VERSION_STAGE
    TAGS := ${VERSION_MINOR}-${VERSION_STAGE} ${TAGS}
  else
    TAGS := ${VERSION_MINOR} ${TAGS}
  endif
  ifeq ($(IS_LATEST),true)
    TAGS := latest ${TAGS}
  endif
endif

# Show build info.
info:
	@echo "Version: ${BUILD_VERSION}"
	@echo "Image:   ${IMAGE}"
	@echo "Tags:    ${TAGS}"

# build will compile the binaries.
.PHONY: build
BUILD_CMD = GOBIN=$(CURDIR)/bin go install -ldflags "-X ${GO_PACKAGE}/pkg/version.Version=${BUILD_VERSION}" $(GO_PACKAGE)/cmd/${BINARY}
build: clean
	$(foreach BINARY,$(BINARIES),($(BUILD_CMD)) || exit $$?;)

# docker-build will build the docker image.
.PHONY: docker-build
docker-build:
	docker build --pull -t ${IMAGE} .

# docker-save will save the built image.
.PHONY: docker-save
docker-save:
	mkdir -p images
	docker save ${IMAGE} | gzip > images/${IMAGE_FILE}

# docker-load will load the saved docker image.
.PHONY: docker-load
docker-load:
	mkdir -p images
	gzip -cd images/${IMAGE_FILE} | docker load

# docker-test will run tests inside docker container.
.PHONY: docker-test
docker-test:
	docker run --rm ${IMAGE} make test

# docker-push will push the previously build image.
.PHONY: docker-push
PUSH_CMD = docker tag ${IMAGE} ${REPOSITORY}:${TAG} && docker push ${REPOSITORY}:${TAG}
docker-push:
	$(foreach TAG,$(TAGS),($(PUSH_CMD)) || exit $$?;)

.PHONY: docker-push-%
docker-push-%:
	$(eval TAG := $*)
	$(PUSH_CMD)

# clean deletes build artifacts from the project.
.PHONY: clean
clean:
	rm -rf bin

# test will start the project test suites.
.PHONY: test
test:
	echo Running unit tests
	cd pkg && go test ./...
	echo Running integration tests
	cd test && go test ./...

# bootstrap will install project dependencies.
.PHONY: bootstrap
HAS_GLIDE := $(shell command -v glide;)
HAS_GIT := $(shell command -v git;)
HAS_GO := $(shell command -v go;)
bootstrap:
ifndef HAS_GO
	$(error You must install Go)
endif
ifndef HAS_GIT
	$(error You must install Git)
endif
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install ${GLIDE_OPTS}
