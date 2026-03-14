import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import axios from "axios";
import API_URL from "../config";
import { IMG_WOMAN_APPLYING, IMG_SINGLE_BOX, IMG_STACKED_BOXES } from "../assets/images";

export default function Home() {
  const navigate = useNavigate();
  const [categories, setCategories] = useState([]);

  useEffect(() => {
    axios.get(`${API_URL}/api/categories`)
      .then(res => setCategories(res.data.categories || []))
      .catch(err => console.error("Error fetching categories:", err));
  }, []);

  return (
    <main className="home">

      {/* ── HERO ── */}
      <section className="hero">
        <div className="hero-left">
          <div className="hero-tag">Farm to your hands in 24 hours</div>
          <h1 className="hero-title">Pure aloe,<br /><em>exactly as</em><br />nature made it.</h1>
          <p className="hero-subtitle">
            We cut, cube, and ship directly from our farm. No preservatives, no fillers,
            no processing. Just raw, living aloe vera.
          </p>
          <div className="hero-actions">
            <button className="btn-primary" onClick={() => navigate("/product/1")}>Shop Now — ₹299</button>
            <a href="#story" className="btn-ghost">Our Story</a>
          </div>
          <div className="hero-badges">
            <div className="badge"><span className="badge-icon">🌿</span>100% Organic</div>
            <div className="badge"><span className="badge-icon">⚡</span>Same Day Harvest</div>
            <div className="badge"><span className="badge-icon">📦</span>Free Delivery</div>
          </div>
        </div>
        <div className="hero-right">
          <div className="hero-img-frame">
            <img src={IMG_SINGLE_BOX} alt="Aloé Raw product box" />
            <div className="hero-img-tag">
              <span className="big">₹299</span>
              <span className="small">Per Pack · 10 Cubes</span>
            </div>
          </div>
        </div>
      </section>

      {/* ── MARQUEE ── */}
      <div className="marquee-strip">
        <div className="marquee-inner">
          {Array(2).fill(["Farm Fresh","100% Pure Aloe Vera","Zero Processing","Delivered in 24 Hours","No Preservatives","Organically Grown"]).flat()
            .map((t, i) => <span key={i} className="marquee-item">{t} <span className="marquee-dot">✦</span></span>)}
        </div>
      </div>

      {/* ── SHOP BY CATEGORY ── */}
      {categories.length > 0 && (
        <section className="categories-section">
          <div className="section-label centered">Explore</div>
          <h2 className="categories-title">Shop by <em>Category</em></h2>
          <div className="categories-grid">
            {categories.map(cat => (
              <div key={cat.id} className="category-card" onClick={() => navigate("/product/1")}>
                <h3>{cat.name}</h3>
                <p>{cat.description}</p>
                <div className="category-link">View Products →</div>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* ── STORY ── */}
      <section className="story-section" id="story">
        <div className="story-text">
          <div className="section-label">Our Story</div>
          <h2 className="story-title">From our farm<br />to your <em>hands</em><br />in 24 hours.</h2>
          <p className="story-body">Most aloe vera products sit in warehouses for months. By the time they reach you, the beneficial compounds have degraded. We started Aloé Raw to change that.</p>
          <p className="story-body">Our farm in Odisha has been growing aloe vera for three generations. We harvest at peak maturity, cube the gel, seal it, and ship — in less than a day.</p>
          <div className="story-stats">
            <div><span className="stat-num">3</span><span className="stat-label">Generations farming</span></div>
            <div><span className="stat-num">24h</span><span className="stat-label">Farm to delivery</span></div>
            <div><span className="stat-num">0</span><span className="stat-label">Additives</span></div>
          </div>
        </div>
        <div className="story-images">
          <img src={IMG_STACKED_BOXES} alt="Aloé Raw farm" className="story-img-main" />
          <img src={IMG_WOMAN_APPLYING} alt="Applying aloe" className="story-img-secondary" />
        </div>
      </section>

      {/* ── BENEFITS ── */}
      <section className="benefits-section">
        <div className="section-label centered">What Raw Aloe Does</div>
        <h2 className="benefits-title">Nature's most <em>versatile</em> skincare ingredient.</h2>
        <div className="benefits-grid">
          {[
            { icon:"✨", name:"Deep Moisturisation", text:"Fresh aloe gel penetrates three layers of skin, delivering hydration that processed gels simply cannot match." },
            { icon:"🌡️", name:"Soothing & Cooling", text:"Instantly calms sunburn, redness, and inflammation. Works on contact — no waiting for absorption." },
            { icon:"💆", name:"Hair & Scalp Health", text:"Raw aloe enzymes remove dead scalp cells and condition hair follicles for visibly healthier growth." },
            { icon:"🫀", name:"Internal Benefits", text:"Food-grade quality. Blend into juices or consume directly for digestion support and immune boosting." },
          ].map(b => (
            <div className="benefit-card" key={b.name}>
              <span className="benefit-icon">{b.icon}</span>
              <div className="benefit-name">{b.name}</div>
              <p className="benefit-text">{b.text}</p>
            </div>
          ))}
        </div>
      </section>

      {/* ── TESTIMONIALS ── */}
      <section className="testimonials-section">
        <div className="section-label centered">Real Results</div>
        <h2 className="testi-title">What our customers say</h2>
        <div className="testi-grid">
          {[
            { name:"Priya M.", loc:"Bangalore", text:"I've tried every aloe gel on the market. Nothing comes close. You can feel the difference — it's genuinely alive in a way bottled products are not." },
            { name:"Rohit K.", loc:"Mumbai",    text:"Used it on my sunburn the day it arrived. Completely different from any product I've used before. Healed in two days. Genuinely remarkable." },
            { name:"Ananya S.", loc:"Hyderabad", text:"The gel is crystal clear and perfectly set. Even fresher than having my own plant at home. Incredible product, fast delivery." },
          ].map(t => (
            <div className="testi-card" key={t.name}>
              <div className="testi-stars">★★★★★</div>
              <p className="testi-text">"{t.text}"</p>
              <div className="testi-author">{t.name}</div>
              <div className="testi-location">{t.loc} · Verified Buyer</div>
            </div>
          ))}
        </div>
      </section>

      {/* ── CTA ── */}
      <section className="cta-section">
        <div className="section-label cta-label centered">Limited Launch Offer</div>
        <h2 className="cta-title">Try pure aloe,<br /><em>risk free.</em></h2>
        <p className="cta-sub">Order today and experience genuinely fresh aloe vera. Not satisfied? Full refund, no questions asked.</p>
        <button className="btn-light" onClick={() => navigate("/product/1")}>Order Now — ₹299</button>
      </section>

      {/* ── FOOTER ── */}
      <footer>
        <div className="footer-top">
          <div className="footer-brand">
            <span className="nav-logo footer-logo">Aloé <span>Raw</span></span>
            <span className="footer-tagline">Farm Fresh, 100% Pure. Zero Processing.</span>
          </div>
          <div className="footer-col"><h4>Shop</h4><ul><li><a href="/aloe-raw/product/1">Aloe Vera Cubes</a></li></ul></div>
          <div className="footer-col"><h4>Company</h4><ul><li><a href="#story">Our Farm</a></li><li><a href="#story">About</a></li></ul></div>
          <div className="footer-col"><h4>Help</h4><ul><li><a href="#">Shipping</a></li><li><a href="#">Returns</a></li><li><a href="#">FAQ</a></li></ul></div>
        </div>
        <div className="footer-bottom">
          <span>© 2026 Aloé Raw. All rights reserved.</span>
          <span>Made with care in Odisha, India 🌿</span>
        </div>
      </footer>
    </main>
  );
}
