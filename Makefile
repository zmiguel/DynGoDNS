BINARY_NAME=DynGoDNS

build: cloudflare.so
	go build -o build/${BINARY_NAME} cmd/DynGoDNS/main.go

run:
	./build/${BINARY_NAME}

cloudflare.so:
	go build -buildmode=plugin -o build/plugins/cloudflare.so ./cmd/plugins/cloudflare/cloudflare.go

bnr: build run

clean:
	go clean
	rm ${BINARY_NAME}
	rm plugins/*