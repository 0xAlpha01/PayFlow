import { api } from "./api";
export const getNotifications = (params = {}) => api.get("/notifications", { params });
export const getUnreadNotificationCount = () => api.get("/notifications/unread-count");
export const markNotificationRead = (id) => api.post(`/notifications/${id}/read`);
export const markAllNotificationsRead = () => api.post("/notifications/read-all");
