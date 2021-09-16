NAME=go-grpc-example
Q=$(if $V,,@)

all: bin/server-acme bin/server-cert bin/client

bin/server-acme: server-acme/main.go
	$Q go build -o bin/server-acme server-acme/main.go

bin/server-cert: server-cert/main.go
	$Q go build -o bin/server-cert server-cert/main.go

bin/client: client/main.go
	$Q go build -o bin/client client/main.go

clean:
	$Q rm -rf bin

docker:
	$Q docker build -t ${NAME} .

docker-dev: docker
	$Q docker tag ${NAME}:latest localhost:5000/${NAME}:latest
	$Q docker push localhost:5000/${NAME}:latest
