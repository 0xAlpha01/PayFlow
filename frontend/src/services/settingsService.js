import { api } from "./api";
export const getProfile = () => api.get("/settings/profile");
export const updateProfile = (data) => api.put("/settings/profile", data);
export const getBusiness = () => api.get("/settings/business");
export const updateBusiness = (data) => api.put("/settings/business", data);
export const changePassword = (data) =>
  api.post("/settings/change-password", data);
export const getNotifications = () => api.get("/settings/notifications");
export const updateNotifications = (data) =>
  api.put("/settings/notifications", data);
