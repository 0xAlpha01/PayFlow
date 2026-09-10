import { api } from "./api";
export const getWalletBalance = () => api.get("/wallet/balance");
export const sendTransfer = (data) => api.post("/transfers", data);
export const getTransfers = (params = {}) => api.get("/transfers", { params });
