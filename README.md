# Jailpack

A jail management tool inspired by Kubernetes and Docker, but tailored for FreeBSD and its philosophy: simplicity, stability, performance, security.

## Features

- **Build applications** into portable containers (.cage.tar.gz)
- **Run jails** from ready-made containers
- **Manage** running jails
- **Simple configuration** via jailpack.yaml

## Installation

```bash
# Clone repository
git clone <repository-url>
cd jailpack

# Build
go build -o jailpack

# Install (optional)
sudo cp jailpack /usr/local/bin/
```

## Usage

### Building an Application

```bash
# Create Cage from application directory
jailpack build ./my-app

# Specify output filename
jailpack build ./my-app -o my-app.cage.tar.gz
```

### Running a Cage

```bash
# Run with default parameters
jailpack run my-app.cage.tar.gz

# Run with custom parameters
jailpack run my-app.cage.tar.gz --name my-jail --ip 10.0.0.20
```

### Managing Jails

```bash
# View running jails
jailpack list
```

## Configuration

Create `jailpack.yaml` in your project root to configure the build:

```yaml
# jailpack.yaml — configuration for Cage build
name: my-application
version: 1.0.0
description: "Application description"

# Build settings
build:
  output: my-app.cage.tar.gz
  ignore:
    - .git
    - node_modules
    - *.log

# Run settings
run:
  default_name: my-jail
  default_ip: 10.0.0.10
  ports:
    - 8080:8080
```

## Requirements

- FreeBSD 13.0+
- Go 1.22+
- Administrator privileges for jail creation

## Architecture

Jailpack creates a **Cage** — a self-contained container containing:

- `rootfs/` — minimal filesystem
- `app/` — your application
- `app-start.sh` — startup script
- `config.json` — container metadata

## Development Status

Project is under active development. See [TODO.md](TODO.md) for development plans.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Create a Pull Request


## Links

- [FreeBSD Jail Documentation](https://docs.freebsd.org/en/books/handbook/jails/)
