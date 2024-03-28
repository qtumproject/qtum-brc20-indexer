# input variables

REPO_PATH 		:= ${VAR_REPO_PATH}
GIT_ROOT_DIR    := ${VAR_GIT_ROOT_DIR}

# calculated variables

ROOT_DIR 		:= ${GIT_ROOT_DIR}/${REPO_PATH}

# commands

gen_grpc:
	@echo ">> generating grpc"
	$(ROOT_DIR)/scripts/shell/grpc_gen.sh
	@echo ">> generate grpc end"

gen_wire:
	@echo ">> generate wire start"
	cd cmd/server && wire
	@echo ">> generate wire end"

gen: gen_grpc gen_wire
