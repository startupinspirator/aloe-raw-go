import { useState, useEffect } from "react";
import axios from "axios";
import { useAuth } from "../context/AuthContext";
import API_URL from "../config";
import { useNavigate } from "react-router-dom";

export default function Profile() {
  const { user, loading } = useAuth();
  const navigate = useNavigate();
  const [orders, setOrders] = useState([]);
  const [fetching, setFetching] = useState(true);

  useEffect(() => {
    if (!loading) {
      if (!user) {
        navigate("/");
        return;
      }
      fetchMyOrders();
    }
  }, [user, loading, navigate]);

  const fetchMyOrders = async () => {
    try {
      const res = await axios.get(`${API_URL}/api/orders`, { withCredentials: true });
      setOrders(res.data.orders || []);
    } catch (err) {
      console.error(err);
    } finally {
      setFetching(false);
    }
  };

  if (loading || fetching) return <div className="loading">Loading Profile...</div>;

  return (
    <main className="container page-profile fade-in">
      <div className="profile-header">
        {user.avatar ? (
          <img src={user.avatar} alt="Avatar" className="profile-avatar-large" />
        ) : (
          <div className="profile-avatar-fallback">{user.name?.[0]}</div>
        )}
        <div className="profile-info">
          <h1>{user.name}</h1>
          <p>{user.email}</p>
          {user.role === 'admin' && <span className="role-badge admin">Admin</span>}
        </div>
      </div>

      <section className="profile-orders">
        <h2>My Order History</h2>
        {orders.length === 0 ? (
          <div className="empty-state">
            <p>You haven't placed any orders yet.</p>
            <button className="btn-primary" onClick={() => navigate("/product/1")}>Shop Now</button>
          </div>
        ) : (
          <div className="orders-list">
            {orders.map(order => (
              <div key={order.id} className="order-card">
                <div className="order-header">
                  <div>
                    <span className="order-id">Order #{order.id}</span>
                    <span className="order-date">{new Date(order.created_at).toLocaleDateString()}</span>
                  </div>
                  <span className={`order-status ${order.status}`}>{order.status}</span>
                </div>
                <div className="order-items">
                  {order.items.map(item => (
                    <div key={item.id} className="order-item">
                      <span>{item.quantity}x {item.name}</span>
                      <span>₹{item.price_at_purchase}</span>
                    </div>
                  ))}
                </div>
                <div className="order-footer">
                  <strong>Total: ₹{order.total_amount}</strong>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
    </main>
  );
}
