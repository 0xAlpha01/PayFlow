import { api } from "./api";

export const authService = {
  signup: (data) => api.post("/auth/signup", data),
  login: (data) => api.post("/auth/login", data),
  logout: (refreshToken) => api.post("/auth/logout", { refreshToken }),
  refreshToken: (refreshToken) => api.post("/auth/refresh", { refreshToken }),
  getCurrentUser: () => api.get("/auth/me"),
  verifyEmail: (data) => api.post("/auth/verify-email", data),
  resendEmailVerification: (email) =>
    api.post("/auth/resend-email-verification", { email }),
  sendPhoneCode: (phone) => api.post("/auth/send-phone-code", { phone }),
  resendPhoneCode: (phone) => api.post("/auth/resend-phone-code", { phone }),
  verifyPhone: (code) => api.post("/auth/verify-phone", { code }),
  forgotPassword: (email) => api.post("/auth/forgot-password", { email }),
  resetPassword: (data) => api.post("/auth/reset-password", data),
};
