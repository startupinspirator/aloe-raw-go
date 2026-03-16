import { Routes, Route, useSearchParams } from "react-router-dom";
import { useEffect } from "react";
import { AuthProvider } from "./context/AuthContext";
import { CartProvider } from "./context/CartContext";
import Navbar from "./components/Navbar";
import Home from "./pages/Home";
import Product from "./pages/Product";
import Cart from "./pages/Cart";
import Checkout from "./pages/Checkout";
import Orders from "./pages/Orders";

import AdminDashboard from "./pages/AdminDashboard";
import Profile from "./pages/Profile";

function LoginHandler() {
  const [params] = useSearchParams();
  useEffect(() => {
    if (params.get("login") === "success") {
      window.history.replaceState({}, "", "/");
    }
  }, [params]);
  return null;
}

export default function App() {
  return (
    <AuthProvider>
      <CartProvider>
        <LoginHandler />
        <Navbar />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/product/:id" element={<Product />} />
          <Route path="/cart" element={<Cart />} />
          <Route path="/checkout" element={<Checkout />} />
          <Route path="/orders" element={<Orders />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/admin" element={<AdminDashboard />} />
        </Routes>
      </CartProvider>
    </AuthProvider>
  );
}
