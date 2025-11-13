# Desafio Fullstack — Veritas

Este repositório contém o desafio técnico (backend em Go e frontend em React) usado no processo seletivo.

## Requisitos

- Go (1.20+ recomendado, o módulo usa go 1.22)
- Node.js (v16+ recomendado) e npm/yarn
- Git

## Clonar o repositório

Se você estiver no GitHub, clone com:

```bash
git clone https://github.com/IanSantosRodrigues/desafio-fullstack-veritas.git
cd desafio-fullstack-veritas
```

Ou, se já tiver o código localmente, entre na pasta do projeto:

```bash
cd /caminho/para/desafio-fullstack-veritas
```

## Rodando frontend + backend juntos

1. Comece o backend (ex: porta 8080):

```bash
cd backend
go run .
```

2. Em outro terminal, rode o frontend:

```bash
cd frontend
npm install
npm run dev
```

Abra o navegador em `http://localhost:3000` (ou porta exibida pelo dev server do frontend).