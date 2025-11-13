#!/usr/bin/env bash
# Start script para o backend Veritas
# Uso: ./start.sh [port] [data_file]

PORT=${1:-${HTTP_PORT:-8080}}
DATA_FILE=${2:-${DATA_FILE:-backend/data/tasks.json}}

echo "Iniciando backend na porta ${PORT} com arquivo de dados ${DATA_FILE}"

# Exporta variáveis para que `go run` as veja
export HTTP_PORT=${PORT}
export DATA_FILE=${DATA_FILE}

# Cria diretório do data file se necessário
DIR=$(dirname "${DATA_FILE}")
mkdir -p "${DIR}"

# Executa
go run .
