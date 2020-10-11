# Polyglot ![CI](https://github.com/filhodanuvem/polyglot/workflows/CI/badge.svg)

What are your favorite programming languages? 👅
Polyglot tells you that based in a github username.

# How to use

```bash
polyglot --username=filhodanuvem
```

```bash
polyglot -s #To run on server mode
```

Usage:

```bash
Usage:
  Polyglot [flags]

Flags:
  -h, --help              help for Polyglot
  -l, --log string        Log verbosity, options [debug, info, warning, error, fatal] (default "fatal")
  -o, --output string     Path to log in a file
  -p, --path string       Path where to download the repositories (default "/tmp/polyglot")
  -u, --username string   Username
  -s, --server            Run polyglot API Server
  --port string           IP address for the server (default "127.0.0.1")
  --host string           Port for the server (default "8080")
```

Server Usage:

```bash
curl http://127.0.0.1:8080/\?user\=santana125\&limit\=10
```
