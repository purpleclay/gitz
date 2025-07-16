{
  description = "Write fluent interactions to Git. Programmatically crafting git commands becomes a breeze!";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
        pname = "gitz";
      in
        with pkgs; {
          devShells.default = mkShell {
            buildInputs = [
              alejandra
              git
              gofumpt
              golangci-lint
              go_1_23
              nil
              just
              tparse
              typos
            ];
          };

          # TODO: ensure this formatting command works against the code base
          # TODO: do we switch over to nativeBuildInputs reach nix docs
          # TODO: before we check in any code, should we run githooks?

          packages.format = pkgs.writeShellScriptBin "${pname}-format" ''
            echo "Running gofumpt..."
            ${pkgs.gofumpt}/bin/gofumpt -w .
            exit_code=$?
            if [ $exit_code -eq 0 ]; then
              echo "✓ Formatting completed!"
            else
              echo "✗ Formatting failed (exit code: $exit_code)"
            fi
            exit $exit_code
          '';

          packages.lint = pkgs.writeShellScriptBin "${pname}-lint" ''
            echo "Running golangci-lint..."
            ${pkgs.golangci-lint}/bin/golangci-lint run ./...
            exit_code=$?
            if [ $exit_code -eq 0 ]; then
              echo "✓ No linting issues found!"
            else
              echo "✗ Linting issues detected (exit code: $exit_code)"
            fi
            exit $exit_code
          '';

          packages.test = pkgs.writeShellScriptBin "${pname}-test" ''
            export GOMODCACHE="''${GOMODCACHE:-$HOME/go/pkg/mod}"
            export GOCACHE="''${GOCACHE:-$HOME/.cache/go-build}"
            mkdir -p reports

            echo "Running tests..."
            ${pkgs.go_1_23}/bin/go test \
              -C ${./.} \
              -short \
              -race \
              -vet=off \
              -shuffle=on \
              -p 1 \
              -covermode=atomic \
              -json ./... | tee reports/unittest.json | ${pkgs.tparse}/bin/tparse -follow
            exit_code=$?

            if [ $exit_code -eq 0 ]; then
              echo "✓ All tests passed!"
            else
              echo "✗ Some tests failed (exit code: $exit_code)"
            fi
            exit $exit_code
          '';
        }
    );
}
