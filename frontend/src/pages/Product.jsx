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
  const [product, setProduct] = useState(null);
  const [reviews, setReviews] = useState([]);
  const [newReview, setNewReview] = useState({ rating: 5, comment: "" });
  const [reviewError, setReviewError] = useState("");
  const { user, loginWithGoogle } = useAuth();
  
  const [qty, setQty] = useState(1);
  const [adding, setAdding] = useState(false);
  const [added, setAdded] = useState(false);
  const [activeImg, setActiveImg] = useState(0);
  const [error, setError] = useState("");

  const images = [IMG_SINGLE_BOX, IMG_STACKED_BOXES, IMG_WOMAN_APPLYING];

  useEffect(() => {
    axios.get(`${API_URL}/api/products/${id}`)
      .then(res => {
        // Ensure we got a valid JSON object and not an HTML fallback
        if (typeof res.data !== "object" || res.data === null) {
          throw new Error("Invalid response format");
        }
        setProduct(res.data);
      })
      .catch((err) => {
        console.error("Fetch product failed:", err);
        setError("Product not found or Connection Error");
      });
      
    fetchReviews();
  }, [id]);

  const fetchReviews = () => {
    axios.get(`${API_URL}/api/products/${id}/reviews`)
      .then(res => setReviews(res.data.reviews || []))
      .catch(err => console.error("Could not fetch reviews", err));
  };

  const submitReview = async () => {
    try {
      setReviewError("");
      await axios.post(`${API_URL}/api/products/${id}/reviews`, newReview, { withCredentials: true });
      setNewReview({ rating: 5, comment: "" });
      fetchReviews(); // Refresh
    } catch (err) {
      setReviewError(err.response?.data?.error || "Failed to submit review");
    }
  };

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

      {/* Reviews Section */}
      <div className="reviews-section">
        <h2>Customer Reviews</h2>
        <div className="reviews-grid">
          {reviews.length === 0 ? (
            <p className="no-reviews">No reviews yet. Be the first to review!</p>
          ) : (
            reviews.map(r => (
              <div key={r.id} className="review-card">
                <div className="review-header">
                  {r.user_avatar ? (
                    <img src={r.user_avatar} alt="User" />
                  ) : (
                    <div className="avatar-fallback">{r.user_name?.[0]}</div>
                  )}
                  <div>
                    <strong>{r.user_name}</strong>
                    <div className="review-stars">{"★".repeat(r.rating)}{"☆".repeat(5-r.rating)}</div>
                  </div>
                  <span className="review-date">{new Date(r.created_at).toLocaleDateString()}</span>
                </div>
                <p>{r.comment}</p>
              </div>
            ))
          )}
        </div>

        {user ? (
          <div className="review-form">
            <h3>Write a Review</h3>
            {reviewError && <p className="error">{reviewError}</p>}
            <div className="rating-select">
              <span>Rating: </span>
              {[1, 2, 3, 4, 5].map(num => (
                <button
                  key={num}
                  type="button"
                  className={newReview.rating >= num ? "star selected" : "star"}
                  onClick={() => setNewReview({ ...newReview, rating: num })}
                >★</button>
              ))}
            </div>
            <textarea
              placeholder="What did you think of our product?"
              value={newReview.comment}
              onChange={e => setNewReview({ ...newReview, comment: e.target.value })}
            />
            <button className="btn-primary" onClick={submitReview}>Submit Review</button>
          </div>
        ) : (
          <div className="review-login-prompt">
            <p>Please log in to write a review.</p>
            <button className="btn-secondary" onClick={loginWithGoogle}>Sign In</button>
          </div>
        )}
      </div>
    </main>
  );
}
