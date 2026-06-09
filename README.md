# API de Gestão de Tarefas Colaborativas | Alunos: Eduardo Borsoi e Lucas Renato **
API RESTful para gerenciamento de tarefas em equipe, desenvolvida para a disciplina de **Engenharia de Software – Arquitetura e Padrões** (UNISINOS).
**Alunos: Eduardo Borsoi e Lucas Renato**

---

## Visão Geral

O sistema permite que usuários criem, editem, atribuam e concluam tarefas de forma colaborativa. A API é protegida por autenticação JWT e dispara notificações via Webhooks sempre que uma tarefa é criada, atribuída ou concluída.

**Requisitos Complementares implementados:**
- **Filtro Avançado de Tarefas** — filtragem por status, prioridade e prazo
- **Webhooks** — notificações automáticas integráveis com n8n, Slack, Discord ou qualquer sistema via HTTP

---

## Decisões Arquiteturais

### Arquitetura: Clean Architecture (Hexagonal)

O projeto segue **Clean Architecture**, separando o código em três camadas principais:

```
internal/
├── domain/          ← Regras de negócio puras (entidades, interfaces)
│   ├── user/        ← Entidade User, interface Repository
│   └── task/        ← Entidade Task, interface Repository, status, priority
│
├── application/     ← Casos de uso (serviços)
│   ├── auth/        ← Login, geração e validação de JWT
│   ├── user/        ← CRUD de usuários
│   ├── task/        ← CRUD de tarefas + filtros
│   └── webhook/     ← Disparo de notificações
│
└── infrastructure/  ← Implementações concretas (banco de dados)
    └── db/          ← Repositórios PostgreSQL
```

**Por que Clean Architecture?**
- O domínio não depende do banco de dados — se mudar de PostgreSQL para MongoDB, só muda a camada `infrastructure`
- Facilita testes: usamos mocks no lugar dos repositórios reais
- Cada camada tem responsabilidade única

### Linguagem e Framework

- **Go** — performance e tipagem estática
- **Fuego** — gera documentação Swagger/OpenAPI automaticamente a partir do código
- **PostgreSQL** — banco relacional robusto com suporte a UUID e timestamps

### Autenticação

JWT (JSON Web Token) — o servidor assina um token que o cliente envia em cada requisição. Não há sessão no servidor, o que facilita escalabilidade horizontal.

---

## Modelagem de Dados

### Tabela `users`

| Coluna          | Tipo      | Descrição                              |
|----------------|-----------|----------------------------------------|
| `id`            | UUID      | Chave primária                         |
| `name`          | TEXT      | Nome do usuário                        |
| `email`         | TEXT      | Email único (usado no login)           |
| `password_hash` | TEXT      | Hash bcrypt da senha                   |

### Tabela `tasks`

| Coluna        | Tipo        | Descrição                                      |
|--------------|-------------|------------------------------------------------|
| `id`          | UUID        | Chave primária                                 |
| `title`       | TEXT        | Título da tarefa                               |
| `description` | TEXT        | Descrição detalhada                            |
| `status`      | INT         | 0=Pendente, 1=Em andamento, 2=Concluída        |
| `priority`    | INT         | 0=Baixa, 1=Média, 2=Alta                       |
| `due_date`    | TIMESTAMPTZ | Prazo (opcional)                               |
| `created_by`  | UUID        | FK → usuário que criou                         |
| `assigned_to` | UUID        | FK → usuário atribuído (opcional)              |
| `created_at`  | TIMESTAMPTZ | Data/hora de criação                           |

---

## Endpoints

### Autenticação (pública)

| Método | Rota           | Descrição                        |
|--------|----------------|----------------------------------|
| POST   | `/auth/login`  | Login — retorna token JWT        |
| POST   | `/auth/logout` | Logout (cliente descarta token)  |

**Exemplo de login:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "usuario@email.com", "password": "senha123"}'
```
Resposta:
```json
{"token": "eyJhbGciOiJIUzI1NiJ9..."}
```

### Usuários (requer token)

Adicione o header `Authorization: Bearer <token>` em todas as requisições abaixo.

| Método | Rota           | Descrição               |
|--------|----------------|-------------------------|
| POST   | `/users/`      | Criar usuário           |
| GET    | `/users/{id}`  | Buscar usuário por ID   |
| PUT    | `/users/{id}`  | Atualizar usuário       |
| DELETE | `/users/{id}`  | Remover usuário         |

**Criar usuário:**
```bash
curl -X POST http://localhost:8080/users/ \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Eduardo", "email": "edu@email.com", "password": "senha123"}'
```

### Tarefas (requer token)

| Método | Rota           | Descrição                            |
|--------|----------------|--------------------------------------|
| POST   | `/tasks/`      | Criar tarefa                         |
| GET    | `/tasks/{id}`  | Buscar tarefa por ID                 |
| GET    | `/tasks/`      | Listar tarefas com filtros opcionais |
| PUT    | `/tasks/{id}`  | Atualizar tarefa                     |
| DELETE | `/tasks/{id}`  | Remover tarefa                       |

**Filtros disponíveis em `GET /tasks/`:**
- `assignedTo` — UUID do usuário atribuído
- `createdBy` — UUID do criador
- `status` — 0 (Pendente), 1 (Em andamento), 2 (Concluída)
- `priority` — 0 (Baixa), 1 (Média), 2 (Alta)
- `dueBefore` — Data no formato `YYYY-MM-DD`

**Exemplo com filtros:**
```bash
GET /tasks/?status=0&priority=2&dueBefore=2025-12-31
```

---

## Requisitos Complementares

### Filtro Avançado de Tarefas

O endpoint `GET /tasks/` suporta múltiplos parâmetros de filtro combinados, permitindo buscas precisas como "todas as tarefas de alta prioridade ainda pendentes com prazo até junho".

### Webhooks (integração com n8n)

Quando uma tarefa é **criada**, **atribuída** ou **concluída**, a API envia automaticamente um POST HTTP para a URL configurada em `WEBHOOK_URL` no arquivo `.env`.

**Payload enviado:**
```json
{
  "event": "task.created",
  "timestamp": "2025-06-01T10:00:00Z",
  "task": {
    "id": "...",
    "title": "Implementar login",
    "status": 0,
    "priority": 2
  }
}
```

**Como integrar com n8n:**
1. No n8n, crie um workflow com um node **Webhook** (método POST)
2. Copie a URL gerada pelo n8n
3. No arquivo `.env`, configure: `WEBHOOK_URL=https://sua-url-n8n.com/webhook/...`
4. Reinicie a API — agora todas as mudanças de tarefas chegam ao n8n

---

## Configuração e Execução

### Pré-requisitos

- Go 1.21+
- PostgreSQL 14+
- pgAdmin (opcional, para visualizar o banco)

### 1. Configurar o banco de dados

No pgAdmin, crie um banco chamado `tarefas` e execute as migrations na ordem:

```sql
-- Execute cada arquivo em migrations/ na ordem numérica
migrations/000001_create_users.up.sql
migrations/000002_create_tasks.up.sql
migrations/000003_add_user_auth_fields.up.sql
migrations/000004_add_task_priority_duedate.up.sql
```

### 2. Configurar variáveis de ambiente

Edite o arquivo `.env`:

```env
DATABASE_DSN=postgres://usuario:senha@localhost:5432/tarefas?sslmode=disable
JWT_SECRET=sua-chave-secreta-aqui
WEBHOOK_URL=             # opcional — URL do n8n ou outro sistema
```

### 3. Instalar dependências

```bash
go mod vendor
```

### 4. Rodar a API

```bash
go run ./cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`.  
A documentação Swagger estará em `http://localhost:8080/swagger`.

---

## Testes Automatizados

Os testes cobrem as camadas de serviço (application layer) usando repositórios mock — sem precisar de banco de dados real.

**Executar todos os testes:**
```bash
go test ./...
```

**Com cobertura de código:**
```bash
go test ./... -cover
```

**Testes por módulo:**
```bash
go test ./internal/application/user/...
go test ./internal/application/task/...
go test ./internal/application/auth/...
```

### Estratégia de Testes

- **Testes unitários** — testam cada serviço isoladamente com mocks
- **Mocks** — implementações falsas dos repositórios, sem banco de dados
- **Casos testados:** criação, busca, atualização, exclusão, filtros, autenticação, tokens JWT

---

## Boas Práticas Aplicadas

- **Logs estruturados** com `log/slog` (padrão Go 1.21+)
- **Tratamento de erros** com `errors.Join` para rastreabilidade
- **Senhas** armazenadas como hash bcrypt — nunca em texto puro
- **JWT stateless** — sem sessão no servidor
- **Webhooks assíncronos** — executados em goroutine para não bloquear a API
- **Migrations versionadas** — controle de versão do banco de dados
- **Vendor** — dependências vendorizadas para builds reproduzíveis

---

## Estrutura do Projeto

```
.
├── cmd/api/
│   ├── handler/        ← Handlers HTTP (users, tasks, auth)
│   ├── middleware/     ← Middleware de autenticação JWT
│   └── main.go         ← Ponto de entrada da aplicação
├── internal/
│   ├── application/    ← Serviços (lógica de negócio)
│   ├── domain/         ← Entidades e interfaces
│   └── infrastructure/ ← Implementações do banco de dados
├── migrations/         ← Scripts SQL de criação/alteração do banco
├── vendor/             ← Dependências locais
├── .env                ← Configurações (não commitar em produção)
└── go.mod              ← Módulo Go
```
