# Roadmap — Mallki Nueces

Lista de funcionalidades pendientes para el negocio real, priorizadas y con las decisiones
tomadas para cada una. La idea es ir tachando de a una, no todas juntas.

> Nota: este archivo es del producto real, a diferencia de `decisiones.md` / `evidencias.md`
> (esos son artefactos académicos del TP y no se actualizan más). Ver `CLAUDE.md` para el
> contexto completo del proyecto.

## En curso

_(vacío por ahora)_

## Backlog (a priorizar juntos, todavía sin orden ni decisiones)

- **Precios por volumen**: descuento por cantidad, típico de venta mayorista.
- **Carga de fotos reales de producto**: panel simple para reemplazar las fotos genéricas
  actuales — depende del panel admin del login.

## Hecho

### Notificación de pedido nuevo por WhatsApp (2026-09-03)
El admin puede avisarse a sí mismo por WhatsApp cuando entra un pedido nuevo, sin depender de
WhatsApp Business API ni de un número de Twilio (ninguno de los dos está disponible todavía).

**Cómo quedó:**
- Botón "Avisar por WhatsApp" en cada tarjeta de `/pedidos`, que abre un link `wa.me/<numero>`
  con el mensaje del pedido (cliente, contacto, método de pago, items) ya escrito.
- **No es automático**: alguien tiene que tocar "Enviar" en WhatsApp. La alternativa automática
  real (WhatsApp Business API) requiere cuenta verificada y tiene costo por mensaje; las
  librerías no oficiales que automatizan WhatsApp Web violan los términos de uso y arriesgan el
  baneo del número — ninguna de las dos vale la pena para este volumen de pedidos.
- El número del admin está hardcodeado en el frontend (`WHATSAPP_ADMIN` en `Pedidos.jsx`). No es
  un secreto — es el mismo número que un botón público de "chatear por WhatsApp" mostraría en
  cualquier sitio — pero si cambia, hay que tocar código (no hay panel para editarlo).

### Panel de stock editable por el admin (2026-09-03)
Los productos (nombre, descripción, categoría, foto, precio) se cargan directo en la base y no
cambian — pero el stock disponible sí, todos los días.

**Cómo quedó:**
- Nueva página `/stock`, protegida igual que `/pedidos`, con nombre/categoría/precio de solo
  lectura y el stock (`stock_kg`) editable por producto.
- Nuevo `PATCH /api/productos/{id}/stock`, endpoint dedicado (no reusa el `PUT` genérico) para
  no arriesgar pisar otros campos con datos viejos del frontend.
- De paso se protegieron `POST/PUT/DELETE /api/productos`, que estaban públicos sin querer desde
  el día uno — no hay ningún flujo del catálogo que los necesite sin login.

### Nav protegido por sesión (2026-09-03)
El link "Pedidos" del header aparecía siempre, aunque `/pedidos` ya rechazara la carga sin
sesión con un redirect a `/login`. Ahora directamente no se muestra sin sesión activa.

**Cómo quedó:**
- Nuevo `GET /api/me` (protegido) para que el frontend le pregunte al backend si la cookie de
  sesión sigue siendo válida — no se puede leer la cookie httpOnly desde JS.
- `AuthProvider`/`useAuth()` centralizan ese estado; login, logout y una sesión vencida (401 en
  cualquier fetch protegido) lo actualizan al toque, sin esperar un refresh de página.
- El link que decía "Admin" ahora dice "Iniciar sesión".

### Estados de pedido editables por el admin (2026-09-03)
Antes, `pedidos.estado` existía en el modelo pero nadie lo podía cambiar: todo pedido quedaba
`"pendiente"` para siempre. Ahora el admin puede actualizarlo desde `/pedidos`.

**Cómo quedó:**
- Cuatro estados válidos: `pendiente`, `confirmado`, `en_preparacion`, `entregado` — mostrados en
  la UI como "Pendiente", "Confirmado", "En preparación", "Entregado". Cualquier otro valor
  (incluido vacío) se rechaza con `400` y mensaje en español.
- **Transiciones libres, sin máquina de estados**: el admin puede poner cualquier estado válido en
  cualquier momento, sin restricción de orden ni de secuencia. Decisión explícita para permitir
  corregir errores de carga (por ejemplo, volver un pedido de "entregado" a "confirmado" si se
  marcó mal) sin agregar la complejidad de una máquina de estados que todavía no hace falta.
- Endpoint nuevo `PATCH /api/pedidos/{id}/estado`, protegido con el mismo middleware de sesión que
  `GET /api/pedidos` (listar) — es una acción de gestión del admin, no algo que haga el cliente.
  `404` si el pedido no existe.
- El cliente en `/armar-pedido` no ve ni elige estado — sigue siendo un campo de gestión interna.
- Frontend: selector (`<select>`) por pedido en `/pedidos`, junto al resto de la info de cada
  tarjeta. Al cambiar, actualiza solo esa fila en el estado local de React (sin refetch completo
  de la lista); si falla, muestra el error debajo del selector de ese pedido puntual.

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
