/**
 * Divisor orgánico tipo "trazo de pincel" para romper la línea recta entre
 * secciones (hero de foto -> contenido). El color se toma de `currentColor`,
 * así que se controla con una clase de texto de Tailwind (p. ej. text-primary-50)
 * en vez de una prop de color aparte.
 */
export function SectionDivider({ className = '', flip = false }) {
  return (
    <div
      aria-hidden="true"
      className={`pointer-events-none absolute inset-x-0 bottom-0 h-12 w-full overflow-hidden sm:h-20 ${
        flip ? '-scale-y-100' : ''
      } ${className}`}
    >
      <svg viewBox="0 0 1200 120" preserveAspectRatio="none" className="h-full w-full">
        <path
          d="M0,32 C120,80 260,0 420,38 C580,76 700,10 860,42 C1000,70 1100,20 1200,48 L1200,120 L0,120 Z"
          fill="currentColor"
        />
      </svg>
    </div>
  );
}
