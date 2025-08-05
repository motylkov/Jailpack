# Jailpack

A tool for jail management, inspired by Kubernetes and Docker, but tailored for FreeBSD and its philosophy: simplicity, stability, performance, security.

### Development Goals

- Create CLI with Cobra: jailpack build, run, deploy, list, logs, shell, push, pull, destroy
- Implement jailpack build — packaging application into .jail.tar.gz cage
- Implement jailpack run — running cage as jail
- Implement jailpack list — showing running jails (via jls)
- Implement jailpack logs and shell — access to logs and shell
- Add support for flags: --name, --ip, --port, --env, --output

### Build and Cage
- Support for `jailpack.yaml` for declarative configuration
- Create **Cage** `.jail.tar.gz` or `.cage.tar.gz` with `rootfs/`, `config.json`, `startup.sh` (or other variants)
- Dependencies
- Support for `.jailignore` (like `.dockerignore`) if needed

### Deployment and Orchestration
- `jailpack deploy -f deployment.yaml` — running Cage
- Network support
- `jailpack compose` — orchestration of multiple Cages
- Healthcheck

### Storage and Security
- `jailpack push` / `pull` — sending and receiving **Cage**
- Cage signing (GPG, sha256, ...)
- Using ZFS for storing and cloning Cages

### Integration
- Integration with `FreeBSD-Command-Manager`: calling `jailpack` as backend

### Documentation and Examples
- Migration guide from Docker
- Best practices for jail containerization
- Examples: Go, Python, Node.js, ... 