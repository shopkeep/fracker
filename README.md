# fracker

`fracker` is a small utility for converting directory hierarchies in etcd to shell environment variables.

# Usage

```
fracker <key> [<key> <key> ...]
```

`fracker` reads the directory trees at the given keys, collects all the leaf values and writes environment 
variable declarations to standard out.

Assume you have the following directory structure in etcd:

```
/foo
  |__ /bar
  |     |__ /val = "1234"
  |
  |__ /baz
        |__ /woo = "abcd"
        |__ /hoo = "efgh"
```

`fracker foo` will output the following to standard out:

```
FOO_BAR_VAL=1234
FOO_BAZ_WOO=abcd
FOO_BAZ_HOO=efgh
```

Whereas `fracker foo/bar` will only output `FOO_BAR_VAL=1234`.

The output is intended to be redirected into a file such as `/etc/environment` for service configuration.
