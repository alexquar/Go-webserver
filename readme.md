# U-Watchlist

U-Watchlist is a custom web server project built with Go. It provides a platform to manage and track watchlists, such as movies or TV shows, through a web interface.

## Features

- Custom HTTP server implemented in Go
- RESTful API for managing watchlist items
- Simple web UI for adding, viewing, and removing items
- Persistent storage using local files or a database

## Tech Stack

- Go
- HTMX
- SQLServer

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.18 or newer

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/U-Watchlist.git
    cd U-Watchlist
    ```
2. Build the server:
    ```sh
    go build -o uwatchlist
    ```
3. Run the server:
    ```sh
    ./uwatchlist
    ```

4. Open your browser and navigate to `http://localhost:8080`
