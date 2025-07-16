# Chuck

![Chuck Logo](https://img.shields.io/badge/Chuck-Container%20Update%20Checker-blue?style=for-the-badge&logo=docker)
![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

Chuck is a lightweight, self-contained utility written in Go that helps you keep your Docker container images up-to-date. It scans running Docker containers, parses their image names, queries their respective registries for available tags, and identifies if newer semantic versions are available.

## Getting Started

### Prerequisites

* Go 1.24 or higher installed.
* Docker Desktop or Docker Engine running on your system.
* Chuck requires access to the Docker daemon (usually via `/var/run/docker.sock` or `DOCKER_HOST` environment variable).

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/FedericoAntoniazzi/chuck.git
    cd chuck
    ```

2.  **Build the binary:**
    ```bash
    go build -o chuck
    ```
    This will create an executable named `chuck` in your current directory.

### Usage

Run Chuck from your terminal:

```bash
./chuck [flags]
```

Example
```shell
❯ chuck -output tab
2025-07-17T00:22:04.663+0200	warn	chuck/main.go:152	skipping unsupported registry (quay.io) for image quay.io/nginx/nginx-unprivileged:1.25
CONTAINER_NAME	IMAGE	CURRENT TAG	LATEST TAG
mywebserver	nginx	1.21		1.29.0
```

## Roadmap

### Phase 1: Basic Update Detection (Current / Completed)
- [x] Scan running Docker containers.
- [x] Parse image names into registry, namespace, name, and tag.
- [x] Query Docker Hub for available image tags.
- [x] Perform semantic version comparison to detect updates.
- [x] Report updates to standard output/log file.
- [x] Gracefully skip non-semantic version tags.
- [x] Show results in a tabbed format (such as `kubectl`)

### Phase 2: Configuration & Status Reporting
- [ ] Develop configuration management using a `chuck.yaml` file, supporting XDG Base Directory Specification for config location.
- [ ] Implement token/credential management within `chuck.yaml` for future authentication needs.
- [ ] Develop status reporting to a text file (YAML, JSON, or CSV format, user-selectable). This will be the base for notifications.

### Phase 3: Notifications
- [ ] Develop notification integration for Telegram, utilizing the configuration from Phase 2 for tokens.
- [ ] Implement a generic notification interface to allow for easy expansion to other platforms (e.g., Slack, Email).

### Phase 4: Daemon Mode
- [ ] Develop a mode for running Chuck as a daemon.
- [ ] Implement periodic execution at a configurable interval.
- [ ] Integrate with the notification system to alert users about updates while running as a daemon.

### Phase 5: Persistence
- [ ] Implement SQLite database integration to:
    - [ ] Store discovered updates (current version, latest available, last checked time).
    - [ ] Track acknowledged updates or ignored images.
    - [ ] Prevent repetitive notifications for already known updates.

### Phase 6: Custom Registries & Authentication
- [ ] Implement specific clients for additional custom/self-hosted registries (e.g., Nexus, Artifactory, Harbor).
- [ ] Enhance registry clients with robust authentication mechanisms for private repositories and to overcome public registry rate limits (e.g., Docker Hub authentication flow, basic auth, token support) using the configuration from Phase 2.

### Phase 7: Advanced Features & Usability
- [ ] Add filtering capabilities (e.g., exclude certain images/containers, include only specific registries).
- [ ] Output customization beyond file formats (e.g., custom templates).
- [ ] Support for other container runtimes (e.g., Containerd, Podman).
- [ ] Potentially implement options for scheduling and more continuous monitoring (though daemon mode covers much of this).

