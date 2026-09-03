import { useState } from 'react';
import { useProductos } from '../hooks/useProductos';
import { createPedido } from '../api';
import { ProductoCard } from './ProductoCard';
import { CarritoResumen } from './CarritoResumen';
import { PageBanner } from './PageBanner';
import bannerFoto from '../assets/fotos/nueces-tela.jpg';

export function ArmarPedido() {
  const { productos, loading, error: errorCatalogo } = useProductos();

  // Carrito en memoria: array de {producto_id, nombre, precio_por_kg, cantidad_kg}
  const [carrito, setCarrito] = useState([]);
  const [clienteNombre, setClienteNombre] = useState('');
  const [clienteContacto, setClienteContacto] = useState('');
  const [metodoPago, setMetodoPago] = useState('');
  const [enviando, setEnviando] = useState(false);
  const [errorEnvio, setErrorEnvio] = useState(null);
  const [confirmacion, setConfirmacion] = useState(null);
  const [erroresValidacion, setErroresValidacion] = useState([]);

  function agregarAlCarrito(producto, cantidadKg) {
    setCarrito((prev) => {
      const existente = prev.find((item) => item.producto_id === producto.id);
      if (existente) {
        return prev.map((item) =>
          item.producto_id === producto.id
            ? { ...item, cantidad_kg: item.cantidad_kg + cantidadKg }
            : item,
        );
      }
      return [
        ...prev,
        {
          producto_id: producto.id,
          nombre: producto.nombre,
          precio_por_kg: producto.precio_por_kg,
          cantidad_kg: cantidadKg,
        },
      ];
    });
    setConfirmacion(null);
  }

  function editarCantidad(productoId, nuevaCantidad) {
    setCarrito((prev) =>
      prev.map((item) => (item.producto_id === productoId ? { ...item, cantidad_kg: nuevaCantidad } : item)),
    );
  }

  function eliminarDelCarrito(productoId) {
    setCarrito((prev) => prev.filter((item) => item.producto_id !== productoId));
  }

  function validar() {
    const errores = [];
    if (!clienteNombre.trim()) errores.push('El nombre del cliente es requerido.');
    if (!clienteContacto.trim()) errores.push('El contacto del cliente es requerido.');
    if (!metodoPago) errores.push('Elegí un método de pago.');
    if (carrito.length === 0) errores.push('Agregá al menos un producto al pedido.');
    if (carrito.some((item) => !item.cantidad_kg || item.cantidad_kg <= 0)) {
      errores.push('Todas las cantidades deben ser mayores a 0.');
    }
    return errores;
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setErrorEnvio(null);
    setConfirmacion(null);

    const errores = validar();
    setErroresValidacion(errores);
    if (errores.length > 0) return;

    setEnviando(true);
    try {
      const pedido = await createPedido({
        cliente_nombre: clienteNombre.trim(),
        cliente_contacto: clienteContacto.trim(),
        metodo_pago: metodoPago,
        items: carrito.map((item) => ({
          producto_id: item.producto_id,
          cantidad_kg: item.cantidad_kg,
        })),
      });
      setConfirmacion(pedido);
      setCarrito([]);
      setClienteNombre('');
      setClienteContacto('');
      setMetodoPago('');
      setErroresValidacion([]);
    } catch (err) {
      setErrorEnvio(err.message);
    } finally {
      setEnviando(false);
    }
  }

  return (
    <div className="flex flex-col">
      <PageBanner titulo="Armar pedido" foto={bannerFoto} alt="Nueces sobre una tela clara" />

      <div className="mx-auto flex w-full max-w-5xl flex-col gap-8 px-4 py-10 sm:px-6">
        <section>
          {loading && <p className="text-secondary-500">Cargando catálogo...</p>}
          {errorCatalogo && (
            <p className="rounded-md bg-red-50 px-4 py-3 text-red-700">
              Error al cargar productos: {errorCatalogo}
            </p>
          )}

          {!loading && !errorCatalogo && (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {productos.map((producto) => (
                <ProductoCard key={producto.id} producto={producto} onAgregar={agregarAlCarrito} />
              ))}
            </div>
          )}
        </section>

        <section>
          <h2 className="mb-3 font-heading text-3xl text-secondary-800">Carrito</h2>
          <CarritoResumen items={carrito} onEditarCantidad={editarCantidad} onEliminar={eliminarDelCarrito} />
        </section>

        <section>
          <h2 className="mb-3 font-heading text-3xl text-secondary-800">Datos del pedido</h2>
          <form
            onSubmit={handleSubmit}
            className="flex max-w-md flex-col gap-4 rounded-2xl bg-white p-5 shadow-sm ring-1 ring-secondary-100 sm:p-6"
          >
            <label className="flex flex-col gap-1 text-sm font-semibold text-secondary-700">
              Nombre del cliente
              <input
                type="text"
                value={clienteNombre}
                onChange={(e) => setClienteNombre(e.target.value)}
                placeholder="Ej: Juan Pérez"
                className="rounded-md border border-secondary-200 px-3 py-2 text-base font-normal text-secondary-800 focus:border-primary-500 focus:outline-none"
              />
            </label>
            <label className="flex flex-col gap-1 text-sm font-semibold text-secondary-700">
              Contacto (teléfono o email)
              <input
                type="text"
                value={clienteContacto}
                onChange={(e) => setClienteContacto(e.target.value)}
                placeholder="Ej: 3511234567"
                className="rounded-md border border-secondary-200 px-3 py-2 text-base font-normal text-secondary-800 focus:border-primary-500 focus:outline-none"
              />
            </label>
            <fieldset className="flex flex-col gap-1 text-sm font-semibold text-secondary-700">
              <legend className="mb-1">Método de pago</legend>
              <div className="flex gap-4 font-normal text-secondary-800">
                <label className="flex items-center gap-2">
                  <input
                    type="radio"
                    name="metodoPago"
                    value="transferencia"
                    checked={metodoPago === 'transferencia'}
                    onChange={(e) => setMetodoPago(e.target.value)}
                  />
                  Transferencia
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="radio"
                    name="metodoPago"
                    value="efectivo"
                    checked={metodoPago === 'efectivo'}
                    onChange={(e) => setMetodoPago(e.target.value)}
                  />
                  Efectivo
                </label>
              </div>
            </fieldset>

            {erroresValidacion.length > 0 && (
              <ul className="list-disc rounded-md bg-red-50 px-4 py-3 pl-8 text-red-700">
                {erroresValidacion.map((msg) => (
                  <li key={msg}>{msg}</li>
                ))}
              </ul>
            )}

            {errorEnvio && <p className="rounded-md bg-red-50 px-4 py-3 text-red-700">{errorEnvio}</p>}

            {confirmacion && (
              <p className="rounded-md bg-green-50 px-4 py-3 text-green-700">
                ¡Pedido #{confirmacion.id} confirmado! Total de items: {confirmacion.items.length}.
              </p>
            )}

            <button
              type="submit"
              disabled={enviando}
              className="self-start rounded-md bg-primary-500 px-6 py-2.5 text-base font-bold text-secondary-900 transition-colors hover:bg-primary-600 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {enviando ? 'Enviando...' : 'Confirmar pedido'}
            </button>
          </form>
        </section>
      </div>
    </div>
  );
}
