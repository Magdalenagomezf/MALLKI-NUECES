export function Footer() {
  return (
    <footer className="mt-20 border-t-4 border-primary-500 bg-secondary-900 py-14 text-secondary-100">
      <div className="mx-auto flex w-full max-w-5xl flex-col items-center gap-3 px-4 text-center sm:px-6">
        <p className="font-heading text-3xl text-primary-400 sm:text-4xl">Mallki Nueces</p>
        <p className="max-w-md text-sm text-secondary-200">
          Frutos secos de Catamarca, directo del productor. Catamarca ♥ Argentina.
        </p>
        {/* Placeholder: reemplazar por el contacto real del negocio. */}
        <p className="text-sm text-secondary-400">Consultas: contacto@mallkinueces.com.ar</p>
      </div>
    </footer>
  );
}
