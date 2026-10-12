const LOCAL_API_BASE_URL = "http://localhost:8080/api/v1";

// Render service URLs are commonly entered either as the host itself or with an
// `/api` suffix. The Go server exposes its API under `/api/v1`, so normalize
// those two forms before any request or websocket URL is created.
export function resolveApiBaseUrl(value = import.meta.env.VITE_API_URL) {
  const baseUrl = (value || LOCAL_API_BASE_URL).trim().replace(/\/+$/, "");

  try {
    const url = new URL(baseUrl);
    if (url.pathname === "" || url.pathname === "/") url.pathname = "/api/v1";
    else if (url.pathname === "/api") url.pathname = "/api/v1";
    return url.toString().replace(/\/$/, "");
  } catch {
    // Preserve a deliberately supplied relative API path for development.
    return baseUrl;
  }
}

export const apiBaseUrl = resolveApiBaseUrl();
