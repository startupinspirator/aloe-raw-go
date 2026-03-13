import { useEffect, useState } from "react";
import { useSearchParams, Link } from "react-router-dom";
import axios from "axios";
import { useAuth } from "../context/AuthContext";
import API_URL from "../config";

export default function Orders() {
  const { user, loginWithGoogle } = useAuth();
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [params] = useSearchParams();
  const successId = params.get("success");

  useEffect(() => {
    if (!user) return;
    axios.get(`${API_URL}/api/orders`, { withCredentials: true })
      .then(res => setOrders(res.data || []))
      .finally(() => setLoading(false));
  }, [user]);

  if (!user) return (
    <main className="orders-page">
      <div className="empty-cart">
        <h2>Please sign in to view orders</h2>
        <button className="btn-primary" onClick={loginWithGoogle}>Sign in with Google</button>
      </div>
    </main>
  );

  return (
    <main className="orders-page">
      <div className="orders-container">
        {successId && (
          <div className="success-banner">
            🎉 Order #{successId} confirmed! Your fresh aloe vera is on its way.
          </div>
        )}
        <h1 className="page-title">My Orders</h1>
        {loading ? <div className="loading">Loading orders...</div>
          : orders.length === 0 ? (
            <div className="empty-cart">
              <p>No orders yet.</p>
              <Link to="/product/1" className="btn-primary">Shop Now</Link>
            </div>
          ) : (
            <div className="orders-list">
              {orders.map(order => (
                <div className="order-card" key={order.id}>
                  <div className="order-header">
                    <div>
                      <span className="order-id">Order #{order.id}</span>
                      <span className="order-date">
                        {new Date(order.created_at).toLocaleDateString("en-IN",
                          { day:"numeric", month:"long", year:"numeric" })}
                      </span>
                    </div>
                    <span className={`order-status ${order.status}`}>{order.status.toUpperCase()}</span>
                  </div>
                  <div className="order-items">
                    {order.items?.map(item => (
                      <div className="order-item" key={item.id}>
                        <span>{item.name} × {item.quantity}</span>
                        <span>₹{item.price_at_purchase * item.quantity}</span>
                      </div>
                    ))}
                  </div>
                  <div className="order-footer">
                    <span>Shipping to: {order.shipping_city}</span>
                    <span className="order-total">Total: ₹{order.total_amount}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
      </div>
    </main>
  );
}
