include ../../../scripts/make/variables.mk
include ../../../scripts/make/grpc_gen.mk
include ../../../scripts/make/golang.mk
include ../../../scripts/make/skaffold.mk

SERVICE_NAME := ordinal-indexer-task
GRPC_PORT := 50051

manifest-test-local: export CONFIG_FILE_NAME=config.yaml
manifest-test-local: export IMAGE=web3-container-registry.cn-hongkong.cr.aliyuncs.com/${OWNER_NAME}/foxwallet-server-test/${SERVICE_NAME}:${GIT_REVISION}
manifest-test-local: manifest

manifest-test: export CI=true
manifest-test: manifest-test-local

docker-local-prod: export CONFIG_FILE_NAME=config.yaml
docker-local-prod: manifest


test1:
	@echo ">> dry-run the kubernetes specs"
	kustomize build ${KUSTOMIZE_DIR}/overlays/foxwallet | envsubst >| ${BUILD_DIR}/k8s.yaml