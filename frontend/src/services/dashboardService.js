import { api } from "./api";

export const dashboardService = {
  getDashboardSummary: () => api.get("/dashboard/summary"),
  getRevenue: (period = "30d") =>
    api.get("/dashboard/revenue", { params: { period } }),
  getRecentTransactions: () => api.get("/dashboard/transactions"),
  getRecentCustomers: () => api.get("/dashboard/customers"),
};
