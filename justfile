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

migration_status:
	@atlas migrate status \
		--dir "file://migrations" \
		--url "sqlite://georgslauf.db"

migration_generate migration_name="migration":
	@atlas migrate diff {{ migration_name }} \
		--dir "file://migrations" \
		--to "file://schema.hcl" \
		--dev-url "sqlite://file?mode=memory"

migration_apply:
	@atlas migrate apply \
		--dir "file://migrations" \
		--url "sqlite://georgslauf.db"
