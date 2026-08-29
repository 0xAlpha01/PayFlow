import { api } from "./api";
export const getSubscriptions = (params = {}) =>
  api.get("/subscriptions", { params });
export const getSubscription = (id) => api.get(`/subscriptions/${id}`);
export const getSubscriptionSummary = () => api.get("/subscriptions/summary");
export const createSubscription = (data) => api.post("/subscriptions", data);
export const updateSubscription = (id, data) => api.put(`/subscriptions/${id}`, data);
export const cancelSubscription = (id) =>
  api.post(`/subscriptions/${id}/cancel`);
