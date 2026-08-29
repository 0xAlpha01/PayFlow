import { api } from "./api";
export const getPayments = (params = {}) => api.get("/payments", { params });
export const getPayment = (id) => api.get(`/payments/${id}`);
export const getPaymentSummary = () => api.get("/payments/summary");
export const initializePayment = (data) =>
  api.post("/payments/initialize", data);
export const verifyPayment = (reference) =>
  api.get(`/payments/verify/${reference}`);
