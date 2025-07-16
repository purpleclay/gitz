alias l := lint
alias t := test
alias fmt := format

_default:
    @just --list

lint:
    nix run .#lint

test:
    nix run .#test

format:
    nix run .#format
