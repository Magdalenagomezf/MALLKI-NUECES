import { useEffect, useState, useCallback } from 'react';
import { getProductos } from '../api';

/**
 * Hook compartido para cargar el catálogo de productos.
 * Lo usan tanto Catalogo como ArmarPedido, que necesitan la misma
 * lista con el mismo manejo de loading/error.
 */
export function useProductos() {
  const [productos, setProductos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const cargar = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getProductos();
      setProductos(data ?? []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    cargar();
  }, [cargar]);

  return { productos, loading, error, recargar: cargar };
}
