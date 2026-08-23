// react-alice-carousel es CJS sin "exports" map: según cómo lo interprete
// el bundler, el default export puede llegar doble-envuelto (el objeto del
// módulo en vez del componente). Se desenvuelve a mano para no depender de
// que Vite/esbuild lo resuelvan "bien" en dev y en build por igual.
import AliceCarouselImport from 'react-alice-carousel';
import 'react-alice-carousel/lib/alice-carousel.css';

const AliceCarousel = AliceCarouselImport.default ?? AliceCarouselImport;

const RESPONSIVE = {
  0: { items: 1 },
  640: { items: 2 },
  1024: { items: 3 },
};

/**
 * Carrusel de productos destacados para el hero de Catálogo.
 * Es una mejora visual, no crítica: si no hay productos, no se renderiza
 * nada (Catalogo ya maneja el estado vacío general).
 */
export function FeaturedCarousel({ productos }) {
  if (!productos || productos.length === 0) return null;

  const items = productos.map((producto) => (
    <div key={producto.id} className="h-full px-2">
      <div className="flex h-full flex-col justify-between rounded-2xl bg-white p-5 shadow-lg">
        <div>
          {producto.categoria && (
            <span className="inline-block rounded-full bg-primary-100 px-2.5 py-1 text-xs font-bold uppercase tracking-wide text-primary-700">
              {producto.categoria}
            </span>
          )}
          <h3 className="mt-2 text-xl font-bold text-secondary-800">{producto.nombre}</h3>
          {producto.descripcion && (
            <p className="mt-1 line-clamp-2 text-sm text-secondary-500">{producto.descripcion}</p>
          )}
        </div>
        <p className="mt-3 text-lg font-extrabold text-secondary-800">
          ${producto.precio_por_kg.toFixed(2)} <span className="text-sm font-normal text-secondary-400">/ kg</span>
        </p>
      </div>
    </div>
  ));

  return (
    <div className="featured-carousel">
      <AliceCarousel
        items={items}
        responsive={RESPONSIVE}
        autoPlay
        autoPlayInterval={3500}
        infinite
        disableButtonsControls
        disableDotsControls={items.length <= 1}
        mouseTracking
      />
    </div>
  );
}
