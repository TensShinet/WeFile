apidoc:
	swagger generate spec -o ./doc/swagger.json && swagger serve ./doc/swagger.json