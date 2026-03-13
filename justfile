default:
	@air -c .air.toml

generate:
	@go generate

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

# build

container_update: build_arm_container save_arm_image copy_arm_image

build_arm_container:
	@podman build \
		-f arm64.containerfile \
		--platform linux/arm64 \
		-t localhost/georgslauf:latest \
		.

save_arm_image:
	@rm ./georgslauf.tar.gz
	@podman save localhost/georgslauf:latest \
		-o georgslauf.tar.gz

copy_arm_image:
	@echo "confirm"
	@scp ./georgslauf.tar.gz georgslauf:~/containers/georgslauf/

migrations_save:
	@tar -czf migrations.tar.gz migrations

migrations_copy:
	@echo "confirm"
	@scp migrations.tar.gz georgslauf:~/containers/georgslauf/
