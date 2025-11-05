{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  packages = [
    pkgs.go
    pkgs.delve
    pkgs.gcc
    pkgs.pkg-config
    pkgs.SDL2
    pkgs.xorg.libX11.dev
  ];

  hardeningDisable = [ "fortify" ];

  SDL_VIDEODRIVER = "wayland";
}
