run:
	go run main.go
install-swagger:
	which swagger || GO11MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: install-swagger
	GO11MODULE=off swagger generate spec -o ./swagger.yaml --scan-models