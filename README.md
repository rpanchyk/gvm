# Go Version Manager

A version manager for Golang SDKs.

## Installation

1. Unpack installation archive to `$HOME` directory.
1. Add `$HOME/.gvm` to `$PATH` environment variable.

## Configuration

File `$HOME/.gvm/config.toml` is single configuration file.

```toml
# URL where Go versions are looked.
release_url = "https://go.dev/dl/"

# Directory where SDK archives are downloaded.
# Can be absolute or relative path.
download_dir = "./downloads"

# Directory where SDKs are installed.
# Typically, it is GOROOT env.
# Can be absolute or relative path.
install_dir = "./sdk"

# Directory where local binaries and cache located.
# Typically, it is GOPATH env.
# Can be absolute or relative path.
local_dir = "./local"

# Max versions number to show in list.
list_limit = 10

# Show versions having same OS with current environment.
list_filter_os = true

# Show versions having same architecture with current environment.
list_filter_arch = true

# File with list of cached SDKs in JSON format.
list_cache_file = "./cache.json"

# Cache expiration time in minutes.
list_cache_ttl = 1440

# Defines in which file is added loading config string.
# Non-windows property.
# Examples: ~/.profile, ~/.bashrc, ~/.zshrc, etc.
unix_shell_config = "~/.profile"
```

## Usage

The following is shown if `gvm --help` executed:

```
Usage:
  gvm [flags]
  gvm [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  default     Set specified Go version as default
  download    Download specified Go version
  help        Help about any command
  install     Install specified Go version
  list        Shows list of available Go versions
  remove      Remove specified Go version
  update      Update Go to the latest version and set it as default
  version     Shows version of gvm

Flags:
  -h, --help   help for gvm
```

For example, let's add Go new version.

1. Show available Go versions list: `gvm list`

```
  1.22.2 linux amd64
  1.22.1 linux amd64
  1.22.0 linux amd64
  1.21.9 linux amd64
  1.21.8 linux amd64
  1.21.7 linux amd64
  ...
```

2. Install Go `1.22.0` version: `gvm install 1.22.0`

```
SDK has been installed: /home/user/.gvm/sdk/go1.22.0
Local directory created: /home/user/.gvm/local/go1.22.0
```

3. Set Go `1.22.0` version as default: `gvm default 1.22.0`

```
User environment is set to go1.22.0 version as default
```

4. (_Important!_) Restart terminal to apply changes.

5. Check Go versions list again to see the asterisk near default version: `gvm list`

```
  1.22.2 linux amd64
  1.22.1 linux amd64
* 1.22.0 linux amd64 [downloaded] [installed]
  1.21.9 linux amd64
  1.21.8 linux amd64
  1.21.7 linux amd64
  ...
```

6. Verify everything is expected: `go version && go env GOROOT GOPATH && echo $GOROOT && echo $GOPATH`

```
go version go1.22.0 linux/amd64

/home/user/.gvm/sdk/go1.22.0
/home/user/.gvm/local/go1.22.0

/home/user/.gvm/sdk/go1.22.0
/home/user/.gvm/local/go1.22.0
```

## How it works

- Verify specified SDK version is available.
- Download SDK archive from https://go.dev/dl page.
- Verify checksum of fetched file (Not implemented).
- Extract archive to `{install_dir}/go{version}` directory.
- Create `{local_dir}/go{version}` directory.
- Set `GOROOT` and `GOPATH` environment variables.
- Add `GOROOT/bin` and `GOPATH/bin` to `PATH` environment variable.

## Alternatives

This project was started for learning of Go Programming Language purposes. Since author comes from Java world,
initially it was inspired by [SDKMAN!](https://sdkman.io) but later found similar projects, like:

- `moovweb/gvm` - https://github.com/moovweb/gvm
- `GoTV` - https://go101.org/apps-and-libs/gotv.html

Starting in Go `1.21` introduced [Go toolchain](https://go.dev/doc/toolchain), which is the standard library
as well as the compiler, assembler, and other tools.

## Disclaimer

The software is provided "as is", without warranty of any kind, express or
implied, including but not limited to the warranties of merchantability,
fitness for a particular purpose and noninfringement. in no event shall the
authors or copyright holders be liable for any claim, damages or other
liability, whether in an action of contract, tort or otherwise, arising from,
out of or in connection with the software or the use or other dealings in the
software.

## Contribution

If you have any ideas or inspiration for contributing the project,
please create an [issue](https://github.com/rpanchyk/gvm/issues/new)
or a [pull request](https://github.com/rpanchyk/gvm/pulls).
