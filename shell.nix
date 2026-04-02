{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  packages = [
    pkgs.go
    pkgs.gopls
    pkgs.delve
    pkgs.gcc
    pkgs.pkg-config
    pkgs.libX11.dev
    pkgs.libGL
    pkgs.libxcursor
    pkgs.libxi
    pkgs.libxinerama
    pkgs.libxrandr
    pkgs.libxxf86vm
    pkgs.libxkbcommon
    pkgs.wayland
  ];

  hardeningDisable = [ "fortify" ];

  SDL_VIDEODRIVER = "wayland";
}
