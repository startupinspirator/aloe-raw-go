import { useNavigate } from "react-router-dom";
import { useCart } from "../context/CartContext";
import { useAuth } from "../context/AuthContext";
import { IMG_SINGLE_BOX } from "../assets/images";

export default function Cart() {
  const { cart, removeFromCart, addToCart, cartTotal } = useCart();
  const { user, loginWithGoogle } = useAuth();
  const navigate = useNavigate();

  if (cart.length === 0) return (
    <main className="cart-page empty">
      <div className="empty-cart">
        <span className="empty-icon">🛒</span>
        <h2>Your cart is empty</h2>
        <p>Add some fresh aloe vera cubes to get started.</p>
        <button className="btn-primary" onClick={() => navigate("/product/1")}>Shop Now</button>
      </div>
    </main>
  );

  return (
    <main className="cart-page">
      <div className="cart-container">
        <h1 className="page-title">Your Cart</h1>
        <div className="cart-layout">
          <div className="cart-items">
            {cart.map(item => (
              <div className="cart-item" key={item.product_id}>
                <img src={IMG_SINGLE_BOX} alt={item.name} className="cart-item-img" />
                <div className="cart-item-info">
                  <div className="cart-item-name">{item.name}</div>
                  <div className="cart-item-price">₹{item.price} each</div>
                </div>
                <div className="cart-item-qty">
                  <button onClick={() => addToCart({ id: item.product_id, name: item.name, price: item.price }, Math.max(1, item.quantity - 1))}>−</button>
                  <span>{item.quantity}</span>
                  <button onClick={() => addToCart({ id: item.product_id, name: item.name, price: item.price }, item.quantity + 1)}>+</button>
                </div>
                <div className="cart-item-total">₹{item.price * item.quantity}</div>
                <button className="remove-btn" onClick={() => removeFromCart(item.product_id)}>✕</button>
              </div>
            ))}
          </div>

          <div className="cart-summary">
            <h3>Order Summary</h3>
            <div className="summary-row"><span>Subtotal</span><span>₹{cartTotal}</span></div>
            <div className="summary-row"><span>Shipping</span><span className="free">FREE</span></div>
            <div className="summary-row total"><span>Total</span><span>₹{cartTotal}</span></div>
            {!user && (
              <button className="btn-login-small" onClick={loginWithGoogle}>
                Sign in with Google to checkout
              </button>
            )}
            <button className="btn-primary full-width" style={{marginTop:"1rem"}}
              onClick={() => user ? navigate("/checkout") : loginWithGoogle()}>
              {user ? "Proceed to Checkout" : "Sign In to Checkout"}
            </button>
            <div className="cart-guarantee">🛡️ Secure checkout · Free returns</div>
          </div>
        </div>
      </div>
    </main>
  );
}
