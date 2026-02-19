{
	inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

	outputs = { self, nixpkgs }:
		let
			pkgs = nixpkgs.legacyPackages.x86_64-linux.pkgs;
		in {
			# Development environment
			devShells.x86_64-linux.default = pkgs.mkShell {
				nativeBuildInputs = with pkgs; [
					just
					go_1_25 gotools gopls go-outline gopkgs gocode-gomod godef golint
					air sqlc templ atlas nodejs_24 sqlite vscodium
				];
			};

			defaultPackage.x86_64-linux = self.packages.x86_64-linux.my-script;
		};
}
