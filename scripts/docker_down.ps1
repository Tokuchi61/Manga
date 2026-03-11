$ErrorActionPreference = 'Stop'
docker compose -f deploy/docker-compose.yml down --remove-orphans
