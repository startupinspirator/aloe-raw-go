// In development: Vite proxies /api to localhost:8080
// In production (GitHub Pages): calls your Render backend directly
const isDev = import.meta.env.DEV;

export const API_URL = isDev
  ? "" // Vite proxy handles it
  : import.meta.env.VITE_API_URL;

if (!isDev && !API_URL) {
  console.error("CRITICAL: VITE_API_URL is missing. API calls will fail. Check your environment variables.");
}

export default API_URL;
