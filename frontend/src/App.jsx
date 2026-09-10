import { useEffect, useState } from "react";
import {
  BrowserRouter,
  Navigate,
  NavLink,
  Route,
  Routes,
  useLocation,
} from "react-router-dom";
import {
  Area,
  AreaChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import {
  Activity,
  ArrowDownLeft,
  ArrowUpRight,
  Bell,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  CircleDollarSign,
  CreditCard,
  FileText,
  LayoutDashboard,
  Menu,
  Moon,
  MoreHorizontal,
  Plus,
  Search,
  Settings as SettingsIcon,
  Sun,
  Users,
  Wallet,
  X,
} from "lucide-react";
import "./App.css";
import {
  Login,
  Signup,
  VerifyEmail,
  VerifyPhone,
  ForgotPassword,
  ResetPassword,
} from "./pages/Auth";
import ProtectedRoute from "./routes/ProtectedRoute";
import PublicRoute from "./routes/PublicRoute";
import { useAuth } from "./context/AuthContext";
import { dashboardService } from "./services/dashboardService";
import CustomersPage from "./pages/Customers";
import PaymentsPage from "./pages/Payments";
import TransactionsPage from "./pages/Transactions";
import InvoicesPage from "./pages/Invoices";
import PaymentCallback from "./pages/PaymentCallback";
import SettingsPage from "./pages/Settings";
import { getCustomers } from "./services/customerService";
import { cancelSubscription, createSubscription, getSubscriptionSummary, getSubscriptions } from "./services/subscriptionService";
import { getCustomerAnalytics, getPaymentAnalytics, getRevenueAnalytics } from "./services/analyticsService";

const navGroups = [
  {
    label: "Overview",
    items: [{ label: "Dashboard", icon: LayoutDashboard, to: "/dashboard" }],
  },
  {
    label: "Manage",
    items: [
      { label: "Payments", icon: Wallet, to: "/payments" },
      { label: "Transactions", icon: Activity, to: "/transactions" },
      { label: "Customers", icon: Users, to: "/customers" },
      { label: "Invoices", icon: FileText, to: "/invoices" },
      { label: "Subscriptions", icon: CreditCard, to: "/subscriptions" },
    ],
  },
  {
    label: "Insights",
    items: [{ label: "Analytics", icon: CircleDollarSign, to: "/analytics" }],
  },
];

function Badge({ children }) {
  return (
    <span className={`badge badge-${children.toLowerCase()}`}>
      <i />
      {children}
    </span>
  );
}
function Button({ children, secondary = false, onClick }) {
  return (
    <button
      className={`button ${secondary ? "button-secondary" : ""}`}
      onClick={onClick}
    >
      {children}
    </button>
  );
}
function StatCard({ title, value, change, icon: Icon, tone = "green" }) {
  return (
    <div className="stat-card">
      <div className="stat-top">
        <span>{title}</span>
        <span className={`stat-icon ${tone}`}>
          <Icon size={18} />
        </span>
      </div>
      <strong>{value}</strong>
      <p className={change.startsWith("-") ? "negative" : ""}>
        <span>
          {change.startsWith("-") ? (
            <ArrowDownLeft size={14} />
          ) : (
            <ArrowUpRight size={14} />
          )}
        </span>
        {change} <em>vs last month</em>
      </p>
    </div>
  );
}
function Sidebar({ open, setOpen, dark, setDark }) {
  const { user, logout } = useAuth();
  const [loggingOut, setLoggingOut] = useState(false);
  async function handleLogout() {
    if (loggingOut) return;
    setLoggingOut(true);
    try {
      await logout();
    } catch {
      /* local state is cleared by AuthContext */
    }
  }
  return (
    <aside className={`sidebar ${open ? "open" : ""}`}>
      <div className="brand">
        <div className="brand-mark">P</div>
        <span>
          pay<span>flow</span>
        </span>
        <button className="mobile-close" onClick={() => setOpen(false)}>
          <X size={18} />
        </button>
      </div>
      <div className="workspace">
        <div className="workspace-avatar">{(user?.businessName || user?.firstName || "P")[0].toUpperCase()}</div>
        <div>
          <strong>{user?.businessName || "Your business"}</strong>
          <small>Business account</small>
        </div>
        <ChevronDown size={15} />
      </div>
      <nav>
        {navGroups.map((group) => (
          <div className="nav-group" key={group.label}>
            <small>{group.label}</small>
            {group.items.map(({ label, icon: Icon, to }) => (
              <NavLink
                key={label}
                to={to}
                onClick={() => setOpen(false)}
                className={({ isActive }) => (isActive ? "active" : "")}
              >
                <Icon size={18} />
                {label}
                {label === "Payments" && <span className="nav-count">8</span>}
              </NavLink>
            ))}
          </div>
        ))}
      </nav>
      <div className="sidebar-bottom">
        <NavLink to="/settings/profile">
          <SettingsIcon size={18} />
          Settings
        </NavLink>
        <button className="theme-toggle" onClick={() => setDark(!dark)}>
          {dark ? <Sun size={18} /> : <Moon size={18} />}{" "}
          {dark ? "Light mode" : "Dark mode"}
        </button>
        <button
          className="user"
          onClick={handleLogout}
          disabled={loggingOut}
          title="Sign out"
        >
          <div className="avatar">
            {user?.firstName?.[0] || "J"}
            {user?.lastName?.[0] || "D"}
          </div>
          <div>
            <strong>
              {loggingOut
                ? "Logging out..."
                : `${user?.firstName || "Jordan"} ${user?.lastName || "Davis"}`}
            </strong>
            <small>{user?.email || "jordan@acme.io"}</small>
          </div>
          <MoreHorizontal size={18} />
        </button>
      </div>
    </aside>
  );
}
function Topbar({ setOpen }) {
  const location = useLocation();
  const { user } = useAuth();
  const title =
    location.pathname === "/dashboard" || location.pathname === "/"
      ? "Dashboard"
      : location.pathname.split("/")[1].replace("-", " ");
  return (
    <header className="topbar">
      <button className="mobile-menu" onClick={() => setOpen(true)}>
        <Menu size={22} />
      </button>
      <div>
        <h1>{title.charAt(0).toUpperCase() + title.slice(1)}</h1>
        <p>Here’s what’s happening with your business today.</p>
      </div>
      <div className="top-actions">
        <button className="icon-button">
          <Search size={19} />
        </button>
        <button className="icon-button notification">
          <Bell size={19} />
          <i />
        </button>
        <div className="top-avatar">{initials(`${user?.firstName || ""} ${user?.lastName || ""}`.trim()) || "PF"}</div>
      </div>
    </header>
  );
}
function ChartCard({ data = [], loading = false, onPeriodChange }) {
  const [period, setPeriod] = useState("30 Days");
  const chartData = data.map((item) => ({
    name: item.date,
    value: item.revenue,
  }));
  const latestRevenue = chartData.at(-1)?.value || 0;
  return (
    <div className="panel chart-panel">
      <div className="panel-heading">
        <div>
          <h2>Revenue overview</h2>
          <p>Track your revenue performance over time.</p>
        </div>
        <select
          value={period}
          onChange={(e) => {
            const nextPeriod = e.target.value;
            setPeriod(nextPeriod);
            onPeriodChange?.(
              {
                "7 Days": "7d",
                "30 Days": "30d",
                "3 Months": "3m",
                "6 Months": "6m",
                "1 Year": "1y",
              }[nextPeriod],
            );
          }}
        >
          {["7 Days", "30 Days", "3 Months", "6 Months", "1 Year"].map(
            (item) => (
              <option key={item}>{item}</option>
            ),
          )}
        </select>
      </div>
      <div className="chart-legend">
        <strong>{formatMoney(latestRevenue)}</strong>
        <span>
          <ArrowUpRight size={14} /> Live
        </span>
        <small>vs last month</small>
      </div>
      <div className="chart">
        {loading ? (
          <div className="data-loading">Loading revenue...</div>
        ) : chartData.length === 0 ? (
          <div className="data-empty">No revenue data available.</div>
        ) : (
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData}>
              <defs>
                <linearGradient id="revenueFill" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor="#16a66a" stopOpacity={0.28} />
                  <stop offset="100%" stopColor="#16a66a" stopOpacity={0} />
                </linearGradient>
              </defs>
              <XAxis
                dataKey="name"
                axisLine={false}
                tickLine={false}
                tick={{ fill: "#8b96a5", fontSize: 11 }}
                interval={1}
              />
              <YAxis
                axisLine={false}
                tickLine={false}
                tick={{ fill: "#8b96a5", fontSize: 11 }}
                tickFormatter={(value) => `$${value / 1000}k`}
                width={44}
              />
              <Tooltip
                formatter={(value) => [`$${value.toLocaleString()}`, "Revenue"]}
                contentStyle={{ border: "1px solid #e6ebe8", borderRadius: 8 }}
              />
              <Area
                type="monotone"
                dataKey="value"
                stroke="#16a66a"
                strokeWidth={2.5}
                fill="url(#revenueFill)"
              />
            </AreaChart>
          </ResponsiveContainer>
        )}
      </div>
    </div>
  );
}
function formatMoney(value) {
  return `$${Number(value || 0).toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}
function initials(name = "") {
  return name
    .split(" ")
    .map((part) => part[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();
}
function dashboardTransactions(rows = []) {
  return rows.map((row) => ({
    name: row.customerName,
    email: row.currency,
    initials: initials(row.customerName),
    id: `#${row.id?.slice(-8) || "TRANSACTION"}`,
    amount: formatMoney(row.amount),
    method: row.paymentMethod,
    date: new Date(row.createdAt).toLocaleDateString(),
    status: row.status.charAt(0).toUpperCase() + row.status.slice(1),
  }));
}
function dashboardCustomers(rows = []) {
  return rows.map((row) => ({
    name: row.name,
    email: row.email,
    initials: initials(row.name),
    id: row.id,
    spent: formatMoney(row.amountSpent),
    payments: 0,
    date: row.lastPayment
      ? new Date(row.lastPayment).toLocaleDateString()
      : "No payments",
    status: row.status.charAt(0).toUpperCase() + row.status.slice(1),
  }));
}
function Table({ rows = [], customerTable = false }) {
  const [query, setQuery] = useState("");
  const filtered = rows.filter((row) =>
    `${row.name} ${row.email} ${row.id}`
      .toLowerCase()
      .includes(query.toLowerCase()),
  );
  return (
    <div className="panel table-panel">
      <div className="table-tools">
        <div className="search">
          <Search size={17} />
          <input
            placeholder="Search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        <Button secondary>
          Filter <ChevronDown size={15} />
        </Button>
      </div>
      <div className="table-scroll">
        <table>
          <thead>
            <tr>
              <th>{customerTable ? "Customer" : "Transaction"}</th>
              <th>{customerTable ? "Total spent" : "Amount"}</th>
              <th>{customerTable ? "Payments" : "Payment method"}</th>
              <th>{customerTable ? "Last payment" : "Date"}</th>
              <th>Status</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {filtered.map((row) => (
              <tr key={row.id}>
                <td>
                  <div className="person">
                    <div className="person-avatar">{row.initials}</div>
                    <div>
                      <strong>{row.name}</strong>
                      <small>{customerTable ? row.email : row.id}</small>
                    </div>
                  </div>
                </td>
                <td>
                  <strong>{customerTable ? row.spent : row.amount}</strong>
                </td>
                <td>
                  {customerTable ? `${row.payments} payments` : row.method}
                </td>
                <td>{customerTable ? row.date : row.date}</td>
                <td>
                  <Badge>{row.status}</Badge>
                </td>
                <td>
                  <button className="more">
                    <MoreHorizontal size={17} />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="table-footer">
        <span>
          Showing 1–{filtered.length} of {rows.length}
        </span>
        <div>
          <button>
            <ChevronLeft size={16} />
          </button>
          <button className="current">1</button>
          <button>2</button>
          <button>3</button>
          <button>
            <ChevronRight size={16} />
          </button>
        </div>
      </div>
    </div>
  );
}
function Dashboard() {
  const { user } = useAuth();
  const [data, setData] = useState({
    summary: {},
    revenue: [],
    transactions: [],
    customers: [],
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  async function loadDashboard(period = "30d") {
    setLoading(true);
    setError(false);
    try {
      const [summary, revenueData, transactionData, customerData] =
        await Promise.all([
          dashboardService.getDashboardSummary(),
          dashboardService.getRevenue(period),
          dashboardService.getRecentTransactions(),
          dashboardService.getRecentCustomers(),
        ]);
      setData({
        summary: summary.data.data,
        revenue: revenueData.data.data || [],
        transactions: transactionData.data.data || [],
        customers: customerData.data.data || [],
      });
    } catch {
      setError(true);
    } finally {
      setLoading(false);
    }
  }
  useEffect(() => {
    void Promise.resolve().then(() => loadDashboard());
  }, []);
  const summary = data.summary;
  if (error)
    return (
      <div className="dashboard-error">
        <h2>Unable to load dashboard data.</h2>
        <button className="button" onClick={loadDashboard}>
          Try again
        </button>
      </div>
    );
  return (
    <>
      <div className="welcome-row">
        <div>
          <h2>
            Good morning, {user?.firstName || "there"} <span>✦</span>
          </h2>
          <p>Keep an eye on your numbers. You're doing great.</p>
        </div>
        <Button>
          <Plus size={17} /> Create invoice
        </Button>
      </div>
      {loading ? (
        <div className="dashboard-loading">Loading your financial data...</div>
      ) : (
        <div className="stats-grid">
          <StatCard
            title="Total revenue"
            value={formatMoney(summary.totalRevenue)}
            change="Live"
            icon={CircleDollarSign}
          />
          <StatCard
            title="Available balance"
            value={formatMoney(summary.availableBalance)}
            change="Live"
            icon={Wallet}
            tone="blue"
          />
          <StatCard
            title="Pending payments"
            value={formatMoney(summary.pendingPayments)}
            change="Live"
            icon={Activity}
            tone="orange"
          />
          <StatCard
            title="Total customers"
            value={Number(summary.totalCustomers || 0).toLocaleString()}
            change="Live"
            icon={Users}
            tone="purple"
          />
        </div>
      )}
      <ChartCard
        data={data.revenue}
        loading={loading}
        onPeriodChange={loadDashboard}
      />
      <div className="section-heading">
        <div>
          <h2>Recent transactions</h2>
          <p>Keep track of your most recent payments.</p>
        </div>
        <NavLink to="/transactions">
          View all <ArrowUpRight size={15} />
        </NavLink>
      </div>
      {loading ? (
        <div className="panel data-loading">Loading transactions...</div>
      ) : data.transactions.length === 0 ? (
        <div className="panel data-empty">No transactions yet.</div>
      ) : (
        <Table rows={dashboardTransactions(data.transactions)} />
      )}
      <div className="section-heading dashboard-subsection">
        <div>
          <h2>Recent customers</h2>
          <p>Your latest customer activity.</p>
        </div>
        <NavLink to="/customers">
          View all <ArrowUpRight size={15} />
        </NavLink>
      </div>
      {loading ? (
        <div className="panel data-loading">Loading customers...</div>
      ) : data.customers.length === 0 ? (
        <div className="panel data-empty">No customers yet.</div>
      ) : (
        <Table rows={dashboardCustomers(data.customers)} customerTable />
      )}
    </>
  );
}
function Page({ title, description, action, children }) {
  return (
    <>
      <div className="page-heading">
        <div>
          <h2>{title}</h2>
          <p>{description}</p>
        </div>
        {action && (
          <Button>
            <Plus size={17} /> {action}
          </Button>
        )}
      </div>
      {children}
    </>
  );
}
function SubscriptionsPanel() {
  const [rows, setRows] = useState([]), [summary, setSummary] = useState(null), [loading, setLoading] = useState(true), [error, setError] = useState(""), [search, setSearch] = useState(""), [status, setStatus] = useState(""), [page, setPage] = useState(1), [meta, setMeta] = useState({}), [form, setForm] = useState(null), [customers, setCustomers] = useState([]);
  const load = async () => { setLoading(true); setError(""); try { const [list, totals] = await Promise.all([getSubscriptions({ page, limit: 20, search, status }), getSubscriptionSummary()]); setRows(Array.isArray(list.data.data) ? list.data.data : []); setMeta(list.data.pagination || {}); setSummary(totals.data.data || null); } catch { setError("Unable to load subscriptions."); } finally { setLoading(false); } };
  useEffect(() => { const timer = setTimeout(load, search ? 300 : 0); return () => clearTimeout(timer); }, [page, search, status]);
  const create = async () => { try { const response = await getCustomers({ page: 1, limit: 100, status: "active" }); setCustomers(Array.isArray(response.data.data) ? response.data.data : []); setForm({ customerId: "", plan: "", amount: "", billingCycle: "monthly", status: "active", nextPaymentDate: "", currency: "NGN" }); } catch { setError("Unable to load customers."); } };
  const save = async (event) => { event.preventDefault(); try { await createSubscription({ ...form, amount: Number(form.amount) }); setForm(null); load(); } catch { setError("Unable to create subscription."); } };
  return <><div className="page-heading"><div><h2>Subscriptions</h2><p>Keep your recurring revenue growing.</p></div><Button onClick={create}><Plus size={17}/> Create plan</Button></div><div className="stats-grid compact"><StatCard title="Monthly recurring revenue" value={summary ? formatMoney(summary.monthlyRecurringRevenue) : "—"} change="" icon={CircleDollarSign}/><StatCard title="Active subscriptions" value={summary ? String(summary.activeSubscriptions) : "—"} change="" icon={CreditCard} tone="blue"/><StatCard title="Upcoming renewals" value={summary ? String(summary.upcomingRenewals) : "—"} change="" icon={Activity} tone="orange"/></div>{error && <div className="data-empty">{error} <button onClick={load}>Try again</button></div>}<div className="panel table-panel"><div className="table-tools"><div className="search"><Search size={17}/><input placeholder="Search plans" value={search} onChange={(e)=>{setSearch(e.target.value);setPage(1)}}/></div><select value={status} onChange={(e)=>{setStatus(e.target.value);setPage(1)}}><option value="">All statuses</option><option value="active">Active</option><option value="paused">Paused</option><option value="cancelled">Cancelled</option><option value="expired">Expired</option></select></div>{loading?<div className="data-loading">Loading subscriptions...</div>:rows.length===0?<div className="data-empty">No subscriptions found.</div>:<div className="table-scroll"><table><thead><tr><th>Plan</th><th>Amount</th><th>Billing cycle</th><th>Next payment</th><th>Status</th><th/></tr></thead><tbody>{rows.map((item)=><tr key={item.id}><td>{item.plan}</td><td>{formatMoney(item.amount)}</td><td>{item.billingCycle}</td><td>{new Date(item.nextPaymentDate).toLocaleDateString()}</td><td><Badge>{item.status}</Badge></td><td>{item.status!=="cancelled"&&<button className="more" onClick={async()=>{try{await cancelSubscription(item.id);load()}catch{setError("Unable to cancel subscription.")}}}>Cancel</button>}</td></tr>)}</tbody></table></div>}<div className="table-footer"><span>Showing {rows.length} of {meta.total || 0}</span><div><button disabled={page<=1} onClick={()=>setPage(page-1)}><ChevronLeft size={16}/></button><button className="current">{page}</button><button disabled={page >= (meta.totalPages || 0)} onClick={()=>setPage(page+1)}><ChevronRight size={16}/></button></div></div></div>{form&&<div className="panel" style={{marginTop:20}}><form onSubmit={save}><h3>Create plan</h3><div className="table-tools"><select required value={form.customerId} onChange={(e)=>setForm({...form,customerId:e.target.value})}><option value="">Customer</option>{customers.map((customer)=><option key={customer.id} value={customer.id}>{customer.firstName} {customer.lastName}</option>)}</select><input required placeholder="Plan" value={form.plan} onChange={(e)=>setForm({...form,plan:e.target.value})}/><input required type="number" min="0.01" placeholder="Amount" value={form.amount} onChange={(e)=>setForm({...form,amount:e.target.value})}/><select value={form.billingCycle} onChange={(e)=>setForm({...form,billingCycle:e.target.value})}><option value="monthly">Monthly</option><option value="quarterly">Quarterly</option><option value="yearly">Yearly</option></select><input required type="date" value={form.nextPaymentDate} onChange={(e)=>setForm({...form,nextPaymentDate:e.target.value})}/><Button>Save plan</Button><button type="button" className="button button-secondary" onClick={()=>setForm(null)}>Cancel</button></div></form></div>}</>;
}
function AnalyticsPanel() {
  const [period, setPeriod] = useState("30d"), [revenue, setRevenue] = useState([]), [payments, setPayments] = useState(null), [customers, setCustomers] = useState(null), [loading, setLoading] = useState(true), [error, setError] = useState("");
  const load = async () => { setLoading(true); setError(""); try { const [r, p, c] = await Promise.all([getRevenueAnalytics(period), getPaymentAnalytics(period), getCustomerAnalytics(period)]); setRevenue(Array.isArray(r.data.data) ? r.data.data : []); setPayments(p.data.data || null); setCustomers(c.data.data || null); } catch { setError("Unable to load analytics."); } finally { setLoading(false); } };
  useEffect(() => { load(); }, [period]);
  const rate = Number(payments?.paymentSuccessRate || 0), growth = Number(customers?.customerGrowth || 0);
  return <Page title="Analytics" description="Understand your business performance at a glance.">{error && <div className="data-empty">{error} <button onClick={load}>Try again</button></div>}<ChartCard data={revenue} loading={loading} onPeriodChange={setPeriod}/><div className="analytics-grid"><div className="metric-panel"><h3>Payment success rate</h3><strong>{loading ? "—" : `${rate.toFixed(1)}%`}</strong><span>{loading ? "Loading analytics..." : payments?.total ? `${payments.successful} successful payments` : "No payment data available."}</span><div className="progress"><i style={{width:`${Math.min(rate,100)}%`}}/></div></div><div className="metric-panel"><h3>Customer growth</h3><strong>{loading ? "—" : `+${growth}`}</strong><span>{loading ? "Loading analytics..." : customers?.newCustomers ? `${customers.newCustomers} new customers` : "No customer data available."}</span><div className="mini-bars">{Array.from({length:8},(_,i)=><i key={i} style={{height:`${customers?.newCustomers ? (i+1)*12.5 : 0}%`}}/>)}</div></div></div></Page>;
}
function AppShell() {
  const [open, setOpen] = useState(false);
  const [dark, setDark] = useState(false);
  const location = useLocation();
  const isAuthPage = [
    "/login",
    "/signup",
    "/verify-email",
    "/verify-phone",
    "/forgot-password",
    "/reset-password",
  ].includes(location.pathname);
  if (isAuthPage) {
    return (
      <Routes>
        <Route element={<PublicRoute />}>
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<Signup />} />
          <Route path="/verify-email" element={<VerifyEmail />} />
          <Route path="/verify-phone" element={<VerifyPhone />} />
          <Route path="/forgot-password" element={<ForgotPassword />} />
          <Route path="/reset-password" element={<ResetPassword />} />
        </Route>
      </Routes>
    );
  }
  return (
    <div className={dark ? "app dark" : "app"}>
      <Sidebar open={open} setOpen={setOpen} dark={dark} setDark={setDark} />
      <main className="main">
        <Topbar setOpen={setOpen} />
        <div className="content">
          <Routes>
            <Route element={<PublicRoute />}>
              <Route path="/login" element={<Login />} />
              <Route path="/signup" element={<Signup />} />
              <Route path="/verify-email" element={<VerifyEmail />} />
              <Route path="/verify-phone" element={<VerifyPhone />} />
              <Route path="/forgot-password" element={<ForgotPassword />} />
              <Route path="/reset-password" element={<ResetPassword />} />
            </Route>
            <Route element={<ProtectedRoute />}>
              <Route path="/" element={<Dashboard />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route
                path="/customers/:id"
                element={<CustomersPage />}
              />
              <Route
                path="/customers"
                element={<CustomersPage />}
              />
              <Route path="/invoices/create" element={<InvoicesPage />} />
              <Route
                path="/invoices/:id"
                element={<InvoicesPage />}
              />
              <Route path="/invoices" element={<InvoicesPage />} />
              <Route
                path="/transactions"
                element={<TransactionsPage />}
              />
              <Route path="/payments" element={<PaymentsPage />} />
              <Route path="/payments/callback" element={<PaymentCallback />} />
              <Route
                path="/subscriptions"
                element={<SubscriptionsPanel />}
              />
              <Route
                path="/analytics"
                element={<AnalyticsPanel />}
              />
              <Route path="/settings/*" element={<SettingsPage />} />
              <Route
                path="*"
                element={
                  <Page
                    title="Page not found"
                    description="The page you requested does not exist."
                  >
                    <Button onClick={() => window.history.back()}>
                      <ChevronLeft size={16} /> Go back
                    </Button>
                  </Page>
                }
              />
            </Route>
            <Route path="*" element={<Navigate to="/login" replace />} />
          </Routes>
        </div>
      </main>
    </div>
  );
}
export default function App() {
  return (
    <BrowserRouter>
      <AppShell />
    </BrowserRouter>
  );
}
