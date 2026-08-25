import { api } from "./api";
export const getRevenueAnalytics = (period = "30d") =>
  api.get("/analytics/revenue", { params: { period } });
export const getTransactionAnalytics = (period = "30d") =>
  api.get("/analytics/transactions", { params: { period } });
export const getCustomerAnalytics = (period = "30d") =>
  api.get("/analytics/customers", { params: { period } });
export const getPaymentAnalytics = (period = "30d") =>
  api.get("/analytics/payments", { params: { period } });
