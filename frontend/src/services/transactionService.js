import { api } from "./api";
export const getTransactions = (params = {}) =>
  api.get("/transactions", { params });
export const getTransaction = (id) => api.get(`/transactions/${id}`);
