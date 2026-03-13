import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useCart } from "../context/CartContext";

export default function Navbar() {
  const { user, logout, loginWithGoogle } = useAuth();
  const { cartCount } = useCart();
  const [menuOpen, setMenuOpen] = useState(false);
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate("/");
  };

  return (
    <nav className="navbar">
      <Link to="/" className="nav-logo">Aloé <span>Raw</span></Link>

      <ul className="nav-links">
        <li><Link to="/">Home</Link></li>
        <li><Link to="/product/1">Shop</Link></li>
        {user && <li><Link to="/orders">My Orders</Link></li>}
      </ul>

      <div className="nav-right">
        <Link to="/cart" className="cart-icon">
          🛒
          {cartCount > 0 && <span className="cart-badge">{cartCount}</span>}
        </Link>

        {user ? (
          <div className="user-menu">
            {user.avatar
              ? <img src={user.avatar} alt={user.name} className="user-avatar" />
              : <div className="user-avatar-fallback">{user.name?.[0]}</div>
            }
            <div className="user-dropdown">
              <span className="user-name">{user.name}</span>
              <Link to="/orders">My Orders</Link>
              <button onClick={handleLogout}>Sign Out</button>
            </div>
          </div>
        ) : (
          <button className="btn-login" onClick={loginWithGoogle}>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
            </svg>
            Sign in with Google
          </button>
        )}

        <button className="hamburger" onClick={() => setMenuOpen(!menuOpen)}>☰</button>
      </div>

      {menuOpen && (
        <div className="mobile-menu">
          <Link to="/" onClick={() => setMenuOpen(false)}>Home</Link>
          <Link to="/product/1" onClick={() => setMenuOpen(false)}>Shop</Link>
          <Link to="/cart" onClick={() => setMenuOpen(false)}>Cart {cartCount > 0 && `(${cartCount})`}</Link>
          {user && <Link to="/orders" onClick={() => setMenuOpen(false)}>My Orders</Link>}
          {user
            ? <button onClick={() => { handleLogout(); setMenuOpen(false); }}>Sign Out</button>
            : <button onClick={() => { loginWithGoogle(); setMenuOpen(false); }}>Sign in with Google</button>
          }
        </div>
      )}
    </nav>
  );
}
