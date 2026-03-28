{
  description = "georgslauf";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" ];
    in {

      devShells = nixpkgs.lib.genAttrs systems (system:
        let pkgs = import nixpkgs { inherit system; };
        in {
          default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [
              just
              go_1_25 gotools gopls go-outline gopkgs gocode-gomod godef golint
              air atlas nodejs_24 sqlite vscodium
              tailwindcss_3 sqlc templ esbuild # build tools
            ];
          };
        });

      packages = nixpkgs.lib.genAttrs systems (system:
        let pkgs = import nixpkgs { inherit system; };
        in {

          # Build `nix build .#alertmanager-signal`
          georgslauf = pkgs.buildGoModule rec {
            pname = "georgslauf";
            version = "2.0.11";
            vendorHash = null;

            src = pkgs.lib.cleanSource ./.;

            nativeBuildInputs = [
              pkgs.sqlc
              pkgs.templ
              pkgs.tailwindcss_3
              pkgs.esbuild
            ];

            preBuild = ''
              go generate
            '';

            buildFlags = [
              "-trimpath"
            ];

            ldflags = [
              "-X 'georgslauf/internal/handler.version=v${version}'"
              # TODO "-X 'georgslauf/internal/handler.buildTime=${BUILD_TIMESTAMP}'"
              #"-w" # omit DWARF symbol table
              #"-s" # omit symbol table and debug information
              #"-extldflags=-static"
            ];

            env = {
              CGO_ENABLED = 1;
            };
          };

          # # Container Image `nix build .#image`
          # image = pkgs.dockerTools.buildImage {
          #   name = "docker.io/schlauerlauer/alertmanager-webhook-signal";
          #   tag = "1.1.1";
          #   created = "now";
          #   copyToRoot = pkgs.buildEnv {
          #     name = "image-root";
          #     paths = [
          #       self.packages.${system}.alertmanager-signal
          #       pkgs.cacert
          #     ];
          #   };
          #   config = {
          #     Cmd = [ "/bin/alertmanager-signal" ];
          #     Labels = {
          #       "org.opencontainers.image.title" = "alertmanager-webhook-signal";
          #       "org.opencontainers.image.description" = "A simple program to forward Alertmanager and Grafana Alerts to Signal";
          #       "org.opencontainers.image.source" = "https://github.com/schlauerlauer/alertmanager-webhook-signal";
          #       "org.opencontainers.image.revision" = "1.1.1";
          #     };
          #   };
          # };

        }
      );


    };
}
