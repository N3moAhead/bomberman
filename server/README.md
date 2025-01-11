# Bomberman Game Server

This repository currently contains only the basic setup. In the future, this server will be able to handle multiple clients simultaneously and enable players to enjoy Bomberman together from different clients.

## Setup

All development is carried out within Docker containers. This approach ensures that all required packages are installed only within the Docker container, regardless of the host system, thus avoiding versioning issues.

### Prerequisites

The only dependencies at the moment are:  
- `make`  
- `podman`  
- `podman-compose`  

### Getting Started

To start the project, simply run the following command in your terminal:  

```bash
make dev
```
