default:
	@air -c .air.toml

generate:
	@sqlc generate

js:
	@node_modules/.bin/esbuild \
		--bundle \
		--minify \
		--outdir=dist \
		--platform=browser \
		--format=esm \
		./scripts/main.js

css:
	@node_modules/.bin/tailwindcss \
		--input ./styles/main.scss \
		--output dist/main.css \
		--config ./tailwind.config.js

# database

migrate:
	GOOSE_MIGRATION_DIR=migrations \
		goose sqlite3 ./data.db up

add_migration name="migration":
	GOOSE_MIGRATION_DIR=migrations \
		goose sqlite3 ./data.db create {{ name }} sql
