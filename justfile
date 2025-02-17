default:
	@air -c .air.toml

restart: sqlc templ css js

sqlc:
	@sqlc generate -f sqlc.yaml

templ:
	@templ generate -path handler/templates

css:
	@node_modules/.bin/tailwindcss \
		--input ./styles/main.scss \
		--output dist/main.css \
		--config ./tailwind.config.js

js:
	@node_modules/.bin/esbuild \
		--bundle \
		--minify \
		--outdir=dist \
		--platform=browser \
		--format=esm \
		./scripts/main.js

# database

migrate:
	GOOSE_MIGRATION_DIR=migrations \
		goose sqlite3 ./data.db up

add_migration name="migration":
	GOOSE_MIGRATION_DIR=migrations \
		goose sqlite3 ./data.db create {{ name }} sql
