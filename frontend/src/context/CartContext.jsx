import { createContext, useContext, useEffect, useState } from "react";
import axios from "axios";
import { useAuth } from "./AuthContext";
import API_URL from "../config";

const CartContext = createContext(null);

export function CartProvider({ children }) {
  const { user } = useAuth();
  const [cart, setCart] = useState([]);
  const [localCart, setLocalCart] = useState([]);

  useEffect(() => {
    if (user) {
      axios.get(`${API_URL}/api/cart`, { withCredentials: true })
        .then(res => setCart(res.data || []))
        .catch(() => {});
    } else {
      setCart([]);
    }
  }, [user]);

  const activeCart = user ? cart : localCart;

  const addToCart = async (product, quantity = 1) => {
    if (user) {
      await axios.post(`${API_URL}/api/cart`, { product_id: product.id, quantity }, { withCredentials: true });
      const res = await axios.get(`${API_URL}/api/cart`, { withCredentials: true });
      setCart(res.data || []);
    } else {
      setLocalCart(prev => {
        const existing = prev.find(i => i.product_id === product.id);
        if (existing) return prev.map(i => i.product_id === product.id ? { ...i, quantity } : i);
        return [...prev, { product_id: product.id, quantity, name: product.name, price: product.price }];
      });
    }
  };

  const removeFromCart = async (product_id) => {
    if (user) {
      await axios.delete(`${API_URL}/api/cart/${product_id}`, { withCredentials: true });
      setCart(prev => prev.filter(i => i.product_id !== product_id));
    } else {
      setLocalCart(prev => prev.filter(i => i.product_id !== product_id));
    }
  };

  const clearCart = () => { setCart([]); setLocalCart([]); };

  const cartCount = activeCart.reduce((sum, i) => sum + (i.quantity || 1), 0);
  const cartTotal = activeCart.reduce((sum, i) => sum + i.price * (i.quantity || 1), 0);

  return (
    <CartContext.Provider value={{ cart: activeCart, addToCart, removeFromCart, clearCart, cartCount, cartTotal }}>
      {children}
    </CartContext.Provider>
  );
}

export const useCart = () => useContext(CartContext);
