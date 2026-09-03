import { createContext, useContext, useEffect, useState } from 'react';
import { getMe } from './api';

// isAuthenticated: null mientras se resuelve el chequeo inicial contra el
// backend, true/false una vez que se sabe.
export const AuthContext = createContext({
  isAuthenticated: null,
  setAuthenticated: () => {},
});

export function AuthProvider({ children }) {
  const [isAuthenticated, setAuthenticated] = useState(null);

  useEffect(() => {
    getMe()
      .then(() => setAuthenticated(true))
      .catch(() => setAuthenticated(false));
  }, []);

  return (
    <AuthContext.Provider value={{ isAuthenticated, setAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
