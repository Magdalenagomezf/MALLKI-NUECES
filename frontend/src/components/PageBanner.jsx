import { SectionDivider } from './SectionDivider';

/**
 * Banner full-bleed compacto para páginas internas (Armar pedido, Pedidos).
 * No es el hero principal del Catálogo, pero comparte el mismo lenguaje
 * visual: foto real de fondo, overlay oscuro, título bold uppercase y
 * divisor orgánico — así todas las páginas se sienten parte del mismo sitio.
 */
export function PageBanner({ titulo, foto, alt }) {
  return (
    <section className="relative flex h-56 items-center justify-center overflow-hidden bg-secondary-900 text-white sm:h-64">
      <img src={foto} alt={alt} className="absolute inset-0 h-full w-full object-cover" />
      <div className="absolute inset-0 bg-secondary-900/60" />
      <h1 className="relative z-10 px-6 font-sans text-3xl font-extrabold uppercase tracking-tight sm:text-5xl">
        {titulo}
      </h1>
      <SectionDivider className="text-primary-50" />
    </section>
  );
}
