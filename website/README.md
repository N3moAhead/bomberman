# Website

This directory contains the web frontend and matchmaking service for the Bomberman bot competition platform.

## Components

The project consists of two main Go applications:

1.  **`cmd/website`**: A web server that provides the user interface. It handles user authentication via GitHub, bot registration (Docker images), and displays leaderboards and user-specific dashboards.
2.  **`cmd/matchmaker`**: A background service that orchestrates matches. It periodically queries the database for unmatched bot pairs, publishes match jobs to a RabbitMQ queue, and processes the results.

## Technology Stack

*   **Backend**: Go
*   **Web Framework**: Chi
*   **ORM**: GORM (PostgreSQL)
*   **Frontend**: Templ (Go templating), htmx
*   **Styling**: Tailwind CSS, daisyUI
*   **Authentication**: Goth (GitHub OAuth2)
*   **Message Broker**: RabbitMQ
*   **Development**: `air` for live-reloading

## Local Development

### Prerequisites

*   Go
*   Node.js / npm
*   A running PostgreSQL and RabbitMQ instance.

### Setup

1.  **Configuration**:
    Create a `.env` file from the example and populate it with your configuration, including database credentials and GitHub OAuth application details.
    ```bash
    cp .env_example .env
    ```

2.  **Dependencies**:
    Install the required Node.js packages.
    ```bash
    npm install
    ```

### Running the Services

1.  **Web Server**:
    To run the web server with live-reloading enabled, use the `dev` Makefile target. This also installs the `air` dependency if not present.
    ```bash
    make dev
    ```
    The server will be available at the `PORT` specified in your `.env` file (defaults to `:3000`).

2.  **Matchmaker**:
    Run the matchmaking service in a separate terminal:
    ```bash
    go run ./cmd/matchmaker/main.go
    ```
