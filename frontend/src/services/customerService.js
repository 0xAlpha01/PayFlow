import { api } from "./api";
export const getCustomers = (params = {}) => api.get("/customers", { params });
export const getCustomer = (id) => api.get(`/customers/${id}`);
export const createCustomer = (data) => api.post("/customers", data);
export const updateCustomer = (id, data) => api.put(`/customers/${id}`, data);
export const deleteCustomer = (id) => api.delete(`/customers/${id}`);
export const getCustomerPayments = (id) => api.get(`/customers/${id}/payments`);
export const getCustomerTransactions = (id) => api.get(`/customers/${id}/transactions`);
