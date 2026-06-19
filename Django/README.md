# Encurtador de URL — Django + DRF

API de encurtamento de URLs desenvolvida em Python com **Django** e **Django REST Framework**, utilizando PostgreSQL como banco de dados e cache em memória com `threading.Lock`. Todas as rotas são protegidas por autenticação **Basic Auth**.

## Tecnologias

| Componente   | Tecnologia                                    |
|--------------|-----------------------------------------------|
| Framework    | Django 5 + Django REST Framework              |
| Servidor     | Gunicorn / runserver                          |
| Banco        | PostgreSQL (psycopg2-binary)                  |
| Cache        | dict + threading.Lock                         |
| Autenticação | Basic Auth custom (secrets.compare_digest)    |
| ID           | secrets.choice (8 chars, a-zA-Z0-9)           |

## Variáveis de Ambiente

| Variável       | Descrição                         | Exemplo              |
|----------------|-----------------------------------|----------------------|
| `DB_TYPE`      | Tipo do banco de dados            | `postgres`           |
| `DB_NAME`      | Nome do banco                     | `encurtador`         |
| `DB_HOST`      | Host do banco                     | `0.0.0.0`            |
| `DB_PORT`      | Porta do banco                    | `5432`               |
| `DB_USER`      | Usuário do banco                  | `postgres`           |
| `DB_PASS`      | Senha do banco                    | `postgres`           |
| `HTTP_BASE`    | Prefixo base das rotas            | `/api/`              |
| `HTTP_PORT`    | Porta do servidor                 | `6060`               |
| `TIMEOUT_TIME` | Timeout de leitura (segundos)     | `3`                  |
| `USER`         | Usuário para autenticação         | `user`               |
| `PASS`         | Senha para autenticação           | `pass123`            |

## Autenticação

Todas as rotas utilizam **Basic Auth**. O header `Authorization` deve ser enviado com as credenciais codificadas em Base64 no formato `user:pass`.

```
Authorization: Basic <base64(user:pass)>
```

Exemplo com as credenciais padrão (`user:pass123`):
```
Authorization: Basic dXNlcjpwYXNzMTIz
```

---

## Rotas

Considerando `HTTP_BASE=/api/` e `HTTP_PORT=6060`, o base URL é `http://localhost:6060/api/`.

### `GET /urls`

Retorna todas as URLs encurtadas cadastradas.

**Response** — `200 OK`

```json
{
  "abc12XYZ": {
    "original": "https://www.example.com/pagina-muito-longa",
    "accesses": 5
  },
  "def34ABC": {
    "original": "https://www.google.com",
    "accesses": 12
  }
}
```

---

### `POST /urls`

Cria uma nova URL encurtada.

**Request Body**

```json
{
  "url": "https://www.example.com/pagina-muito-longa"
}
```

| Campo | Tipo   | Obrigatório | Descrição                  |
|-------|--------|-------------|----------------------------|
| `url` | string | Sim         | URL original a ser encurtada |

**Response** — `201 Created` (nova URL criada)

```json
{
  "id": "abc12XYZ",
  "url": "https://www.example.com/pagina-muito-longa"
}
```

**Response** — `200 OK` (URL já existente, retorna o ID já cadastrado)

```json
{
  "id": "abc12XYZ",
  "url": "https://www.example.com/pagina-muito-longa"
}
```

**Response** — `400 Bad Request` (URL inválida ou body malformado)

---

### `GET /urls/{id}`

Retorna informações de uma URL encurtada específica pelo seu ID.

**Path Parameters**

| Parâmetro | Tipo   | Obrigatório | Descrição        |
|-----------|--------|-------------|------------------|
| `id`      | string | Sim         | ID da URL curta  |

**Response** — `200 OK`

```json
{
  "original": "https://www.example.com/pagina-muito-longa",
  "accesses": 5
}
```

**Response** — `400 Bad Request` (ID não encontrado)

---

### `GET /{id}`

Redireciona para a URL original correspondente ao ID informado. Incrementa o contador de acessos.

**Path Parameters**

| Parâmetro | Tipo   | Obrigatório | Descrição        |
|-----------|--------|-------------|------------------|
| `id`      | string | Sim         | ID da URL curta  |

**Response** — `302 Found`

Redirecionamento HTTP para a URL original salva.

**Response** — `404 Not Found` (ID não encontrado)

---

### `DELETE /{id}`

Deleta uma URL encurtada pelo seu ID.

**Path Parameters**

| Parâmetro | Tipo   | Obrigatório | Descrição        |
|-----------|--------|-------------|------------------|
| `id`      | string | Sim         | ID da URL curta  |

**Response** — `200 OK` (deleção bem-sucedida, corpo vazio)

**Response** — `400 Bad Request` (ID não encontrado)

---

## Executando o Projeto

```bash
pip install -r requirements.txt
bash scripts/main.sh
```

Ou manualmente:

```bash
export HTTP_PORT=6060 HTTP_BASE=/api/ DB_HOST=... # demais variáveis
python manage.py makemigrations api
python manage.py migrate
python manage.py runserver 0.0.0.0:$HTTP_PORT
```
