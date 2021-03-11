install_swagger:
	which swagger || GO11MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: install_swagger
	GO11MODULE=off swagger generate spec -o ./swagger.yaml --scan-models