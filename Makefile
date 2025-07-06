build:
	cd frontend && pnpm build
	ENV=prod go build -buildvcs=false -o ./bin/whats-valid ./main.go

dev:
	cd frontend && pnpm dev & air && fg

start:
	ENV=prod ./bin/whats-valid
