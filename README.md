# chuck (Container Image Update Checker)

A command-line tool that helps you stay up-to-date with your container images by checking for newer versions.

## Overview

`chuck` simplifies container image maintenance by checking if newer versions of your container images are available.

## Features

- Check for newer versions of container images
- Support for multiple container registries
- Simple and intuitive command-line interface

## Installation

```bash
# Using go install
go install github.com/FedericoAntoniazzi/chuck

# From source
git clone https://github.com/FedericoAntoniazzi/chuck
cd chuck
go build .
```

## Usage

```bash
$ chuck
chuck fetches the images from running containers and shows eventual updates

Usage:
  chuck [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Run Chuck

Flags:
  -h, --help   help for chuck

Use "chuck [command] --help" for more information about a command.
```

## Examples

Check if a newer version exists:
```bash
$ chuck check
CONTAINER	IMAGE			VERSION UPDATE
my-web-server	docker.io/nginx:1.23	1.27.4
```

Change output format
```bash
$ chuck check --output text
Container my-web-server (docker.io/nginx:1.23) can be upgraded to 1.27.4
```

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support
If you encounter any issues or have questions, please file an issue on the GitHub repository.
