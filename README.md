# Emacs Launcher

Use the emacs daemon to control actions using `emacs-client`

## Usage
For this to work correctly, the emacs daemon has to be running (`(server-start)` or `emacs --daemon`)
```bash
elauncher '(<elisp-function>)'
```

## Functionality
1. First search for an existing `emacs-client` instance and launch `elisp-function` in that window
2. If non exist, create a new frame with `elisp-function`

## Installation
1. Using go: Clone the repository and run `go install`
2. Using flake:
   - Include the flake in your configuration
     ```nix
     elauncher.url = "github:chriswifn/elauncher";
     elauncher.inputs.nixpkgs.follows = "nixpkgs";
	 ```
   - Install the package as a regular package
     ```nix
     inputs.elauncher.packages.${pkgs.system}.default
	 ```
