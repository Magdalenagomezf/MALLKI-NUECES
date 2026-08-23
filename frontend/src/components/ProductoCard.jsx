import { useState } from 'react';

/**
 * Muestra un producto del catálogo. Si se pasa onAgregar, además
 * renderiza un input de cantidad_kg y un botón para agregarlo al carrito
 * (usado desde ArmarPedido). Sin onAgregar es solo de lectura (Catalogo).
 */
export function ProductoCard({ producto, onAgregar }) {
  const [cantidad, setCantidad] = useState('');
  const sinStock = producto.stock_kg <= 0;

  function handleAgregar() {
    const cantidadKg = parseFloat(cantidad);
    if (!cantidadKg || cantidadKg <= 0) return;
    onAgregar(producto, cantidadKg);
    setCantidad('');
  }

  return (
    <div
      className={`group flex flex-col gap-3 rounded-2xl bg-white p-5 shadow-sm ring-1 ring-secondary-100 transition-all hover:-translate-y-1 hover:shadow-lg${
        sinStock ? ' opacity-60' : ''
      }`}
    >
      <div className="flex flex-col gap-2">
        <div className="flex items-start justify-between gap-2">
          <h3 className="text-lg font-bold text-secondary-800">{producto.nombre}</h3>
          {producto.categoria && (
            <span className="shrink-0 rounded-full bg-primary-100 px-2.5 py-1 text-xs font-bold uppercase tracking-wide text-primary-700">
              {producto.categoria}
            </span>
          )}
        </div>
        {producto.descripcion && <p className="text-sm text-secondary-500">{producto.descripcion}</p>}
        <div className="mt-1 flex items-center justify-between">
          <span className="text-xl font-extrabold text-secondary-800">
            ${producto.precio_por_kg.toFixed(2)}{' '}
            <span className="text-sm font-normal text-secondary-400">/ kg</span>
          </span>
          <span className={sinStock ? 'text-xs font-semibold text-red-600' : 'text-xs font-normal text-secondary-400'}>
            {sinStock ? 'Sin stock' : `${producto.stock_kg} kg disp.`}
          </span>
        </div>
      </div>

      {onAgregar && (
        <div className="mt-1 flex gap-2">
          <input
            type="number"
            min="0.01"
            step="0.01"
            placeholder="kg"
            value={cantidad}
            disabled={sinStock}
            onChange={(e) => setCantidad(e.target.value)}
            className="w-20 rounded-lg border border-secondary-200 px-2 py-1.5 text-sm focus:border-primary-500 focus:outline-none disabled:bg-secondary-50"
          />
          <button
            type="button"
            onClick={handleAgregar}
            disabled={sinStock || !cantidad}
            className="flex-1 rounded-lg bg-primary-500 px-3 py-1.5 text-sm font-bold text-secondary-900 transition-colors hover:bg-primary-600 disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:bg-primary-500"
          >
            Agregar
          </button>
        </div>
      )}
    </div>
  );
}
