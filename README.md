# 🔗 URL Shortener Escalable

[![CI/CD Pipeline](https://github.com/username/url-shortener/actions/workflows/pipeline.yml/badge.svg)](https://github.com/username/url-shortener/actions)
[![Deploy to Fly.io](https://github.com/username/url-shortener/actions/workflows/deploy-fly.yml/badge.svg)](https://github.com/username/url-shortener/actions)
[![Docker Image Size](https://img.shields.io/badge/docker%20image-13.8MB-brightgreen)](https://hub.docker.com)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)

> Un acortador de URLs de alto rendimiento diseñado para demostrar dominio en FullStack y DevOps moderno.

## 🌐 Demo en Producción

| Servicio | URL |
|----------|-----|
| **Frontend** | https://url-shortener-xtaxx12.vercel.app |
| **Backend API** | https://url-short-sebas-api.fly.dev |
| **Auditoría** | [Ver Reporte de Calidad ✅](./REVISIÓN_FINAL.md) |
| **Health Check** | https://url-short-sebas-api.fly.dev/health |

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────────────────┐
│                         NGINX (Load Balancer)                    │
│                      Round Robin Distribution                    │
└─────────────────────┬───────────────────────┬───────────────────┘
                      │                       │
         ┌────────────▼────────────┐ ┌────────▼────────────┐
         │   Go API (Replica 1)    │ │   Go API (Replica 2) │
         │      Port 8081          │ │      Port 8082       │
         └────────────┬────────────┘ └────────┬────────────┘
                      │                       │
         ┌────────────▼───────────────────────▼────────────┐
         │                    Redis Cache                   │
         │              (Fast Redirections)                 │
         └────────────────────────┬────────────────────────┘
                                  │
         ┌────────────────────────▼────────────────────────┐
         │                   PostgreSQL                     │
         │              (Persistent Storage)                │
         └─────────────────────────────────────────────────┘
```

## 🛠️ Stack Tecnológico

| Componente | Tecnología |
|------------|------------|
| **Backend** | Go 1.21+ con Gin Framework |
| **Frontend** | React 18 + Vite + TailwindCSS |
| **Base de Datos** | PostgreSQL 15 |
| **Caché** | Redis 7 |
| **Load Balancer** | Nginx |
| **Containerización** | Docker (Multi-stage builds) |
| **CI/CD** | GitHub Actions |

## 📁 Estructura del Proyecto

```
url-shortener/
├── backend/
│   ├── cmd/api/
│   │   └── main.go              # Entry point con graceful shutdown
│   ├── internal/
│   │   ├── domain/url.go        # Entidades y Ports (interfaces)
│   │   ├── handler/             # HTTP handlers + tests
│   │   ├── repository/          # Adapters: PostgreSQL + Redis
│   │   └── service/             # Business logic (caché-first)
│   ├── pkg/shortener/           # Generador de códigos + tests
│   ├── Dockerfile               # Multi-stage (scratch) <50MB
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── App.jsx              # UI moderna con glassmorphism
│   │   └── index.css            # Estilos premium
│   ├── Dockerfile               # Multi-stage con NGINX
│   └── package.json
├── infra/
│   ├── nginx/nginx.conf         # Load Balancer Round Robin
│   ├── terraform/               # (Opcional) IaC para AWS
│   └── scripts/
│       ├── deploy.sh            # Script de despliegue
│       └── init-db.sql          # Esquema inicial
├── .github/workflows/
│   └── pipeline.yml             # CI/CD completo
├── docker-compose.yml           # Stack completo local
└── README.md
```

## 🚀 Quick Start

### Requisitos Previos
- Docker & Docker Compose
- Go 1.21+ (para desarrollo local)
- Node.js 18+ (para desarrollo frontend)

### Levantar el entorno completo

```bash
# Clonar el repositorio
git clone https://github.com/username/url-shortener.git
cd url-shortener

# Levantar todos los servicios
docker-compose up -d

# Verificar que todo esté corriendo
docker-compose ps
```

### Endpoints disponibles

| Servicio | URL |
|----------|-----|
| Frontend | http://localhost:3000 |
| API (via Nginx) | http://localhost:80/api |
| API Directa (Replica 1) | http://localhost:8081 |
| API Directa (Replica 2) | http://localhost:8082 |

## 📊 API Endpoints

### Crear URL corta
```bash
POST /api/shorten
Content-Type: application/json

{
  "url": "https://example.com/muy-larga-url"
}

# Response
{
  "short_code": "abc123",
  "short_url": "http://localhost/abc123",
  "original_url": "https://example.com/muy-larga-url"
}
```

### Redireccionar
```bash
GET /:code

# Redirige a la URL original (HTTP 301)
```

### Obtener estadísticas
```bash
GET /api/stats/:code

# Response
{
  "short_code": "abc123",
  "original_url": "https://example.com/muy-larga-url",
  "clicks": 42,
  "created_at": "2024-01-15T10:30:00Z"
}
```

## 🐳 Docker

### Tamaño de imagen optimizado

```bash
# Construir imagen del backend
docker build -t url-shortener-api ./backend

# Verificar tamaño (objetivo: <50MB)
docker images url-shortener-api
```

La imagen utiliza **multi-stage builds** con `scratch` base para lograr un tamaño mínimo.

## ☁️ Despliegue en Producción (Opcional)

> **Nota:** Los archivos de Terraform están incluidos en `infra/terraform/` para uso futuro cuando se requiera desplegar en AWS.

### Despliegue con Docker (Recomendado para desarrollo)

```bash
# Levantar todo el stack localmente
docker-compose up -d --build

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down
```

### Despliegue en VPS/Servidor

```bash
# En tu servidor (con Docker instalado)
git clone https://github.com/username/url-shortener.git
cd url-shortener
docker-compose up -d
```

<details>
<summary>📦 Terraform para AWS (Uso Futuro)</summary>

Los archivos de Terraform en `infra/terraform/` permiten desplegar la infraestructura en AWS Free Tier:

```bash
cd infra/terraform
terraform init
terraform plan
terraform apply
```

Recursos que crea:
- EC2 t2.micro (Free Tier)
- Security Group (SSH, HTTP, HTTPS)
- User Data para instalar Docker

</details>

## 🔄 CI/CD Pipeline

El pipeline de GitHub Actions ejecuta:

1. ✅ **Test** - Unit tests de Go
2. 🐳 **Build** - Construcción de imagen Docker
3. 📦 **Push** - Push a Docker Hub
4. 🚀 **Deploy** - Despliegue simulado

## 📈 Métricas de Rendimiento

- **Latencia de redirección**: <10ms (cache hit)
- **Latencia de creación**: <50ms
- **Throughput**: >10,000 req/s por réplica
- **Tamaño de imagen Docker**: <50MB

## 📝 Licencia

MIT License - Ver [LICENSE](LICENSE) para más detalles.
