# Go Version Manager

A version manager for Golang SDKs.

## Structure

The `gvm` is installed to `~/.gvm` directory.

### Executables

Directory `bin` contains executables.

### Configuration

File `config.toml` is a configuration file.

```toml
[main]
# Absolute or relative path where SDKs will be installed in `go{version}` directory.
sdk_dir = "./sdk"
```

## How it works

- Verify specified SDK version is available.
- Download archive from https://go.dev/dl page.
- (Not implemented) Verify checksum of fetched archive file.
- Extract archive to `{sdk_dir}/go{version}` directory.
- (Not implemented) Set `GOROOT` and `GOPATH` env variables.
- (Not implemented) Update `PATH` env variable.

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
