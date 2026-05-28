{
  description = "Trabalho 2 de Engenharia de Software";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
        go-migrate-postgres = pkgs.buildGoModule rec {
          pname = "golang-migrate";
          version = "4.19.1";
          src = pkgs.fetchFromGitHub {
            owner = "golang-migrate";
            repo = "migrate";
            rev = "v${version}";
            sha256 = "sha256-Z8ufA2z5XeJ80Jfd6NSls/SurR8rMTO4zq88fQYGGpA=";
          };
          subPackages = [ "cmd/migrate" ];
          tags = [ "postgres" ];
          vendorHash = "sha256-gGwdRyq8uzDwuq6JyxhEp/7M68GN4HG/vQ+ynhxbU1w=";
        };
      in
      {
        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              gotools
              go-tools
              go-migrate-postgres
            ];
          };
        };
      }
    );
}
