import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Header } from './components/Header';
import { Footer } from './components/Footer';
import { Catalogo } from './components/Catalogo';
import { ArmarPedido } from './components/ArmarPedido';
import { Pedidos } from './components/Pedidos';

function App() {
  return (
    <BrowserRouter>
      <div className="flex min-h-screen flex-col bg-primary-50">
        <Header />
        {/* Sin max-width/padding acá a propósito: cada página decide su propio
            ancho, así el Catálogo puede tener un hero full-bleed. */}
        <main className="flex-1">
          <Routes>
            <Route path="/" element={<Catalogo />} />
            <Route path="/armar-pedido" element={<ArmarPedido />} />
            <Route path="/pedidos" element={<Pedidos />} />
          </Routes>
        </main>
        <Footer />
      </div>
    </BrowserRouter>
  );
}

export default App;
