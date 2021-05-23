RUNTIME=podman
DB_TAG=10.5.1
WEB_TAG=7.4.3-apache-buster

prepare_ubuntu:
ifeq ($(RUNTIME),podman)
	sudo usermod --add-subuids 100000-150000 --add-subgids 100000-150000 $(USER)
	podman system migrate
else
	$(error not necessary for $(RUNTIME))
endif

pull:
	$(info pulling mariadb image)
	$(RUNTIME) pull docker.io/mariadb:$(DB_TAG)
	$(info pulling php-apache image)
	$(RUNTIME) pull docker.io/php:$(WEB_TAG)