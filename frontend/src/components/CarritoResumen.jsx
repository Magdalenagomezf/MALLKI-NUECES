/**
 * Lista de items en el carrito con su subtotal, permite editar la
 * cantidad_kg o eliminar un item, y muestra el total acumulado.
 */
export function CarritoResumen({ items, onEditarCantidad, onEliminar }) {
  const total = items.reduce((acc, item) => acc + item.cantidad_kg * item.precio_por_kg, 0);

  if (items.length === 0) {
    return <p className="text-secondary-500">Todavía no agregaste productos al pedido.</p>;
  }

  return (
    <div className="overflow-x-auto rounded-2xl bg-white p-4 shadow-sm ring-1 ring-secondary-100 sm:p-6">
      <table className="w-full border-collapse text-sm">
        <thead>
          <tr>
            <th className="border-b border-secondary-200 px-2 py-2 text-left">Producto</th>
            <th className="border-b border-secondary-200 px-2 py-2 text-left">Cantidad (kg)</th>
            <th className="border-b border-secondary-200 px-2 py-2 text-left">Precio/kg</th>
            <th className="border-b border-secondary-200 px-2 py-2 text-left">Subtotal</th>
            <th className="border-b border-secondary-200 px-2 py-2"></th>
          </tr>
        </thead>
        <tbody>
          {items.map((item) => (
            <tr key={item.producto_id}>
              <td className="border-b border-secondary-100 px-2 py-2">{item.nombre}</td>
              <td className="border-b border-secondary-100 px-2 py-2">
                <input
                  type="number"
                  min="0.01"
                  step="0.01"
                  value={item.cantidad_kg}
                  onChange={(e) => onEditarCantidad(item.producto_id, parseFloat(e.target.value) || 0)}
                  className="w-20 rounded-md border border-secondary-200 px-2 py-1 focus:border-primary-500 focus:outline-none"
                />
              </td>
              <td className="border-b border-secondary-100 px-2 py-2">${item.precio_por_kg.toFixed(2)}</td>
              <td className="border-b border-secondary-100 px-2 py-2 font-semibold">
                ${(item.cantidad_kg * item.precio_por_kg).toFixed(2)}
              </td>
              <td className="border-b border-secondary-100 px-2 py-2">
                <button
                  type="button"
                  onClick={() => onEliminar(item.producto_id)}
                  className="rounded-md border border-red-600 px-2 py-1 text-xs font-semibold text-red-600 transition-colors hover:bg-red-50"
                >
                  Quitar
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <p className="mt-2 text-right text-lg font-bold text-secondary-800">Total: ${total.toFixed(2)}</p>
    </div>
  );
}
