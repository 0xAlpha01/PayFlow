import axios from "axios";

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1",
  headers: { "Content-Type": "application/json" },
  withCredentials: true,
});

let accessToken = null;
let refreshToken = null;
let onRefresh = null;
let onLogout = null;
let refreshing = null;

export function configureAuth(
  nextAccessToken,
  nextRefreshToken,
  refreshHandler,
  logoutHandler,
) {
  accessToken = nextAccessToken;
  refreshToken = nextRefreshToken;
  onRefresh = refreshHandler;
  onLogout = logoutHandler;
}

api.interceptors.request.use((config) => {
  if (accessToken) config.headers.Authorization = `Bearer ${accessToken}`;
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config;
    if (
      error.response?.status !== 401 ||
      original?._retry ||
      !refreshToken ||
      original?.url?.includes("/auth/refresh")
    )
      return Promise.reject(error);
    original._retry = true;
    try {
      if (!refreshing)
        refreshing = onRefresh ? onRefresh() : Promise.reject(error);
      const nextToken = await refreshing;
      refreshing = null;
      original.headers.Authorization = `Bearer ${nextToken}`;
      return api(original);
    } catch (refreshError) {
      refreshing = null;
      onLogout?.();
      window.location.assign("/login");
      return Promise.reject(refreshError);
    }
  },
);
