# Mallki — ecommerce mayorista de nueces

Catálogo, armado de pedidos y listado de pedidos para el emprendimiento de nueces de mi papá.
Backend en Go, frontend en React + Vite, PostgreSQL como base de datos — todo corre en contenedores.

## Arrancar desde cero

Requisitos: [Docker Desktop](https://docs.docker.com/get-docker/) instalado y corriendo.

1. Cloná el repositorio:
   ```
   git clone https://github.com/magdalenagomezf/ingsoft3-tp01.git
   cd ingsoft3-tp01
   ```

2. Copiá la plantilla de variables de entorno y poné una contraseña:
   ```
   cp .env.example .env
   ```
   Editá `.env` y cambiá el valor de `DB_PASSWORD` por cualquier contraseña.

3. Levantá todo el sistema:
   ```
   docker compose up -d --build
   ```

4. Esperá a que los tres servicios estén listos:
   ```
   docker compose ps
   ```
   `db` debe figurar como `healthy`.

5. Abrí `http://localhost:3000` en el navegador.

## Variante: usando las imágenes publicadas (sin compilar)

En vez del paso 3, con el mismo `.env` ya creado:
```
docker compose -f docker-compose.registry.yml up -d
```
Esto descarga las imágenes ya publicadas en ghcr.io en vez de construirlas localmente.

## Apagar el sistema

```
docker compose down        # apaga, conserva los datos
docker compose down -v     # apaga y borra también los datos (el volumen)
```

## Stack

- **Backend**: Go, arquitectura en capas (`handler` → `service` → `repository`), PostgreSQL vía `lib/pq`
- **Frontend**: React + Vite + Tailwind, servido por nginx en producción
- **Base de datos**: PostgreSQL 16
- **Contenedores**: Dockerfiles multi-stage para backend y frontend, orquestados con Docker Compose