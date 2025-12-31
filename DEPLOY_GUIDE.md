# 🚀 GUÍA DE DESPLIEGUE - URL SHORTENER

## Stack de Producción

```
┌─────────────────────────────────────────────────────────────┐
│                    ARQUITECTURA CLOUD                        │
└─────────────────────────────────────────────────────────────┘

     Usuario
        │
        ▼
  ┌───────────┐
  │  VERCEL   │ ◄── Frontend React (CDN Global)
  │  (Free)   │     https://url-shortener.vercel.app
  └─────┬─────┘
        │ Proxy /api/* y /s/*
        ▼
  ┌───────────┐
  │  FLY.IO   │ ◄── Backend Go (Auto-scaling)
  │  (Free)   │     https://url-shortener-api.fly.dev
  └─────┬─────┘
        │
   ┌────┴────┐
   ▼         ▼
┌──────┐  ┌──────┐
│UPSTASH│ │ NEON │
│Redis │  │Postgres│
│(Free)│  │(Free) │
└──────┘  └──────┘
```

---

## 📋 CHECKLIST DE CONFIGURACIÓN

### ☐ Paso 1: Crear cuenta en los servicios (5 min cada uno)

- [ ] **Neon** → https://neon.tech (PostgreSQL)
- [ ] **Upstash** → https://upstash.com (Redis)
- [ ] **Fly.io** → https://fly.io (Backend)
- [ ] **Vercel** → https://vercel.com (Frontend)

---

## 1️⃣ NEON (PostgreSQL)

### Crear base de datos:

1. Ir a https://console.neon.tech
2. Click en **"New Project"**
3. Configurar:
   - **Name:** `url-shortener`
   - **Region:** `US East (Ohio)` (o el más cercano)
   - **Postgres Version:** `15`
4. Click **"Create Project"**

### Obtener credenciales:

Después de crear, verás algo como:

```
postgres://username:password@ep-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require
```

### ☐ Guardar estos valores:

| Variable | Valor |
|----------|-------|
| `DATABASE_URL` | `postgres://user:pass@ep-xxx.aws.neon.tech/neondb?sslmode=require` |

---

## 2️⃣ UPSTASH (Redis)

### Crear base de datos Redis:

1. Ir a https://console.upstash.com
2. Click en **"Create Database"**
3. Configurar:
   - **Name:** `url-shortener-cache`
   - **Region:** `US-East-1` (Virginia)
   - **Type:** `Regional` (gratis)
4. Click **"Create"**

### Obtener credenciales:

En la página de la base de datos, buscar:
- **UPSTASH_REDIS_REST_URL**
- **UPSTASH_REDIS_REST_TOKEN**

O para conexión tradicional:
- **Endpoint:** `xxx.upstash.io`
- **Port:** `6379`
- **Password:** `xxxxx`

### ☐ Guardar estos valores:

| Variable | Valor |
|----------|-------|
| `REDIS_URL` | `xxx.upstash.io:6379` |
| `REDIS_PASSWORD` | `tu-password-aqui` |

---

## 3️⃣ FLY.IO (Backend)

### Instalar Fly CLI:

```bash
# Windows (PowerShell como Admin)
powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"

# Mac/Linux
curl -L https://fly.io/install.sh | sh
```

### Login y crear app:

```bash
# Login (abre navegador)
fly auth login

# Ir al directorio del backend
cd backend

# Crear la app (solo la primera vez)
fly apps create url-shortener-api

# Verificar que se creó
fly apps list
```

### Configurar secretos:

```bash
# Configurar variables de entorno en Fly.io
fly secrets set DATABASE_URL="postgres://user:pass@ep-xxx.aws.neon.tech/neondb?sslmode=require"
fly secrets set REDIS_URL="xxx.upstash.io:6379"
fly secrets set REDIS_PASSWORD="tu-password-de-upstash"
fly secrets set BASE_URL="https://url-shortener.vercel.app"
fly secrets set GIN_MODE="release"

# Verificar que se guardaron
fly secrets list
```

### Desplegar:

```bash
# Desde el directorio backend/
fly deploy

# Ver logs en tiempo real
fly logs

# Verificar status
fly status
```

### ☐ Verificar que funciona:

```bash
# Health check
curl https://url-shortener-api.fly.dev/health
# Debe retornar: {"status":"healthy",...}
```

### Obtener token para CI/CD:

```bash
# Generar token para GitHub Actions
fly tokens create deploy -x 999999h

# Copiar el token que aparece (empieza con "FlyV1...")
```

---

## 4️⃣ VERCEL (Frontend)

### Opción A: Deploy automático (Recomendado)

1. Ir a https://vercel.com/new
2. Click **"Import Git Repository"**
3. Seleccionar tu repo de GitHub
4. Configurar:
   - **Framework Preset:** `Vite`
   - **Root Directory:** `frontend`
   - **Build Command:** `npm run build`
   - **Output Directory:** `dist`
5. Agregar variable de entorno:
   - `VITE_API_URL` = `https://url-shortener-api.fly.dev/api`
6. Click **"Deploy"**

### Opción B: Deploy con CLI

```bash
# Instalar Vercel CLI
npm i -g vercel

# Login
vercel login

# Ir al directorio frontend
cd frontend

# Deploy (sigue las instrucciones)
vercel

# Deploy a producción
vercel --prod
```

### Obtener IDs para CI/CD:

1. Ir a **Project Settings** en Vercel
2. En **General**, copiar:
   - **Project ID**
3. Ir a tu perfil → **Settings** → **Tokens**
4. Crear token con scope de deploy

---

## 5️⃣ GITHUB SECRETS

### Agregar secretos en GitHub:

1. Ir a tu repo en GitHub
2. **Settings** → **Secrets and variables** → **Actions**
3. Click **"New repository secret"**

### ☐ Secretos requeridos:

| Nombre | Valor | Para |
|--------|-------|------|
| `FLY_API_TOKEN` | `FlyV1...` (del paso 3) | Deploy a Fly.io |
| `VERCEL_TOKEN` | Token de Vercel | Deploy a Vercel |
| `VERCEL_ORG_ID` | `team_ZnGmSWKKAwiEW6XAQdI7KVAZ` | Deploy a Vercel |
| `VERCEL_PROJECT_ID` | `prj_oA4P727s7zerwZ0slc1DM8HguRnK` | Deploy a Vercel |
| `DATABASE_URL` | URL de Neon | Tests de integración |

---

## 6️⃣ VERIFICACIÓN FINAL

### ☐ Checklist de verificación:

1. [ ] **Backend funciona:**
   ```bash
   curl https://url-shortener-api.fly.dev/health
   ```

2. [ ] **Frontend funciona:**
   - Visitar https://url-shortener.vercel.app

3. [ ] **Crear URL funciona:**
   ```bash
   curl -X POST https://url-shortener-api.fly.dev/api/shorten \
     -H "Content-Type: application/json" \
     -d '{"url":"https://github.com"}'
   ```

4. [ ] **Redirección funciona:**
   - Visitar la URL corta generada

5. [ ] **CI/CD funciona:**
   - Hacer un push a main y verificar que se despliega

---

## 📊 COSTOS ESTIMADOS

| Servicio | Tier Gratuito | Límites |
|----------|---------------|---------|
| **Neon** | ✅ Free | 0.5 GB storage, 1 proyecto |
| **Upstash** | ✅ Free | 10K comandos/día, 256MB |
| **Fly.io** | ✅ Free | 3 VMs shared, 160GB transfer |
| **Vercel** | ✅ Free | 100GB bandwidth, unlimited deploys |

**Total: $0/mes** (dentro de los límites gratuitos)

---

## 🔧 COMANDOS ÚTILES

### Fly.io:
```bash
fly status              # Ver estado de la app
fly logs                # Ver logs en tiempo real
fly logs -a url-shortener-api  # Logs de una app específica
fly ssh console         # SSH al contenedor
fly scale count 2       # Escalar a 2 instancias
fly secrets list        # Ver secretos configurados
fly open                # Abrir app en navegador
```

### Vercel:
```bash
vercel                  # Deploy preview
vercel --prod           # Deploy producción
vercel logs             # Ver logs
vercel env pull         # Descargar variables de entorno
```

### Neon:
```bash
# Conectar con psql
psql "postgres://user:pass@ep-xxx.aws.neon.tech/neondb?sslmode=require"

# Ver tablas
\dt

# Ver datos
SELECT * FROM urls LIMIT 10;
```

---

## 🚨 TROUBLESHOOTING

### Error: "No machines in group app"
```bash
fly scale count 1
```

### Error: "Connection refused" en PostgreSQL
- Verificar que el SSL está habilitado: `?sslmode=require`
- Verificar que la IP de Fly.io está permitida en Neon

### Error: "Redis connection failed"
- Verificar password de Upstash
- Usar formato: `rediss://default:PASSWORD@HOST:PORT` para TLS

### Frontend no conecta con backend
- Verificar `VITE_API_URL` en Vercel
- Verificar rewrites en `vercel.json`

---

## ✅ ¡LISTO!

Una vez completados todos los pasos, tu aplicación estará disponible en:

- **Frontend:** https://url-shortener.vercel.app
- **Backend API:** https://url-shortener-api.fly.dev
- **Health Check:** https://url-shortener-api.fly.dev/health

Cada push a `main` desplegará automáticamente:
- Cambios en `backend/` → Fly.io
- Cambios en `frontend/` → Vercel
