<h1 align="center">
    Presto
</h1>

<div align="center">
    <img src="https://wakatime.com/badge/github/aktshually/presto.svg" alt="wakatime" style="max-width:100%">
</div>

Presto is a wrapperless Discord bot made in Go that aims to be efficient yet multipurpose.

## Summary

- [📦 Installation](#-installation)
    - [🔧 Prerequisites](#-prerequisites)
    - [🚀 Getting Started](#-getting-started)

## 📦 Installation

### 🔧 Prerequisites
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### 🚀 Getting Started

1. Clone the repository and navigate to the project directory:
```bash
git clone https://github.com/aktshually/presto.git && cd presto
```
Or with SSH:
```bash
git clone git@github.com:aktshually/presto.git && cd presto
```

2. Rename the `.env.example` file to `.env.docker`:
```bash
mv ./.env.example ./.env.docker
```

3. Fill the empty data

4. Start the development environment:
```bash
docker compose up
```

5. After following the previous steps correctly, your bot should be online and should have registered the commands. Notice that if `PRESTO_ENVIRONMENT` is `production`, it might take some time to register the commands.
