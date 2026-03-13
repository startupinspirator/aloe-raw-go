import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import axios from "axios";
import { useCart } from "../context/CartContext";
import { useAuth } from "../context/AuthContext";
import { IMG_SINGLE_BOX, IMG_STACKED_BOXES, IMG_WOMAN_APPLYING } from "../assets/images";
import API_URL from "../config";

export default function Product() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { addToCart } = useCart();
  const { loginWithGoogle } = useAuth();
  const [product, setProduct] = useState(null);
  const [qty, setQty] = useState(1);
  const [adding, setAdding] = useState(false);
  const [added, setAdded] = useState(false);
  const [activeImg, setActiveImg] = useState(0);
  const [error, setError] = useState("");

  const images = [IMG_SINGLE_BOX, IMG_STACKED_BOXES, IMG_WOMAN_APPLYING];

  useEffect(() => {
    axios.get(`${API_URL}/api/products/${id}`)
      .then(res => setProduct(res.data))
      .catch(() => setError("Product not found"));
  }, [id]);

  const handleAddToCart = async () => {
    setAdding(true);
    await addToCart(product, qty);
    setAdding(false);
    setAdded(true);
    setTimeout(() => setAdded(false), 2500);
  };

  const handleBuyNow = async () => {
    await addToCart(product, qty);
    navigate("/checkout");
  };

  if (error) return <div className="loading">{error}</div>;
  if (!product) return <div className="loading">Loading product...</div>;

  const discount = Math.round((1 - product.price / product.original_price) * 100);

  return (
    <main className="product-page">
      <div className="product-container">
        {/* Gallery */}
        <div className="product-gallery">
          <div className="gallery-thumbs">
            {images.map((img, i) => (
              <img key={i} src={img} alt={`View ${i+1}`}
                className={`thumb ${activeImg === i ? "active" : ""}`}
                onClick={() => setActiveImg(i)} />
            ))}
          </div>
          <div className="gallery-main">
            <img src={images[activeImg]} alt={product.name} />
            <div className="discount-badge">-{discount}%</div>
          </div>
        </div>

        {/* Info */}
        <div className="product-info">
          <div className="section-label">Featured Product</div>
          <h1 className="product-title">{product.name}</h1>
          <p className="product-tagline">{product.tagline}</p>
          <div className="rating-row">
            <span className="stars">★★★★★</span>
            <span className="rating-count">4.9 · 128 reviews</span>
          </div>
          <div className="pricing-row">
            <span className="price-main">₹{product.price}</span>
            {product.original_price > 0 && <span className="price-old">₹{product.original_price}</span>}
            <span className="price-note">Launch Price — {discount}% off</span>
          </div>
          <p className="product-desc">{product.description}</p>
          <ul className="product-features">
            {["10 individually sealed cubes per pack",
              "Harvested & shipped within 24 hours",
              "Suitable for skin, hair, and internal use",
              "Refrigerate up to 7 days · Freeze 3 months",
              "Certified organically grown, zero pesticides",
            ].map(f => <li key={f}>{f}</li>)}
          </ul>
          <div className="qty-add">
            <div className="qty-control">
              <button className="qty-btn" onClick={() => setQty(q => Math.max(1, q-1))}>−</button>
              <span className="qty-num">{qty}</span>
              <button className="qty-btn" onClick={() => setQty(q => Math.min(10, q+1))}>+</button>
            </div>
            <button className={`cart-btn ${added ? "added" : ""}`} onClick={handleAddToCart} disabled={adding}>
              {added ? "✓ Added to Cart" : adding ? "Adding..." : "Add to Cart"}
            </button>
          </div>
          <button className="btn-primary full-width" onClick={handleBuyNow}>
            Buy Now — ₹{product.price * qty}
          </button>
          <div className="product-guarantee">
            🛡️ &nbsp;100% freshness guaranteed. Full refund if not satisfied.
          </div>
        </div>
      </div>
    </main>
  );
}
