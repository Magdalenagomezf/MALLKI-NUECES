// Wrapper delgado sobre fetch para hablar con el backend Go en /api.
// En dev, Vite proxea /api hacia http://localhost:8080 (ver vite.config.js).

const BASE_URL = '/api';

/**
 * Ejecuta un fetch contra el backend y devuelve el JSON parseado.
 * Si la respuesta no es ok, intenta leer {error: "..."} del body
 * y lanza un Error con ese mensaje (o uno genérico si no hay body).
 */
async function request(path, options = {}) {
  let response;
  try {
    response = await fetch(`${BASE_URL}${path}`, {
      headers: { 'Content-Type': 'application/json' },
      ...options,
    });
  } catch {
    throw new Error('No se pudo conectar con el servidor. Verificá tu conexión e intentá de nuevo.');
  }

  if (response.status === 204) {
    return null;
  }

  let body = null;
  try {
    body = await response.json();
  } catch {
    // body vacío o no-JSON; lo ignoramos si la respuesta fue ok
  }

  if (!response.ok) {
    const message = body?.error || `Error inesperado (${response.status})`;
    throw new Error(message);
  }

  return body;
}

// --- Productos ---

export function getProductos() {
  return request('/productos');
}

export function getProducto(id) {
  return request(`/productos/${id}`);
}

export function createProducto(producto) {
  return request('/productos', {
    method: 'POST',
    body: JSON.stringify(producto),
  });
}

export function updateProducto(id, producto) {
  return request(`/productos/${id}`, {
    method: 'PUT',
    body: JSON.stringify(producto),
  });
}

export function deleteProducto(id) {
  return request(`/productos/${id}`, { method: 'DELETE' });
}

// --- Pedidos ---

export function getPedidos() {
  return request('/pedidos');
}

export function createPedido(pedido) {
  return request('/pedidos', {
    method: 'POST',
    body: JSON.stringify(pedido),
  });
}
