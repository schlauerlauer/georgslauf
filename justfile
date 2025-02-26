default:
	@air -c .air.toml

restart: generate css js

generate:
	@go generate

css:
	@node_modules/.bin/tailwindcss \
		--input ./dist/main.scss \
		--output resources/main.css \
		--config ./tailwind.config.js

js:
	@node_modules/.bin/esbuild \
		--bundle \
		--minify \
		--outdir=resources \
		--platform=browser \
		--format=esm \
		./dist/main.js

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
