# Crowd Vote API & Sample Web App 📊

Crowd Vote is a lightweight, containerized REST API service designed for real-time reporting and monitoring of location-based crowdedness. It now also ships with a simple sample web app for local testing. 🌐

## Features ✨

- **Group-based Hierarchy**: Manage locations under specific groups. 🏢
- **Real-time Reporting**: Users can submit crowdedness levels (0-3). 🗳️
- **Dynamic Aggregation**: Automatically calculates crowdedness rates (average levels and total votes) based on user-defined time windows. ⏱️
- **Pluggable Storage**: Supports SQLite by default, with an abstraction layer ready for NoSQL (MongoDB). 💾
- **Container-Ready**: Fully dockerized with `docker-compose` support for quick deployment. 🐳
- **Sample Web App**: Provides a browser-based UI for creating demo locations and submitting votes against the local API. 🧪

## Getting Started 🚀

### Prerequisites

- Docker & Docker Compose

### Running the API and Sample Web App

1. Clone the repository.
2. Build and run the services:
   ```bash
   ./launch-sample-webapp.sh
   ```
   You can also run `docker compose up --build` directly.
3. Open `http://localhost:8080` for the sample web app.
4. The REST API is also available at `http://localhost:8080`.
5. Stop the local stack with `Ctrl+C`, then run `docker compose down` if you want to remove the containers.

## API Usage 📡

### Locations

- `GET /locations/{group}`: List all locations in a group.
- `POST /locations/{group}`: Create a new location.
- `GET /locations/{group}/{id}`: Get location details.

### Votes

- `POST /votes/{group}/{id}`: Submit a vote for a location.
- `GET /votes/{group}/{id}?before=10&unit=minute`: Get aggregated crowdedness rate.
- `GET /votes/{group}`: List crowdedness rates for all locations in a group.

### Histories

- `GET /histories/{group}/{id}?before=1&unit=hour`: Get vote history aggregated.

## Development 🛠️

- **Language**: Go 1.25+
- **Database**: SQLite (default), extensible to MongoDB.
- **Testing**: Run tests with `go test -cover ./...`.

## License 📜

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
