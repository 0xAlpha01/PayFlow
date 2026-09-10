import { useEffect, useRef, useState } from "react";
import { ArrowDownLeft, ArrowUpRight, Wallet as WalletIcon } from "lucide-react";
import { getWalletBalance, getTransfers, sendTransfer } from "../services/walletService";
import { searchChatUsers } from "../services/chatService";
import { useAuth } from "../context/AuthContext";

const money = (value) => new Intl.NumberFormat("en-NG", { style: "currency", currency: "NGN" }).format(Number(value || 0));

export default function WalletPage() {
  const { user, refreshUser } = useAuth();
  const [balance, setBalance] = useState(null);
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [form, setForm] = useState({ recipientEmail: "", amount: "", note: "" });
  const [query, setQuery] = useState("");
  const [results, setResults] = useState([]);
  const [showResults, setShowResults] = useState(false);
  const [saving, setSaving] = useState(false);
  const searchTimer = useRef(null);
  const searchBox = useRef(null);

  async function load() {
    setLoading(true);
    setError("");
    try {
      const [balanceResponse, transfersResponse] = await Promise.all([
        getWalletBalance(),
        getTransfers({ page: 1, limit: 20 }),
      ]);
      setBalance(balanceResponse.data.data.balance);
      setRows(transfersResponse.data.data || []);
    } catch {
      setError("Unable to load your wallet.");
    } finally {
      setLoading(false);
    }
  }
  useEffect(() => {
    load();
  }, []);
  useEffect(() => {
    function onOutside(event) {
      if (searchBox.current && !searchBox.current.contains(event.target)) setShowResults(false);
    }
    document.addEventListener("mousedown", onOutside);
    return () => document.removeEventListener("mousedown", onOutside);
  }, []);

  function onQueryChange(value) {
    setQuery(value);
    setForm((current) => ({ ...current, recipientEmail: "" }));
    clearTimeout(searchTimer.current);
    if (!value.trim()) {
      setResults([]);
      return;
    }
    searchTimer.current = setTimeout(async () => {
      try {
        const { data } = await searchChatUsers(value.trim());
        setResults(data.data || []);
        setShowResults(true);
      } catch {
        setResults([]);
      }
    }, 300);
  }
  function selectRecipient(person) {
    setForm((current) => ({ ...current, recipientEmail: person.email }));
    setQuery(`${person.firstName} ${person.lastName} · ${person.email}`);
    setShowResults(false);
  }

  async function submit(event) {
    event.preventDefault();
    setError("");
    setNotice("");
    if (!form.recipientEmail) {
      setError("Search for a recipient and select them from the list, or enter their exact email.");
      return;
    }
    setSaving(true);
    try {
      const clientReference = crypto.randomUUID();
      await sendTransfer({ ...form, amount: Number(form.amount), clientReference });
      setNotice(`You sent ${money(form.amount)} successfully.`);
      setForm({ recipientEmail: "", amount: "", note: "" });
      setQuery("");
      await Promise.all([load(), refreshUser()]);
    } catch (requestError) {
      setError(requestError.response?.data?.message || "Unable to send this transfer.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <>
      <div className="page-heading">
        <div>
          <h2>Wallet</h2>
          <p>Send and receive money instantly between PayFlow accounts.</p>
        </div>
      </div>
      <div className="stats-grid compact">
        <div className="stat-card">
          <div className="stat-top">
            <span>Wallet balance</span>
            <span className="stat-icon blue">
              <WalletIcon size={18} />
            </span>
          </div>
          <strong>{balance === null ? "—" : money(balance)}</strong>
        </div>
      </div>
      {error && <div className="data-empty">{error}</div>}
      {notice && <div className="auth-alert success" style={{ marginBottom: 16 }}>{notice}</div>}
      <div className="panel" style={{ padding: 25 }}>
        <form onSubmit={submit}>
          <h3 style={{ marginTop: 0 }}>Send money</h3>
          <p style={{ color: "var(--muted)", fontSize: 12, marginTop: -8 }}>
            Search by name or email to find another PayFlow user.
          </p>
          <div className="form-grid">
            <div ref={searchBox} style={{ position: "relative" }}>
              <label>
                Recipient
                <input
                  placeholder="Search by name or email"
                  value={query}
                  onChange={(e) => onQueryChange(e.target.value)}
                  onFocus={() => results.length > 0 && setShowResults(true)}
                  autoComplete="off"
                  required
                />
              </label>
              {showResults && results.length > 0 && (
                <div className="user-search-results">
                  {results.map((person) => (
                    <button type="button" key={person.id} onClick={() => selectRecipient(person)}>
                      <strong>{person.firstName} {person.lastName}</strong>
                      <small>{person.email}</small>
                    </button>
                  ))}
                </div>
              )}
              {showResults && results.length === 0 && (
                <div className="user-search-results">
                  <div className="user-search-empty">No matching users.</div>
                </div>
              )}
            </div>
            <label>
              Amount (NGN)
              <input
                required
                type="number"
                min="0.01"
                step="0.01"
                value={form.amount}
                onChange={(e) => setForm({ ...form, amount: e.target.value })}
              />
            </label>
            <label>
              Note (optional)
              <input value={form.note} onChange={(e) => setForm({ ...form, note: e.target.value })} placeholder="What's this for?" />
            </label>
          </div>
          <div className="form-actions" style={{ padding: "20px 0 0" }}>
            <button className="button" disabled={saving}>{saving ? "Sending..." : "Send money"}</button>
          </div>
        </form>
      </div>
      <div className="panel table-panel" style={{ marginTop: 20 }}>
        <div className="table-tools">
          <strong>Transfer history</strong>
        </div>
        {loading ? (
          <div className="data-loading">Loading transfers...</div>
        ) : rows.length === 0 ? (
          <div className="data-empty">No transfers yet.</div>
        ) : (
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>Direction</th>
                  <th>Amount</th>
                  <th>Note</th>
                  <th>Date</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((item) => {
                  const sent = item.senderId === user?.id;
                  return (
                    <tr key={item.id}>
                      <td>
                        {sent ? <ArrowUpRight size={15} color="#c0392b" /> : <ArrowDownLeft size={15} color="#16a66a" />}{" "}
                        {sent ? "Sent" : "Received"}
                      </td>
                      <td>{money(item.amount)}</td>
                      <td>{item.note || "—"}</td>
                      <td>{new Date(item.createdAt).toLocaleString()}</td>
                      <td>
                        <span className={`badge badge-${item.status}`}><i />{item.status}</span>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </>
  );
}
