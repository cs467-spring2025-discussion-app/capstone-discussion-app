# OpenForum

A project submitted to fulfill the Oregon State University CS467 Capstone Project course requirements.

## Team members:

- Alessandro Ceserani
- Samuel Hoover
- Joseph Keehnast
- Adam Leatherman

# Quickstart

Our application was developed locally using interacting Docker services. We used a `docker-compose` configuration to set up the services the services to communicate over an app network.

The latest build of our app uses the following minimal Docker versions:

```
Docker version 28.1.1
Docker Compose version 2.36.1
```

Check your docker versions with:

```bash
docker --version
docker compose version
```

If you have Docker > `28.1.1` installed on your system, you can run:

```bash
docker compose up --build
```

Then go to

```
http://localhost:4173
```

to see a local build for your app.

# Installation

## Set up Docker

### Install the latest version of docker

General installation instructions can be found at [https://docs.docker.com/get-started/get-docker/](https://docs.docker.com/get-started/get-docker/)

#### Linux

We generally recommend the convenience script provided by Docker to install the latest version:

```bash
# Install curl if you don't already have it
sudo apt install curl
# Download the docker installation script
curl -fsSL https://get.docker.com -o get-docker.sh
# Run the script
sudo sh test-docker.sh
```

Or use your distribution's package manager.

```bash
# Arch Linux
pacman -Syu
pacman -S docker
```

For distributions using `apt`, we recommend following the official installation instructions to guarantee the latest stable build:

[https://docs.docker.com/engine/install/debian/#install-using-the-repository](https://docs.docker.com/engine/install/debian/#install-using-the-repository)

#### macOS:

On macOS, we recommend installing the native app to ensure the Docker services will start without issues:

Follow the download links here:

[docs.docker.com/desktop/setup/install/mac-install](docs.docker.com/desktop/setup/install/mac-install)

or use the Homebrew package manager: [https://brew.sh/](https://brew.sh/)

#### Windows

On Windows, we recommend installing Docker Desktop:

[https://docs.docker.com/desktop/setup/install/windows-install/](https://docs.docker.com/desktop/setup/install/windows-install/)

This requires Windows Subsystem For Linux (WSL)

[https://learn.microsoft.com/en-us/windows/wsl/install](https://learn.microsoft.com/en-us/windows/wsl/install)

```bash
brew install --cask docker
```

### Start Docker

#### Linux

Start the daemon (assumes your Linux distro is using systemd)

```bash
sudo systemctl start docker
```

#### macOS

Open the Docker app to ensure the services are running in the background.

## Build the app!

```bash
 # Build and start the services
docker compose up -d --build
# (optional) Follow the output of the logs
docker logs -f

# Stop the services
docker compose down
```

Depending on your user profile, you **may** need to run the docker commands with elevated privileges:

```bash

sudo docker compose up -d --build
```

### justfile

We used a `justfile` for development convenience. Installation instructions for just can be found at

[https://just.systems/man/en/](https://just.systems/man/en/)

Here are some useful commands:

```bash
just start # build the services with any new changes and start the services
just stop  # stop all active services
just flogs # follow the logs of active services

just empty {{db}} # pass `db` or `test_db` to empty the database volumes
just fresh        # build the app completely from scratch
```

## See the app!

### Production build

The Docker service runs a production build of the Vite-React app at

```
http://localhost:4173
```

### Dev build

Development was done on a Vite development build. You can see this by going to the `frontend/` directory of the project and running the appropriate dev command:

```bash
npm run dev
```

This requires `Node.js`. System specific instructions to install Node can be found at

[https://nodejs.org/en/download](https://nodejs.org/en/download)
