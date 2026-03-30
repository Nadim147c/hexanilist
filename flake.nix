{
  inputs.nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";

  outputs =
    { nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      perSystem = f: nixpkgs.lib.genAttrs systems (system: f (import nixpkgs { inherit system; }));
    in
    {
      packages = perSystem (pkgs: rec {
        api = pkgs.callPackage ./api/default.nix { };
        api-docker = pkgs.dockerTools.buildLayeredImage {
          name = "api";
          tag = "latest";
          config.Cmd = [ "${api}/bin/api" ];
        };
      });

      devShells = perSystem (pkgs: {
        default = pkgs.mkShell {
          name = "media";
          buildInputs = with pkgs; [
            typescript-go
            prettier
            just-lsp
            just
          ];
        };
      });
    };
}
