import { useEffect, useState } from "react";
import { getTransactions } from "../services/transactionService";
import { formatCurrency } from "../utils/currency";

export default function TransactionsPage() {
  const [rows, setRows] = useState([]), [loading, setLoading] = useState(true), [error, setError] = useState("");
  const load = async () => { setLoading(true); setError(""); try { const response = await getTransactions({ page: 1, limit: 100 }); const data = response.data.data; setRows(Array.isArray(data?.items) ? data.items : Array.isArray(data) ? data : []); } catch { setError("Unable to load transactions."); } finally { setLoading(false); } };
  useEffect(() => { load(); }, []);
  return <><div className="page-heading"><div><h2>Transactions</h2><p>Review your transaction history.</p></div></div><div className="panel table-panel">{loading ? <div className="data-loading">Loading transactions...</div> : error ? <div className="data-empty">{error} <button onClick={load}>Retry</button></div> : rows.length === 0 ? <div className="data-empty">No transactions found.</div> : <div className="table-scroll"><table><thead><tr><th>Reference</th><th>Amount</th><th>Method</th><th>Date</th><th>Status</th></tr></thead><tbody>{rows.map((item) => <tr key={item.id}><td>{item.paystackReference || item.id}</td><td>{formatCurrency(item.amount)}</td><td>{item.paymentMethod || "—"}</td><td>{new Date(item.createdAt).toLocaleDateString()}</td><td><span className={`badge badge-${item.status}`}><i />{item.status}</span></td></tr>)}</tbody></table></div>}</div></>;
}
