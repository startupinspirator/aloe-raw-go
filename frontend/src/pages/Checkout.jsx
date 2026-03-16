import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import { useEffect } from "react";
import { useCart } from "../context/CartContext";
import { useAuth } from "../context/AuthContext";
import API_URL from "../config";
import { loadScript } from "../utils/loadScript";

export default function Checkout() {
  const { cart, cartTotal, clearCart } = useCart();
  const { user, loginWithGoogle } = useAuth();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({});
  const [form, setForm] = useState({
    name: user?.name || "", phone: "", address: "", city: "", pincode: "",
  });

  useEffect(() => {
    loadScript("https://checkout.razorpay.com/v1/checkout.js")
      .catch(err => console.error("Razorpay SDK load failed:", err));
  }, []);

  if (!user) {
    return (
      <main className="checkout-page">
        <div className="empty-cart">
          <h2>Please sign in to checkout</h2>
          <button className="btn-primary" onClick={loginWithGoogle}>Sign in with Google</button>
        </div>
      </main>
    );
  }

  const validate = () => {
    const e = {};
    if (!form.name.trim()) e.name = "Name required";
    if (!/^[6-9]\d{9}$/.test(form.phone)) e.phone = "Enter valid 10-digit mobile";
    if (!form.address.trim()) e.address = "Address required";
    if (!form.city.trim()) e.city = "City required";
    if (!/^\d{6}$/.test(form.pincode)) e.pincode = "Enter valid 6-digit pincode";
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleChange = e => {
    setForm(f => ({ ...f, [e.target.name]: e.target.value }));
    setErrors(er => ({ ...er, [e.target.name]: "" }));
  };

  const handlePayment = async () => {
    if (!validate()) return;
    setLoading(true);
    try {
      const { data } = await axios.post(`${API_URL}/api/payment/create-order`,
        { shipping: form }, { withCredentials: true });

      const options = {
        key: data.keyId,
        amount: data.amount * 100,
        currency: "INR",
        name: "Aloé Raw",
        description: "Farm Fresh Aloe Vera Cubes",
        order_id: data.razorpayOrderId,
        prefill: { name: user.name, email: user.email, contact: form.phone },
        theme: { color: "#4a6741" },
        handler: async (response) => {
          try {
            await axios.post(`${API_URL}/api/payment/verify`, {
              razorpay_order_id: response.razorpay_order_id,
              razorpay_payment_id: response.razorpay_payment_id,
              razorpay_signature: response.razorpay_signature,
              orderId: data.orderId,
            }, { withCredentials: true });
            clearCart();
            navigate(`/orders?success=${data.orderId}`);
          } catch {
            alert("Payment verification failed. Contact support.");
          }
        },
        modal: { ondismiss: () => setLoading(false) },
      };
      const rzp = new window.Razorpay(options);
      rzp.open();
    } catch (err) {
      alert(err.response?.data?.error || "Payment failed. Try again.");
      setLoading(false);
    }
  };

  const Field = ({ label, name, placeholder, maxLength }) => (
    <div className="form-group">
      <label>{label}</label>
      <input name={name} value={form[name]} onChange={handleChange}
        placeholder={placeholder} maxLength={maxLength} />
      {errors[name] && <span className="form-error">{errors[name]}</span>}
    </div>
  );

  return (
    <main className="checkout-page">
      <div className="checkout-container">
        <h1 className="page-title">Checkout</h1>
        <div className="checkout-layout">
          <div className="shipping-form">
            <h3>Shipping Details</h3>
            <div className="form-row">
              <Field label="Full Name" name="name" placeholder="Ayush Kumar" />
              <Field label="Mobile Number" name="phone" placeholder="9876543210" maxLength={10} />
            </div>
            <Field label="Address" name="address" placeholder="House/Flat No, Street, Area" />
            <div className="form-row">
              <Field label="City" name="city" placeholder="Bhubaneswar" />
              <Field label="Pincode" name="pincode" placeholder="751001" maxLength={6} />
            </div>
          </div>

          <div className="checkout-summary">
            <h3>Order Summary</h3>
            {cart.map(item => (
              <div className="summary-item" key={item.product_id}>
                <span>{item.name} × {item.quantity}</span>
                <span>₹{item.price * item.quantity}</span>
              </div>
            ))}
            <div className="summary-divider" />
            <div className="summary-row"><span>Subtotal</span><span>₹{cartTotal}</span></div>
            <div className="summary-row"><span>Shipping</span><span className="free">FREE</span></div>
            <div className="summary-row total"><span>Total</span><span>₹{cartTotal}</span></div>
            <button className="btn-primary full-width" style={{marginTop:"1.5rem"}}
              onClick={handlePayment} disabled={loading || cart.length === 0}>
              {loading ? "Opening Payment..." : `Pay ₹${cartTotal} with Razorpay`}
            </button>
            <p className="secure-note">🔒 Secured by Razorpay · UPI, Cards, NetBanking</p>
          </div>
        </div>
      </div>
    </main>
  );
}
