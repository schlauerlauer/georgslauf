{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
	buildInputs = [
		go
		gotools
		gopls
		go-outline
		gopkgs
		gocode-gomod
		godef
		golint

		air
		sqlc
		templ
		atlas

		nodejs_23

		sqlite
	];
}

