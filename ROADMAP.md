# Roadmap — Mallki Nueces

Lista de funcionalidades pendientes para el negocio real, priorizadas y con las decisiones
tomadas para cada una. La idea es ir tachando de a una, no todas juntas.

> Nota: este archivo es del producto real, a diferencia de `decisiones.md` / `evidencias.md`
> (esos son artefactos académicos del TP y no se actualizan más). Ver `CLAUDE.md` para el
> contexto completo del proyecto.

## En curso

_(vacío por ahora)_

## Backlog (a priorizar juntos, todavía sin orden ni decisiones)

- **Estados de pedido**: pasar de solo `"pendiente"` a confirmado / en preparación / entregado.
  El campo `pedidos.estado` ya existe en el modelo, pensado para esto.
- **Precios por volumen**: descuento por cantidad, típico de venta mayorista.
- **Notificación de pedido nuevo por WhatsApp**: más simple y más usado acá que email.
- **Carga de fotos reales de producto**: panel simple para reemplazar las fotos genéricas
  actuales — depende del panel admin del login.

## Hecho

### Método de pago por pedido (2026-09-03)
El cliente ahora elige el método de pago al armar el pedido en `/armar-pedido`, sin pasarela de
pago todavía — solo registro de la elección.

**Cómo quedó:**
- Campo nuevo `pedidos.metodo_pago`, requerido, mismo nivel que `cliente_nombre` /
  `cliente_contacto`. Solo dos valores válidos: `"transferencia"` o `"efectivo"` — cualquier otro
  valor (incluido vacío) se rechaza con `400` y mensaje en español.
- `POST /api/pedidos` sigue siendo público; ahí se valida el campo antes de crear el pedido.
- Sin edición desde el panel admin todavía — se define una vez que exista, se muestra en
  `/pedidos` junto a la fecha del pedido.
- Frontend: selector de radio buttons ("Transferencia" / "Efectivo") en el formulario de
  `/armar-pedido`, obligatorio igual que los otros campos del cliente.

**Ojo con la migración:** como no hay herramienta de migraciones real, el esquema se aplica con
`CREATE TABLE IF NOT EXISTS`, que es un no-op si la tabla `pedidos` ya existe. Para que una base
ya corriendo (Docker local con pedidos de prueba) también reciba la columna nueva sin pisar datos
ni requerir un wipe del volumen, se agregó un `ALTER TABLE pedidos ADD COLUMN IF NOT EXISTS
metodo_pago TEXT NOT NULL DEFAULT 'efectivo'` idempotente justo después del `CREATE TABLE`.
Verificado contra el Postgres de `docker compose` local: con una fila pre-existente en `pedidos`,
tras reconstruir y reiniciar el backend la columna se agregó y esa fila quedó con
`metodo_pago = 'efectivo'` (el default), y `POST /api/pedidos` rechaza correctamente valores
inválidos o vacíos de `metodo_pago`.

### Login del admin
Antes, `/pedidos` era una ruta pública sin protección — cualquiera con el link veía todos los
pedidos. Ahora requiere sesión.

**Cómo quedó:**
- Un solo usuario admin, credencial (usuario + hash bcrypt) en variables de entorno
  (`ADMIN_USER`, `ADMIN_PASSWORD_HASH`) — sin tabla nueva ni migración.
- Sesión con cookie **httpOnly** (`SameSite=Lax`), guardada en memoria en el backend. No hace
  falta manejar cross-origin: Vite (dev) y nginx (prod) proxean `/api` al backend desde el mismo
  origen del frontend, así que el `fetch` del frontend no necesita `credentials: 'include'`.
- Rate limit de login: 5 intentos fallidos por IP bloquean esa IP 15 minutos (`429`).
- `POST /api/pedidos` (armar pedido) sigue público a propósito — solo `GET /api/pedidos`
  (listar, usado por `/pedidos`) quedó protegido.
- Frontend: pantalla `/login`, link "Admin" en el Header, `/pedidos` redirige a `/login` si el
  fetch da 401, botón "Cerrar sesión".
- Pensado como la puerta de entrada a un futuro **panel admin** más grande (gestión de
  productos, carga de fotos — ver backlog).

**Ojo si cambiás la contraseña más adelante:** un hash bcrypt tiene `$` (`$2a$10$...`). Tanto
`.env` como `docker-compose.yml` interpretan `$algo` como referencia a una variable, así que hay
que escapar cada `$` como `$$` en el `.env` o el hash queda corrupto y ningún login funciona.
Verificado end-to-end con Docker local: login, logout, guard de `/pedidos`, rate limit al 6to
intento, y catálogo/armar pedido siguen públicos sin sesión.
