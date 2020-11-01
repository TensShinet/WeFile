.PHONY: apidoc proto

PROTO_DIR = "./service/id_generator/proto"
PWD = $(shell pwd)

apidoc:
	swagger generate spec -o ./doc/swagger.json && swagger serve ./doc/swagger.json

proto:
	cd ${PROTO_DIR} && protoc --proto_path=. --go_out=${GOPATH}/src --micro_out=${GOPATH}/src ./*.proto

build_auth:
	cd $(PWD)/service/auth && go build -o auth_service main.go

build_db:
	cd $(PWD)/service/db && go build -o db_service main.go

build_idg:
	cd $(PWD)/service/id_generator && go build -o idg_service main.go

build_file:
	cd $(PWD)/service/file && go build -o file_service main.go

build_base:
	cd $(PWD)/service/base && go build -o base_service main.go

run_auth:
	cd $(PWD)/service/auth && ./auth_service &

run_db:
	cd $(PWD)/service/db && ./db_service &

run_idg:
	cd $(PWD)/service/id_generator && ./idg_service &

run_file:
	cd $(PWD)/service/file && ./file_service &

run_base:
	cd $(PWD)/service/base && ./base_service &

build_all: build_auth build_db build_idg build_file build_base

run_all: run_auth run_db run_idg run_file run_base