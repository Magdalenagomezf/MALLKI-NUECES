import { Link } from 'react-router-dom';
import { useProductos } from '../hooks/useProductos';
import { ProductoCard } from './ProductoCard';
import { FeaturedCarousel } from './FeaturedCarousel';
import { SectionDivider } from './SectionDivider';
import heroFoto from '../assets/fotos/saco-arpillera.jpg';

const MAX_DESTACADOS = 6;

export function Catalogo() {
  const { productos, loading, error } = useProductos();

  return (
    <div className="flex flex-col">
      {/* Hero full-bleed. Copy de ejemplo — el dueño del negocio puede reescribirla libremente. */}
      <section className="relative flex min-h-[75vh] items-center justify-center overflow-hidden bg-secondary-900 text-white sm:min-h-[85vh]">
        <img
          src={heroFoto}
          alt="Nueces recién cosechadas asomando de un saco de arpillera"
          className="absolute inset-0 h-full w-full object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-secondary-900 via-secondary-900/55 to-secondary-900/10" />
        <div className="relative z-10 flex max-w-3xl flex-col items-center gap-5 px-6 text-center">
          <p className="font-heading text-4xl text-primary-400 sm:text-5xl">Frutos secos</p>
          <h1 className="font-sans text-4xl font-extrabold uppercase leading-[1.05] tracking-tight sm:text-6xl">
            Directo de Catamarca a tu mesa
          </h1>
          <p className="max-w-xl text-base text-secondary-100 sm:text-lg">
            Nueces mariposas, cuartos o con cascara. Pedís por kilo, coordinamos la entrega.
          </p>
          <Link
            to="/armar-pedido"
            className="mt-2 rounded-full bg-primary-500 px-8 py-3 text-base font-bold text-secondary-900 shadow-lg transition-transform hover:-translate-y-0.5 hover:bg-primary-400"
          >
            Armar pedido
          </Link>
        </div>
        <SectionDivider className="text-primary-50" />
      </section>

      {!loading && !error && productos.length > 0 && (
        <section className="bg-secondary-800 py-14 text-white">
          <div className="mx-auto w-full max-w-5xl px-4 sm:px-6">
            <p className="font-heading text-3xl text-primary-400 sm:text-4xl">Destacados</p>
            <p className="mb-6 mt-1 text-sm text-secondary-200">Lo que más piden esta semana.</p>
            <FeaturedCarousel productos={productos.slice(0, MAX_DESTACADOS)} />
          </div>
        </section>
      )}

      <section className="mx-auto w-full max-w-5xl px-4 py-14 sm:px-6">
        <h2 className="font-heading text-3xl text-secondary-800 sm:text-4xl">Catálogo completo</h2>
        <p className="mb-6 mt-1 text-sm text-secondary-500">Todo lo que tenemos disponible hoy.</p>

        {loading && <p className="text-secondary-500">Cargando catálogo...</p>}
        {error && (
          <p className="rounded-md bg-red-50 px-4 py-3 text-red-700">Error al cargar productos: {error}</p>
        )}
        {!loading && !error && productos.length === 0 && (
          <p className="text-secondary-500">No hay productos cargados todavía.</p>
        )}

        {!loading && !error && productos.length > 0 && (
          <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {productos.map((producto) => (
              <ProductoCard key={producto.id} producto={producto} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
