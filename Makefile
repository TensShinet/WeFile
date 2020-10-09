.PHONY: apidoc proto

PROTOL_DIR = "./service/id_generator/proto"

apidoc:
	swagger generate spec -o ./doc/swagger.json && swagger serve ./doc/swagger.json

proto:
	cd ${PROTOL_DIR} && protoc --proto_path=. --go_out=. --micro_out=. ./*.proto
