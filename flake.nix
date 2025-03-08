{
	description = "georgslauf";

	#inputs.nixpkgs.url = "github:nixos/nixpkgs/release-24.11";
	inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

	outputs = { self, nixpkgs }:
		let
			pkgs = nixpkgs.legacyPackages.x86_64-linux.pkgs;
		in {
			# Development environment
			devShells.x86_64-linux.default = pkgs.mkShell {
				name = "sendit app development";
				nativeBuildInputs = with pkgs; [
					just
					go_1_24 gotools gopls go-outline gopkgs gocode-gomod godef golint
					air sqlc templ atlas nodejs_23 sqlite
				];
				shellHook = ''
					echo "Welcome in $name"
				'';
			};

			defaultPackage.x86_64-linux = self.packages.x86_64-linux.my-script;
		};
}
