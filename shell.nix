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
