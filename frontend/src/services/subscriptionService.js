import { api } from "./api";
export const getSubscriptions = (params = {}) =>
  api.get("/subscriptions", { params });
export const getSubscription = (id) => api.get(`/subscriptions/${id}`);
export const createSubscription = (data) => api.post("/subscriptions", data);
export const cancelSubscription = (id) =>
  api.post(`/subscriptions/${id}/cancel`);
