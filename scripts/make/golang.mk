# input variables

REPO_PATH 			:= ${VAR_REPO_PATH}
PROJECT_IDENTIFIER	:= ${VAR_PROJECT_IDENTIFIER}
RELATIVE_PATH		:= ${VAR_RELATIVE_PATH}
GIT_REVISION 		:= ${VAR_GIT_REVISION}
GIT_ROOT_DIR    	:= ${VAR_GIT_ROOT_DIR}
COMPILE_TIME 		:= ${VAR_COMPILE_TIME}
ARGS_GO_BUILD 		:= ${VAR_ARGS_GO_BUILD}
APP_RELATIVE_PATH := ${VAR_APP_RELATIVE_PATH}
# calculated variables

ROOT_DIR 			:= ${GIT_ROOT_DIR}/${REPO_PATH}
BUILD_DIR 			:= ${RELATIVE_PATH}/output
MAIN_ENTRY 			:= ${RELATIVE_PATH}/cmd/server

# required environment variables

export REPO_URL=${VAR_REPO_URL}
export OWNER_NAME=${VAR_OWNER_NAME}

# commands

run: gen
	@echo ">> run service"
	@go run ${MAIN_ENTRY}

init:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/protobuf/cmd/protoc-gen-go

	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

	go get -u github.com/google/wire/cmd/wire
	go install github.com/google/wire/cmd/wire

	go get -u github.com/mohuishou/protoc-gen-go-gin@latest
	go install github.com/mohuishou/protoc-gen-go-gin@latest

	go get -u github.com/favadi/protoc-go-inject-tag@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

	go get github.com/grpc-ecosystem/grpc-health-probe
	go install github.com/grpc-ecosystem/grpc-health-probe

build: gen
	@echo ">> build the binary package"
	go build -o ${BUILD_DIR}/app ${MAIN_ENTRY}

# set the os and arch for the binary package
build-linux: gen
	@echo ">> build the binary package for linux(docker)"
	mkdir -p output && \
	mkdir -p output/bin output/configs && \
    cp -r configs output/ && \
	$(ARGS_GO_BUILD) go build -a -v -ldflags "-s -X main.GitSHA=$(VERSION) -X main.BuildTime=$(WHEN)" -o ${BUILD_DIR}/bin/server ${MAIN_ENTRY}

build-docker: init build-linux

docker: export BUILD_CONTEXT=${GIT_ROOT_DIR}
docker: export APP_RELATIVE_PATH=${VAR_APP_RELATIVE_PATH}
docker: docker-build

docker-build:
	@echo ">> build the docker image"
	$(ROOT_DIR)/scripts/shell/docker_build.sh

# cleanup the grpc stubs and binary
.PHONY:clean
clean:
	@rm -rf ${BUILD_DIR}
	@rm -f ${RELATIVE_PATH}/*.pb.go
