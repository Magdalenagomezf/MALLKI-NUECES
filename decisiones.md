# Decisiones — Proyecto IngSoft3

## TP1 — Control de versiones

### Por qué Git no pudo resolver el conflicto solo
Las ramas `feature/titulo-a` y `feature/titulo-b` partieron las dos del mismo commit de `main`, 
y ambas modificaron la misma línea del `README.md`, cada una con un texto distinto
("versión A" vs "versión B"). El merge de tres vías que usa Git compara el commit común (la base),
la rama A y la rama B — y solo puede fusionar automático cuando los cambios no se pisan. Como las
dos tocaron exactamente la misma línea con contenido distinto, Git no tiene forma de decidir y necesita que una persona lo resuelva.

La rama B tendría que haber traído los cambios de `main` (ya con A mergeada) antes de abrir su
 propio PR — así el conflicto aparecía localmente, antes de llegar a GitHub, sabiendo de antemano qué había en A.

### Problemas encontrados y cómo los resolví
Para que el conflicto entre las ramas A y B apareciera de verdad (y no se mergeara solo), las dos
  tenían que partir del mismo commit de `main` y modificar la misma línea del README. Si una rama
  partía de la otra en vez de `main`, Git no generaba ningún conflicto — hubo que prestar atención
  al orden: crear ambas ramas desde `main`, y recién mergear la primera después de tener las dos
  abiertas.

### Uso de IA
Usé Claude como guía durante todo el TP1: para entender la protección de rama, interpretar los
marcadores de conflicto antes de resolverlos, redactar la descripción de la release, y resolver
dudas puntuales. Verifiqué cada paso ejecutándolo yo misma y mirando el resultado real en GitHub.

## TP2 — Contenedores

### Qué app elegí y por qué
Mallki Nueces: un ecommerce mayorista de nueces para el trabajo de mi papá. Es una idea que ya tenia pendiente y decidi aprovechar la oportunidad de hacerla, ademas cumple los criterios de la guía: backend con API (Go), frontend SPA (React + Vite), 
base de datos relacional (PostgreSQL), corre local sin nada raro, y el tamaño es chico a propósito — 
catálogo, armado de pedido y listado
de pedidos, sin login ni pagos (eso queda para más adelante, fuera del alcance de este TP).

### Decisiones de contenerización
- **Imágenes base**: `golang:1.24-alpine` para compilar el backend (coincide con la versión de mi
  `go.mod`), y `alpine:3.20` para correrlo — como el binario de Go es estático
  (`CGO_ENABLED=0`), no necesita ningún runtime de Go instalado en la imagen final. Para el
  frontend, `node:22-alpine` compila el build de Vite, y `nginx:alpine` sirve los estáticos
  resultantes.
- **Multi-stage en los dos**: cada Dockerfile tiene una etapa de compilación (pesada, con el SDK/
  toolchain) y una etapa final mínima que solo recibe el resultado ya compilado —
  `golang:1.24-alpine` pesa 83.5MB, mi imagen final del backend, 9.11MB.
- **Qué persiste y qué no**: solo los datos de Postgres, en el volumen nombrado `db_data`. Backend
  y frontend son efímeros a propósito — su código vive copiado *dentro* de la imagen (no montado
  como volumen), así que cualquier cambio requiere reconstruir la imagen, no editar un archivo
  vivo. Confirmé esto con la prueba de `down`/`up` vs `down -v` en `evidencias.md`.
- **Red y nombres**: dentro de compose, el backend le habla a la base por el nombre del servicio
  (`db:5432`, el puerto interno de Postgres) — no por `localhost` ni por el puerto publicado hacia
  afuera, que es un número distinto y solo le sirve a mi máquina, no a otros contenedores.

### Problemas encontrados y cómo los resolví
- Al probar `docker-compose.registry.yml` (que baja las imágenes en vez de compilarlas), corrimos
  `docker builder prune -af` para forzar una descarga limpia de verdad — eso también me borró las
  imágenes base (`golang:1.24-alpine`, `node:22-alpine`, etc.) que necesitaba para la comparación de
  tamaños de `evidencias.md`. Tuve que volver a bajarlas con `docker pull` una por una.

### Uso de IA
Usé Claude Code para escribir el scaffolding inicial del backend (arquitectura en capas) y del
frontend, a partir de un prompt donde definí el alcance, el modelo de datos y las restricciones
explícitas (sin pagos, sin login, sin Docker generado por la IA). Usé Claude para entender la
teoría de contenedores (namespaces, layers, multi-stage, redes de compose) y para guiarme paso a
paso escribiendo los Dockerfiles, el `docker-compose.yml` y el `docker-compose.registry.yml` — pero
cada comando lo corrí yo misma, y verifiqué el resultado real en mi terminal y en el navegador. Cuando algo no
coincidía con lo esperado, lo diagnostiqué con mis propios comandos antes de pedir ayuda para interpretarlo.

## TP3 — Planificación

**Duración del sprint:** 1 semana. Elegí ese numero porque coincide con el ritmo real
de entrega de la materia (un TP por clase), alinear el sprint con eso hace que cada clase sea el cierre de un sprint, y refleja lo que de verdad estoy entregando, no un número arbitrario.

**Limite de trabajo en progreso:** 2. Es la regla de arranque de la guia (cantidad de personas + 1);
trabajando sola, 1 + 1 = 2. La herramienta no bloquea que sigas agregando tarjetas, pero el contador se pone en rojo cuando te pasás del límite — así el problema se nota en vez de acumularse en silencio.
Si con el tiempo, nunca alcanzo ese limite, esta demasiado alto y debo bajarlo.

**Diagnostico de la historia mal escrita:** 
"Como desarrollador quiero crear la tabla
usuarios" mezcla una tarea tecnica (crear una tabla) con el formato de historia.Nadie quiere una tabla en la base de datos, es un detalle de implementacion, no un requisito para alguien. Le falta un beneficio real- el "para que" le importa a un usuario, no a la base de datos.
Cómo la reescribiría: "Como usuario quiero ver mi perfil con los datos que cargué antes, para no tener que completarlos de nuevo cada vez que entro" — con un criterio verificable (por ejemplo: los datos del perfil siguen ahí después de cerrar sesión y volver a entrar). "Crear la tabla usuarios" pasaría a ser una **tarea** técnica dentro de esa historia.

**Problemas que encontre:**
- Al token de gh le faltaba el scope 'project'; se resolvio con `gh auth refresh -s project`.
- cmd.exe no interpreta comillas simples como GitHub CLI espera (ej: `--owner '@me'` tiraba "unknown owner type"); hubo que sacar las comillas o cambiarlas por dobles.
- Se me subio `ci.yml` vacio por no confirmar el contenido antes
  de commitear — lo solucione revisando el diff en "Files changed" del PR antes de  mergear, en vez de mergear a ciegas.
- Las dos tareas quedaron colgando de la épica en vez de la historia — al agregarlas por la web, seguí parada en la página de la épica en vez de ir   a la de la historia (#7).
  Lo noté comparando el tablero contra el esquema de la Clase 3 (épica → historia → tarea)  y lo corregí con `gh issue edit 8 --parent 7` y `gh issue edit 9 --parent 7`.

**Uso de IA:** Usé la IA (Claude) para guiarme paso a paso en toda la configuración de GitHub Projects, traduciendo los comandos del video del profesor a cmd.exe de Windows, para pensar variantes de la historia mal escrita hasta llegar a una que cumpliera los cuatro criterios, y para diagnosticar dos problemas reales: el token de `gh` sin el scope `project`, y el archivo `ci.yml` que se subió vacío dos veces seguidas por no guardarlo antes de comitear.
Las decisiones (duración del sprint, límite de WIP) las tomé yo. Verifiqué cada paso mirando el estado real en GitHub.

## TP4 — CI: Pipelines as Code

### Estructura elegida del pipeline
Dos jobs, `build-backend` y `build-frontend`, corriendo en paralelo. La separación no es arbitraria: coincide con que desde el TP2 tengo dos Dockerfiles distintos (`backend/Dockerfile` y `frontend/Dockerfile`), así que el pipeline refleja cómo está armada la app de verdad. Van en paralelo porque son builds independientes — uno no necesita el resultado del otro, así que correrlos en serie solo sumaría tiempo de espera sin ninguna ganancia.

### Qué cachea y qué pasa si desaparece
El cache guarda las capas de Docker de cada Dockerfile, en particular las más caras: la descarga de dependencias de Go (`go mod download`) en el backend y la instalación de paquetes (`npm ci`) en el frontend. Lo probé con un commit vacío entre dos corridas del mismo PR, y las 7 capas del backend y las 6 del frontend salieron `CACHED` por completo. Cada job usa su propio `scope` (`backend` y `frontend`) para que no se pisen el cache entre sí — si compartieran uno solo, el último job en terminar sobrescribiría el cache del otro. Si el cache desaparece en algún momento (la plataforma lo puede desalojar), el pipeline sigue funcionando exactamente igual, solo que más lento — no hay ninguna dependencia funcional de que exista, es pura optimización de tiempo.

### Por qué el pipeline construye con mi Dockerfile en vez de compilar por su cuenta
El workflow no tiene ni una línea de Go ni de npm — usa exactamente los mismos `backend/Dockerfile` y `frontend/Dockerfile` que ya tenía del TP2. Si el pipeline compilara por su cuenta (con `go build` o `npm run build` directo en el YAML), tendría dos definiciones distintas de cómo se construye la app: la que usa el pipeline para verificar, y la que después uso para desplegar. Tarde o temprano esas dos definiciones divergen, y terminaría verificando algo que no es lo que realmente corre.

### Problemas encontrados y cómo los resolví
- Para la demostración de romper el build a propósito, agregué un import de un paquete que no existe en `main.go`. VS Code tiene activado el formateador automático de Go (`goimports`), que borra los imports no usados apenas guardás el archivo — así que mi import "roto" desaparecía solo antes de poder commitearlo. Lo resolví usando un import en blanco (`_ "nueces-backend/noexiste"`), que el formateador no toca porque está marcado explícitamente como intencional.
- Al mergear con squash el PR que rompía y arreglaba el build (dos commits en la misma rama), la pestaña "Files changed" mostró "No changes to show" — al principio pensé que algo se había perdido, pero es esperable: el resultado neto entre los dos commits, comparado contra `main`, es cero cambios (agregué una línea y la misma línea la saqué después).

### Uso de IA
Usé Claude para traducir el video y la guía del profesor (escritos sobre .NET) a mi stack en Go, y para guiarme paso a paso. También me ayudó a diagnosticar el problema del formateador de Go que borraba el import roto. Cada paso lo corrí yo misma y verifiqué el resultado real en GitHub (los checks, el cache en el log, el badge en el README) antes de seguir.