genbuf:
	buf lint && buf generate

run:
	go run cmd/server/main.go

srun:
	cd web && npm run build && cd ../
	rm -rf cmd/server/out && cp -R web/out cmd/server
	export CLIENT_HOST=http://localhost:8080 && go run cmd/server/main.go

drun:
	cd web && npm run dev

test:
	go test -race ./...