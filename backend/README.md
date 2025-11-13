# Backend — Veritas (Instruções de execução)

Este diretório contém o backend em Go. A implementação atual suporta persistência simples em JSON (arquivo) e pode ser configurada via variáveis de ambiente.

Principais variáveis de ambiente

- `HTTP_PORT` — porta HTTP onde o servidor vai escutar. Padrão: `8080`.
- `DATA_FILE` — caminho para o arquivo de persistência (JSON). Padrão: `backend/data/tasks.json`.
- Futuro: `STORAGE_TYPE` (json|sqlite|postgres) — seleciona o tipo de armazenamento.

SQLite (opcional, build-tag)

É possível usar SQLite como backend. Para não forçar dependências em ambientes offline, o suporte a SQLite está separado por build-tags.

- Para rodar com SQLite, instale o driver e compile/execute com a tag `sqlite`:

```bash
cd backend
go get github.com/mattn/go-sqlite3
DATABASE_URL=backend/data/tasks.db STORAGE_TYPE=sqlite go run -tags sqlite .
```

Ou usando `go build`:

```bash
cd backend
go get github.com/mattn/go-sqlite3
go build -tags sqlite -o veritas-backend .
DATABASE_URL=backend/data/tasks.db ./veritas-backend
```

Observações:
- O arquivo `storage/sqlite_enabled.go` contém a implementação que usa o driver `github.com/mattn/go-sqlite3` e só é compilada quando a tag `sqlite` é fornecida.
- Se a tag não for passada (padrão), o projeto usa o armazenamento JSON (DATA_FILE) e os testes continuam verdes sem precisar do driver.

Rodando localmente (desenvolvimento)

1. Instale dependências (uma vez):

```bash
cd backend
go mod tidy
```

2. Rodar em modo desenvolvimento (usando arquivo JSON padrão):

```bash
cd backend
go run .
```

3. Rodar com variáveis customizadas:

```bash
HTTP_PORT=3000 DATA_FILE=/tmp/veritas-tasks.json go run .
```

Script útil

- `start.sh` — script simples para iniciar o backend definindo valores padrão (não é obrigatório, é um atalho). Torne-o executável:

```bash
chmod +x backend/start.sh
./backend/start.sh
```

Persistência de dados

- Por padrão os dados são salvos em `backend/data/tasks.json`. Esse caminho pode ser mudado com `DATA_FILE`.
- Em produção, recomenda-se usar um banco de dados (Postgres) ou SQLite com um volume persistente.

Boas práticas

- Não comite o arquivo de dados. O repositório já contém regras em `.gitignore` para `backend/data/`.
- Proteja credenciais de banco em variáveis de ambiente ou secret managers.
- Faça backup periódico dos dados.

Próximos passos (opcional)

- Implementar `STORAGE_TYPE` e adapters para SQLite/Postgres.
- Adicionar um Dockerfile e `docker-compose` para facilitar deploy.

Se quiser, eu implemento o adapter SQLite agora (troca pequena e útil) — me diga se prefere isso.