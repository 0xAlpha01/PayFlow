import { api } from "./api";
export const searchChatUsers = (q) => api.get("/chat/users", { params: { q } });
export const startConversation = (userId) => api.post("/chat/conversations", { userId });
export const getConversations = () => api.get("/chat/conversations");
export const getMessages = (conversationId, params = {}) =>
  api.get(`/chat/conversations/${conversationId}/messages`, { params });
export const sendMessage = (conversationId, content) =>
  api.post(`/chat/conversations/${conversationId}/messages`, { content });
export const markConversationRead = (conversationId) =>
  api.post(`/chat/conversations/${conversationId}/read`);
export const getUnreadMessageCount = () => api.get("/chat/unread-count");
