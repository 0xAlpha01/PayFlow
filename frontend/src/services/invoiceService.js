import { api } from "./api";
export const getInvoices = (params = {}) => api.get("/invoices", { params });
export const getInvoice = (id) => api.get(`/invoices/${id}`);
export const createInvoice = (data) => api.post("/invoices", data);
export const updateInvoice = (id, data) => api.put(`/invoices/${id}`, data);
export const deleteInvoice = (id) => api.delete(`/invoices/${id}`);
export const sendInvoice = (id) => api.post(`/invoices/${id}/send`);
export const markInvoicePaid = (id) => api.post(`/invoices/${id}/mark-paid`);
