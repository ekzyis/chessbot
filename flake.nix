{
  description = "chessbot - play chess on SN";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

  outputs = { self, nixpkgs }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in
    {
      packages.x86_64-linux.default = pkgs.buildGoModule {
        name = "chessbot";
        src = ./.;

        vendorHash = "sha256-xEj9D18AGDQXSjW/DhViBr1vcAcIVfZXL9jOw4ez18U=";

        buildInputs = [ pkgs.sqlite ];

        preBuild = ''
          export CGO_ENABLED=1
          export CGO_CFLAGS="-I${pkgs.sqlite.dev}/include"
          export CGO_LDFLAGS="-L${pkgs.sqlite}/lib"
        '';
      };

      apps.x86_64-linux.default = {
        type = "app";
        program = "${self.packages.x86_64-linux.default}/bin/chessbot";
      };

      devShells.x86_64-linux.default = pkgs.mkShell {
        buildInputs = with pkgs; [ go sqlite gcc ];
        CGO_ENABLED = "1";
        shellHook = ''
          export CGO_CFLAGS="-I${pkgs.sqlite.dev}/include"
          export CGO_LDFLAGS="-L${pkgs.sqlite}/lib"
        '';
      };
    };
}
