# TODO: C-Bindings for MacOS Security Framework

Because kind must be installed as a package on macos in order to bundle it, that portion of the process must be run with admin privileges. This is not a requirement for other nix operating system, and potentially not for windows either, though additional research is necessary to confirm. In any case, ideally we would request elevated permissions only when necessary. In most nix operating systems this is as simple as a `polkit` query. MacOS however, handles security via their built-in security libraries. The C code and yaml here represent most of what is required to create appropriate bindings for this process in go, however additional research is required for how to use the `c-for-go` binary to generate these bindings correctly. For now it is categorized as a non-essential item.