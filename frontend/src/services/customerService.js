import { customers } from "../data/customers";
export const getCustomers = async () => customers;
export const getCustomer = async (id) =>
  customers.find((customer) => customer.id === id);
