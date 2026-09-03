import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { getProductos, updateProductoStock } from '../api';
import { useAuth } from '../auth';
import { PageBanner } from './PageBanner';
import bannerFoto from '../assets/fotos/saco-arpillera.jpg';

export function Stock() {
  const navigate = useNavigate();
  const { setAuthenticated } = useAuth();
  const [productos, setProductos] = useState([]);
  const [drafts, setDrafts] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [stockErrors, setStockErrors] = useState({});

  useEffect(() => {
    async function cargar() {
      setLoading(true);
      setError(null);
      try {
        const data = await getProductos();
        const lista = data ?? [];
        setProductos(lista);
        setDrafts(Object.fromEntries(lista.map((p) => [p.id, String(p.stock_kg)])));
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

  function handleDraftChange(productoId, valor) {
    setDrafts((prev) => ({ ...prev, [productoId]: valor }));
  }

  async function guardarStock(productoId) {
    const valor = parseFloat(drafts[productoId]);
    if (Number.isNaN(valor) || valor < 0) {
      setStockErrors((prev) => ({ ...prev, [productoId]: 'Ingresá un valor válido (0 o mayor)' }));
      return;
    }

    setStockErrors((prev) => ({ ...prev, [productoId]: null }));
    try {
      const actualizado = await updateProductoStock(productoId, valor);
      setProductos((prev) =>
        prev.map((p) => (p.id === productoId ? { ...p, stock_kg: actualizado.stock_kg } : p))
      );
      setDrafts((prev) => ({ ...prev, [productoId]: String(actualizado.stock_kg) }));
    } catch (err) {
      setStockErrors((prev) => ({ ...prev, [productoId]: err.message }));
    }
  }

  function handleKeyDown(event) {
    if (event.key === 'Enter') {
      event.currentTarget.blur();
    }
  }

  return (
    <div className="flex flex-col">
      <PageBanner titulo="Stock" foto={bannerFoto} alt="Nueces guardadas en un saco de arpillera" />

      <div className="mx-auto w-full max-w-4xl px-4 py-10 sm:px-6">
        {loading && <p className="text-secondary-500">Cargando productos...</p>}
        {error && (
          <p className="rounded-md bg-red-50 px-4 py-3 text-red-700">Error al cargar productos: {error}</p>
        )}
        {!loading && !error && productos.length === 0 && (
          <p className="text-secondary-500">Todavía no hay productos cargados.</p>
        )}

        {!loading && !error && productos.length > 0 && (
          <div className="overflow-x-auto rounded-2xl bg-white shadow-sm ring-1 ring-secondary-100">
            <table className="w-full border-collapse text-sm">
              <thead>
                <tr>
                  <th className="border-b border-secondary-200 px-4 py-3 text-left">Producto</th>
                  <th className="border-b border-secondary-200 px-4 py-3 text-left">Categoría</th>
                  <th className="border-b border-secondary-200 px-4 py-3 text-left">Precio/kg</th>
                  <th className="border-b border-secondary-200 px-4 py-3 text-left">Stock (kg)</th>
                </tr>
              </thead>
              <tbody>
                {productos.map((producto) => (
                  <tr key={producto.id}>
                    <td className="border-b border-secondary-100 px-4 py-3 font-semibold text-secondary-800">
                      {producto.nombre}
                    </td>
                    <td className="border-b border-secondary-100 px-4 py-3 text-secondary-600">
                      {producto.categoria}
                    </td>
                    <td className="border-b border-secondary-100 px-4 py-3 text-secondary-600">
                      ${producto.precio_por_kg.toFixed(2)}
                    </td>
                    <td className="border-b border-secondary-100 px-4 py-3">
                      <div className="flex flex-col gap-1">
                        <input
                          type="number"
                          min="0"
                          step="0.1"
                          value={drafts[producto.id] ?? ''}
                          onChange={(e) => handleDraftChange(producto.id, e.target.value)}
                          onBlur={() => guardarStock(producto.id)}
                          onKeyDown={handleKeyDown}
                          className="w-28 rounded-md border border-secondary-200 px-2 py-1 text-sm text-secondary-700 focus:border-secondary-400 focus:outline-none"
                        />
                        {stockErrors[producto.id] && (
                          <span className="text-xs text-red-700">{stockErrors[producto.id]}</span>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
