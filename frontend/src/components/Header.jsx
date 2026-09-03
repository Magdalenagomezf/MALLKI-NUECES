import { useState } from 'react';
import { NavLink } from 'react-router-dom';
import logoMallki from '../assets/logo-mallki.jpeg';

const LINKS = [
  { to: '/', label: 'Catálogo' },
  { to: '/armar-pedido', label: 'Armar pedido' },
  { to: '/pedidos', label: 'Pedidos' },
  { to: '/login', label: 'Admin' },
];

function navLinkClasses({ isActive }) {
  return [
    'rounded-full px-4 py-2 text-base font-semibold transition-colors',
    isActive
      ? 'bg-primary-500 text-secondary-900'
      : 'text-secondary-100 hover:bg-secondary-700 hover:text-white',
  ].join(' ');
}

export function Header() {
  const [menuAbierto, setMenuAbierto] = useState(false);

  return (
    <header className="sticky top-0 z-20 border-b border-secondary-700 bg-secondary-800 shadow-md">
      <div className="mx-auto flex max-w-5xl items-center justify-between gap-4 px-4 py-3 sm:px-6">
        <div className="flex items-center gap-3">
          <img
            src={logoMallki}
            alt="Mallki Nueces"
            className="h-12 w-12 rounded-full border-2 border-primary-500 object-cover sm:h-14 sm:w-14"
          />
          <span className="font-heading text-3xl leading-none text-primary-500 sm:text-4xl">Mallki Nueces</span>
        </div>

        <nav className="hidden items-center gap-2 sm:flex">
          {LINKS.map((link) => (
            <NavLink key={link.to} to={link.to} end={link.to === '/'} className={navLinkClasses}>
              {link.label}
            </NavLink>
          ))}
        </nav>

        <button
          type="button"
          className="flex h-10 w-10 items-center justify-center rounded-md border border-secondary-600 text-primary-500 sm:hidden"
          aria-label="Abrir menú de navegación"
          aria-expanded={menuAbierto}
          onClick={() => setMenuAbierto((abierto) => !abierto)}
        >
          <span className="text-2xl leading-none">{menuAbierto ? '✕' : '☰'}</span>
        </button>
      </div>

      {menuAbierto && (
        <nav className="flex flex-col gap-1 border-t border-secondary-700 px-4 pb-3 pt-2 sm:hidden">
          {LINKS.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              end={link.to === '/'}
              className={navLinkClasses}
              onClick={() => setMenuAbierto(false)}
            >
              {link.label}
            </NavLink>
          ))}
        </nav>
      )}
    </header>
  );
}
