import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { getPedidos, logout, updatePedidoEstado } from '../api';
import { useAuth } from '../auth';
import { PageBanner } from './PageBanner';
import bannerFoto from '../assets/fotos/mitades-luz.jpg';

function formatearFecha(fechaISO) {
  try {
    return new Date(fechaISO).toLocaleString('es-AR', {
      dateStyle: 'short',
      timeStyle: 'short',
    });
  } catch {
    return fechaISO;
  }
}

const ESTADOS = [
  { value: 'pendiente', label: 'Pendiente' },
  { value: 'confirmado', label: 'Confirmado' },
  { value: 'en_preparacion', label: 'En preparación' },
  { value: 'entregado', label: 'Entregado' },
];
const ESTADO_LABELS = Object.fromEntries(ESTADOS.map((e) => [e.value, e.label]));

const ESTADO_BADGE_CLASSES = {
  pendiente: 'bg-yellow-100 text-yellow-800',
  confirmado: 'bg-blue-100 text-blue-800',
  en_preparacion: 'bg-orange-100 text-orange-800',
  entregado: 'bg-green-100 text-green-800',
};
const ESTADO_BADGE_DEFAULT = 'bg-secondary-100 text-secondary-700';

// Número de WhatsApp del admin (formato internacional, sin "+" ni espacios,
// con el "9" que WhatsApp requiere para números móviles de Argentina).
const WHATSAPP_ADMIN = '5493834003867';

function armarMensajeWhatsApp(pedido) {
  const items = pedido.items
    .map((item) => `- ${item.producto_nombre}: ${item.cantidad_kg} kg`)
    .join('\n');
  return (
    `Nuevo pedido #${pedido.id}\n` +
    `Cliente: ${pedido.cliente_nombre} (${pedido.cliente_contacto})\n` +
    `Pago: ${pedido.metodo_pago}\n\n` +
    `${items}`
  );
}

export function Pedidos() {
  const navigate = useNavigate();
  const { setAuthenticated } = useAuth();
  const [pedidos, setPedidos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [estadoErrors, setEstadoErrors] = useState({});

  useEffect(() => {
    async function cargar() {
      setLoading(true);
      setError(null);
      try {
        const data = await getPedidos();
        setPedidos(data ?? []);
      } catch (err) {
        if (err.status === 401) {
          setAuthenticated(false);
          navigate('/login');
          return;
        }
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }
    cargar();
  }, [navigate, setAuthenticated]);

  async function handleLogout() {
    try {
      await logout();
    } finally {
      setAuthenticated(false);
      navigate('/login');
    }
  }

  async function handleEstadoChange(pedidoId, nuevoEstado) {
    setEstadoErrors((prev) => ({ ...prev, [pedidoId]: null }));
    try {
      const actualizado = await updatePedidoEstado(pedidoId, nuevoEstado);
      setPedidos((prev) =>
        prev.map((p) => (p.id === pedidoId ? { ...p, estado: actualizado.estado } : p))
      );
    } catch (err) {
      setEstadoErrors((prev) => ({ ...prev, [pedidoId]: err.message }));
    }
  }

  return (
    <div className="flex flex-col">
      <PageBanner titulo="Pedidos" foto={bannerFoto} alt="Mitades de nuez en primer plano" />

      <div className="mx-auto w-full max-w-5xl px-4 py-10 sm:px-6">
        <div className="mb-6 flex justify-end">
          <button
            type="button"
            onClick={handleLogout}
            className="rounded-md border border-secondary-200 px-4 py-2 text-sm font-semibold text-secondary-700 transition-colors hover:bg-secondary-100"
          >
            Cerrar sesión
          </button>
        </div>

        {loading && <p className="text-secondary-500">Cargando pedidos...</p>}
        {error && (
          <p className="rounded-md bg-red-50 px-4 py-3 text-red-700">Error al cargar pedidos: {error}</p>
        )}
        {!loading && !error && pedidos.length === 0 && (
          <p className="text-secondary-500">Todavía no hay pedidos registrados.</p>
        )}

        {!loading && !error && pedidos.length > 0 && (
          <div className="flex flex-col gap-4">
            {pedidos.map((pedido) => (
              <article
                key={pedido.id}
                className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-secondary-100 transition-shadow hover:shadow-md"
              >
                <header className="flex items-center justify-between">
                  <h3 className="text-lg font-bold text-secondary-800">Pedido #{pedido.id}</h3>
                  <span
                    className={`rounded-full px-3 py-1 text-xs font-bold uppercase tracking-wide ${
                      ESTADO_BADGE_CLASSES[pedido.estado] ?? ESTADO_BADGE_DEFAULT
                    }`}
                  >
                    {ESTADO_LABELS[pedido.estado] ?? pedido.estado}
                  </span>
                </header>
                <p className="mt-1 text-secondary-700">
                  <strong>{pedido.cliente_nombre}</strong> — {pedido.cliente_contacto}
                </p>
                <p className="mt-0.5 text-sm text-secondary-500">
                  {formatearFecha(pedido.fecha_creacion)} · Pago: {pedido.metodo_pago}
                </p>
                <div className="mt-3 flex flex-wrap items-center gap-2">
                  <a
                    href={`https://wa.me/${WHATSAPP_ADMIN}?text=${encodeURIComponent(armarMensajeWhatsApp(pedido))}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="rounded-md border border-green-600 px-3 py-1 text-sm font-semibold text-green-700 transition-colors hover:bg-green-50"
                  >
                    Avisar por WhatsApp
                  </a>
                  <label htmlFor={`estado-${pedido.id}`} className="text-sm font-semibold text-secondary-700">
                    Estado:
                  </label>
                  <select
                    id={`estado-${pedido.id}`}
                    value={pedido.estado}
                    onChange={(e) => handleEstadoChange(pedido.id, e.target.value)}
                    className="rounded-md border border-secondary-200 px-2 py-1 text-sm text-secondary-700 focus:border-secondary-400 focus:outline-none"
                  >
                    {ESTADOS.map((estadoOpt) => (
                      <option key={estadoOpt.value} value={estadoOpt.value}>
                        {estadoOpt.label}
                      </option>
                    ))}
                  </select>
                  {estadoErrors[pedido.id] && (
                    <span className="text-sm text-red-700">{estadoErrors[pedido.id]}</span>
                  )}
                </div>
                <div className="mt-3 overflow-x-auto">
                  <table className="w-full border-collapse text-sm">
                    <thead>
                      <tr>
                        <th className="border-b border-secondary-200 px-2 py-2 text-left">Producto</th>
                        <th className="border-b border-secondary-200 px-2 py-2 text-left">Cantidad (kg)</th>
                        <th className="border-b border-secondary-200 px-2 py-2 text-left">Precio/kg</th>
                      </tr>
                    </thead>
                    <tbody>
                      {pedido.items.map((item) => (
                        <tr key={item.id}>
                          <td className="border-b border-secondary-100 px-2 py-2">{item.producto_nombre}</td>
                          <td className="border-b border-secondary-100 px-2 py-2">{item.cantidad_kg}</td>
                          <td className="border-b border-secondary-100 px-2 py-2">
                            ${item.precio_unitario.toFixed(2)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </article>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
