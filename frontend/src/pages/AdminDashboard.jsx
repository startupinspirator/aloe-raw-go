import { useState, useEffect } from "react";
import axios from "axios";
import { useAuth } from "../context/AuthContext";
import API_URL from "../config";
import { useNavigate } from "react-router-dom";

export default function AdminDashboard() {
  const { user, loading } = useAuth();
  const navigate = useNavigate();
  const [orders, setOrders] = useState([]);
  const [fetching, setFetching] = useState(true);

  useEffect(() => {
    if (!loading) {
      if (!user || user.role !== "admin") {
        navigate("/");
        return;
      }
      fetchOrders();
    }
  }, [user, loading, navigate]);

  const fetchOrders = async () => {
    try {
      const res = await axios.get(`${API_URL}/api/admin/orders`, { withCredentials: true });
      setOrders(res.data.orders || []);
    } catch (err) {
      console.error(err);
    } finally {
      setFetching(false);
    }
  };

  if (loading || fetching) return <div className="loading">Loading Admin...</div>;

  return (
    <main className="container page-admin fade-in">
      <div className="section-header">
        <h1>Admin Dashboard</h1>
        <p>Manage all platform orders and products</p>
      </div>

      <div className="admin-content">
        <section className="admin-card full-width">
          <h2>All Orders</h2>
          {orders.length === 0 ? (
            <p>No orders on the platform yet.</p>
          ) : (
            <div className="admin-table-wrapper">
              <table className="admin-table">
                <thead>
                  <tr>
                    <th>Order ID</th>
                    <th>Customer (ID)</th>
                    <th>Date</th>
                    <th>Status</th>
                    <th>Total</th>
                  </tr>
                </thead>
                <tbody>
                  {orders.map(order => (
                    <tr key={order.id}>
                      <td>#{order.id}</td>
                      <td>{order.shipping_name} (User {order.user_id})</td>
                      <td>{new Date(order.created_at).toLocaleDateString()}</td>
                      <td>
                        <span className={`status-badge ${order.status}`}>
                          {order.status}
                        </span>
                      </td>
                      <td>₹{order.total_amount}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
