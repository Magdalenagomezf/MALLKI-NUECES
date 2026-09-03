import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../api';

export function Login() {
  const navigate = useNavigate();
  const [usuario, setUsuario] = useState('');
  const [contrasena, setContrasena] = useState('');
  const [enviando, setEnviando] = useState(false);
  const [error, setError] = useState(null);

  async function handleSubmit(e) {
    e.preventDefault();
    setError(null);
    setEnviando(true);
    try {
      await login(usuario, contrasena);
      navigate('/pedidos');
    } catch (err) {
      setError(err.message);
    } finally {
      setEnviando(false);
    }
  }

  return (
    <div className="mx-auto flex w-full max-w-md flex-col gap-8 px-4 py-16 sm:px-6">
      <h1 className="font-heading text-3xl text-secondary-800">Ingreso admin</h1>

      <form
        onSubmit={handleSubmit}
        className="flex flex-col gap-4 rounded-2xl bg-white p-5 shadow-sm ring-1 ring-secondary-100 sm:p-6"
      >
        <label className="flex flex-col gap-1 text-sm font-semibold text-secondary-700">
          Usuario
          <input
            type="text"
            value={usuario}
            onChange={(e) => setUsuario(e.target.value)}
            autoComplete="username"
            className="rounded-md border border-secondary-200 px-3 py-2 text-base font-normal text-secondary-800 focus:border-primary-500 focus:outline-none"
          />
        </label>
        <label className="flex flex-col gap-1 text-sm font-semibold text-secondary-700">
          Contraseña
          <input
            type="password"
            value={contrasena}
            onChange={(e) => setContrasena(e.target.value)}
            autoComplete="current-password"
            className="rounded-md border border-secondary-200 px-3 py-2 text-base font-normal text-secondary-800 focus:border-primary-500 focus:outline-none"
          />
        </label>

        {error && <p className="rounded-md bg-red-50 px-4 py-3 text-red-700">{error}</p>}

        <button
          type="submit"
          disabled={enviando}
          className="self-start rounded-md bg-primary-500 px-6 py-2.5 text-base font-bold text-secondary-900 transition-colors hover:bg-primary-600 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {enviando ? 'Ingresando...' : 'Ingresar'}
        </button>
      </form>
    </div>
  );
}
