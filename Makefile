.PHONY: apidoc proto

PROTO_DIR = "./service/id_generator/proto"

apidoc:
	swagger generate spec -o ./doc/swagger.json && swagger serve ./doc/swagger.json

proto:
	cd ${PROTO_DIR} && protoc --proto_path=. --go_out=${GOPATH}/src --micro_out=${GOPATH}/src ./*.proto
