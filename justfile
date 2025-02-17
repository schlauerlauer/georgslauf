default:
	air -c .air.toml

generate:
	sqlc generate

migrate:
	GOOSE_MIGRATION_DIR=migrations \
		goose sqlite3 ./data.db up

add_migration name="migration":
	GOOSE_MIGRATION_DIR=migrations \
		goose sqlite3 ./data.db create {{ name }} sql
