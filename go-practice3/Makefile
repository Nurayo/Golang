DB=postgres://nuray:12345@localhost:5432/practice3?sslmode=disable
MIG=internal/db/migrations

migrate-up:
	migrate -path $(MIG) -database "$(DB)" up

migrate-down:
	migrate -path $(MIG) -database "$(DB)" down 1

migrate-version:
	migrate -path $(MIG) -database "$(DB)" version

verify:
	go run ./cmd/verify
