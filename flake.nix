{
  description = "Emacs Launcher for Hyprland";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages.${system}.default = pkgs.buildGoModule {
        pname = "elauncher";
        version = "0.0.1";
        src = ./.;
        vendorHash = "sha256-7K17JaXFsjf163g5PXCb5ng2gYdotnZ2IDKk8KFjNj0=";
      };
      devShells.${system}.default = pkgs.mkShell { buildInputs = [ pkgs.go ]; };
    };
}
