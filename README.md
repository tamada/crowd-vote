# Crowd Vote API 📊

Crowd Vote is a lightweight, containerized REST API service designed for real-time reporting and monitoring of location-based crowdedness. 🌐

## Features ✨

- **Group-based Hierarchy**: Manage locations under specific groups. 🏢
- **Real-time Reporting**: Users can submit crowdedness levels (0-3). 🗳️
- **Dynamic Aggregation**: Automatically calculates crowdedness rates (average levels and total votes) based on user-defined time windows. ⏱️
- **Pluggable Storage**: Supports SQLite by default, with an abstraction layer ready for NoSQL (MongoDB). 💾
- **Container-Ready**: Fully dockerized with `docker-compose` support for quick deployment. 🐳

## Getting Started 🚀

### Prerequisites

- Docker & Docker Compose

### Running the API

1. Clone the repository.
2. Build and run the services:
   ```bash
   docker compose up --build
   ```
3. The API will be available at `http://localhost:8080`.

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
