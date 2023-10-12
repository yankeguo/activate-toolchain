# activate-toolchain

A toolchain download and activation tool with various mirror sites embedded

## Installation

You can either build binary from source, or just download pre-built binary.

* Build from source

    ```shell
   git clone https://github.com/guoyk93/activate-toolchain.git
   cd activate-toolchain
   go build -o activate-toolchain ./cmd/activate-toolchain
    ```

* Download pre-built binaries

  View <https://github.com/guoyk93/activate-toolchain/releases>

## Usage

Only support shell in `POSIX` environment.

```shell
eval "$(activate-toolchain node@16.2 jdk@17 maven@3.8)"
```

## Supported Toolchains and Version Formats

* `node`
    * `node@18`
    * `node@18.18`
    * `node@18.18.1`
* `jdk`
    * `jdk@8`
    * `jdk@8.0.372`
    * `jdk@8.0.372+7`

## Credits

GUO YANKE, MIT License
