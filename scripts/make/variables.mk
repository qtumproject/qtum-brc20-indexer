############## REQUIRED ##############
#
# SERVICE_NAME	:= <service_name>
#
############## REQUIRED ##############

# !!!should be modified when generate new repo from this!!!

VAR_REPO_PATH 			:= .
VAR_OWNER_NAME			:= foxwallet
VAR_PROJECT_IDENTIFIER	:= fox-ordinal
VAR_REPO_URL			:= web3-container-registry.cn-hongkong.cr.aliyuncs.com
TARGET_NAMESPACE        := minikube-dev

# generic variables

VAR_RELATIVE_PATH		:= .
VAR_APP_RELATIVE_PATH   := $(shell a=`basename $$PWD` && cd .. && b=`basename $$PWD` && echo $$b/$$a)


VAR_GIT_REVISION 		:= $(shell date '+%Y%m%d%H%M%S')-$(shell git rev-parse --verify --short HEAD 2>/dev/null)
VAR_GIT_ROOT_DIR    	:= $(shell git rev-parse --show-toplevel)
VAR_COMPILE_TIME 		:= $(shell git log -1 --format="%ad" --date=short)
VAR_ARGS_GO_BUILD 		:= CGO_ENABLED=0 GOOS=linux GOARCH=amd64
