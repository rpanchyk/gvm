# Go Version Manager

A version manager for Golang SDKs.

## Usage

The following is shown if `gvm` executed:

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
  update      Install the latest Go version
  version     Shows version of gvm

Flags:
  -h, --help   help for gvm
```

## Structure

The `gvm` is installed to `~/.gvm` directory.

### Executables

Directory `~/.gvm/bin` contains executables.

### Configuration

File `~/.gvm/config.toml` is a configuration file.

```toml
# URL where Go versions are looked
all_releases_url = "https://go.dev/dl/"

# Directory where SDKs will be downloaded (absolute or relative path).
download_dir = "./downloads"

# Directory where SDKs will be installed.
# Typically, it is GOROOT env.
# Can be absolute or relative path.
install_dir = "./sdk"

# Directory where local binaries and cache located.
# Typically, it is GOPATH env.
# Can be absolute or relative path.
local_dir = "./local"

# Max versions number to show in list
limit = 10

# Show versions having same OS with current environment
filter_os = true

# Show versions having same architecture with current environment
filter_arch = true
```

## How it works

- Verify specified SDK version is available.
- Download archive from https://go.dev/dl page.
- Verify checksum of fetched archive file (Not implemented).
- Extract archive to `{install_dir}/go{version}` directory.
- Create `{local_dir}/go{version}` directory.
- Set `GOROOT` and `GOPATH` env variables (Implemented for Windows only so far).
- Update `PATH` env variable with `GOROOT/bin` and `GOPATH/bin` (Implemented for Windows only so far).

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
