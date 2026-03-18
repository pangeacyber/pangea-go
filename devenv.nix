{
  pkgs,
  lib,
  config,
  inputs,
  ...
}:

{
  languages.go = {
    enable = true;
    lsp.enable = true;
  };
}
