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
    pkgs.libGL
    pkgs.xorg.libXcursor
    pkgs.xorg.libXi
    pkgs.xorg.libXinerama
    pkgs.xorg.libXrandr
    pkgs.xorg.libXxf86vm
    pkgs.libxkbcommon
    pkgs.wayland
  ];

  hardeningDisable = [ "fortify" ];

  SDL_VIDEODRIVER = "wayland";
}
