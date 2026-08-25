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
  ShieldCheck,
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

const transactions = [
  {
    name: "Olivia Rhye",
    email: "olivia@acme.io",
    initials: "OR",
    id: "#TRX-8291",
    amount: "$2,400.00",
    method: "Visa •••• 4242",
    date: "Mar 24, 2024",
    status: "Successful",
  },
  {
    name: "Phoenix Baker",
    email: "phoenix@lattice.com",
    initials: "PB",
    id: "#TRX-8290",
    amount: "$1,820.00",
    method: "Mastercard •••• 8831",
    date: "Mar 24, 2024",
    status: "Pending",
  },
  {
    name: "Lana Steiner",
    email: "lana@catalog.com",
    initials: "LS",
    id: "#TRX-8289",
    amount: "$890.00",
    method: "Visa •••• 1092",
    date: "Mar 23, 2024",
    status: "Successful",
  },
  {
    name: "Demi Wilkinson",
    email: "demi@circooles.com",
    initials: "DW",
    id: "#TRX-8288",
    amount: "$3,240.00",
    method: "Apple Pay",
    date: "Mar 23, 2024",
    status: "Refunded",
  },
  {
    name: "Candice Wu",
    email: "candice@koru.com",
    initials: "CW",
    id: "#TRX-8287",
    amount: "$450.00",
    method: "Visa •••• 2244",
    date: "Mar 22, 2024",
    status: "Failed",
  },
];
const customers = transactions.map((item, index) => ({
  ...item,
  spent: ["$12,450.00", "$8,920.00", "$6,240.00", "$4,810.00", "$3,120.00"][
    index
  ],
  payments: [24, 18, 14, 11, 8][index],
}));
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
        <div className="workspace-avatar">A</div>
        <div>
          <strong>Acme Inc.</strong>
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
        <div className="top-avatar">JD</div>
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
function Table({ rows = transactions, customerTable = false }) {
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
  // The initial request synchronizes this view with the protected API.
  // eslint-disable-next-line react-hooks/set-state-in-effect
  useEffect(() => {
    loadDashboard();
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
function ListPage({ type }) {
  if (type === "customers")
    return (
      <Page
        title="Customers"
        description="Manage your customers and their payment activity."
        action="Add customer"
      >
        <div className="stats-grid compact">
          <StatCard
            title="Total customers"
            value="2,420"
            change="18.2%"
            icon={Users}
          />
          <StatCard
            title="Active customers"
            value="1,892"
            change="10.4%"
            icon={Activity}
            tone="blue"
          />
          <StatCard
            title="Average spend"
            value="$1,240"
            change="6.8%"
            icon={CircleDollarSign}
            tone="orange"
          />
        </div>
        <Table rows={customers} customerTable />
      </Page>
    );
  if (type === "transactions")
    return (
      <Page
        title="Transactions"
        description="Review and manage every payment across your business."
        action="Export"
      >
        <Table />
      </Page>
    );
  if (type === "invoices")
    return (
      <Page
        title="Invoices"
        description="Create, send, and track invoices in one place."
        action="Create invoice"
      >
        <div className="empty-panel">
          <FileText size={28} />
          <h3>Your invoices are ready when you are</h3>
          <p>
            Create your first invoice to start collecting payments from
            customers.
          </p>
          <Button>
            <Plus size={16} /> Create invoice
          </Button>
        </div>
      </Page>
    );
  if (type === "payments")
    return (
      <Page
        title="Payments"
        description="Monitor payment performance and activity."
        action="Create payment"
      >
        <div className="stats-grid compact">
          <StatCard
            title="Successful payments"
            value="$214,820"
            change="12.5%"
            icon={ArrowUpRight}
          />
          <StatCard
            title="Pending payments"
            value="$12,840"
            change="4.1%"
            icon={Activity}
            tone="orange"
          />
          <StatCard
            title="Failed payments"
            value="$3,240"
            change="-2.3%"
            icon={ArrowDownLeft}
            tone="purple"
          />
        </div>
        <Table />
      </Page>
    );
  if (type === "subscriptions")
    return (
      <Page
        title="Subscriptions"
        description="Keep your recurring revenue growing."
        action="Create plan"
      >
        <div className="stats-grid compact">
          <StatCard
            title="Monthly recurring revenue"
            value="$84,240"
            change="14.2%"
            icon={CircleDollarSign}
          />
          <StatCard
            title="Active subscriptions"
            value="1,248"
            change="8.4%"
            icon={CreditCard}
            tone="blue"
          />
          <StatCard
            title="Upcoming renewals"
            value="86"
            change="5.1%"
            icon={Activity}
            tone="orange"
          />
        </div>
        <Table />
      </Page>
    );
  return (
    <Page
      title="Analytics"
      description="Understand your business performance at a glance."
    >
      <ChartCard />
      <div className="analytics-grid">
        <div className="metric-panel">
          <h3>Payment success rate</h3>
          <strong>96.8%</strong>
          <span>↑ 2.4% this month</span>
          <div className="progress">
            <i style={{ width: "96.8%" }} />
          </div>
        </div>
        <div className="metric-panel">
          <h3>Customer growth</h3>
          <strong>+284</strong>
          <span>↑ 18.2% this month</span>
          <div className="mini-bars">
            {[30, 50, 38, 62, 48, 82, 70, 92].map((height, i) => (
              <i key={i} style={{ height: `${height}%` }} />
            ))}
          </div>
        </div>
      </div>
    </Page>
  );
}
function Settings() {
  const [saved, setSaved] = useState(false);
  return (
    <Page
      title="Settings"
      description="Manage your account and business preferences."
    >
      <div className="settings-layout">
        <div className="settings-nav">
          {[
            "Profile",
            "Business",
            "Payment settings",
            "Security",
            "Notifications",
          ].map((item, i) => (
            <button className={i === 0 ? "selected" : ""} key={item}>
              {item}
            </button>
          ))}
        </div>
        <div className="settings-content">
          <div className="settings-section">
            <div>
              <h3>Profile information</h3>
              <p>Update your personal details and contact information.</p>
            </div>
            <div className="form-grid">
              <label>
                Full name
                <input defaultValue="Jordan Davis" />
              </label>
              <label>
                Email address
                <input defaultValue="jordan@acme.io" />
              </label>
              <label>
                Phone number
                <input defaultValue="+1 (555) 012-3456" />
              </label>
              <label>
                Role
                <input defaultValue="Administrator" />
              </label>
            </div>
          </div>
          <div className="settings-section security-row">
            <div className="setting-icon">
              <ShieldCheck size={21} />
            </div>
            <div>
              <h3>Security preferences</h3>
              <p>
                Two-factor authentication and active sessions are managed here.
              </p>
            </div>
            <Button secondary>Manage</Button>
          </div>
          <div className="form-actions">
            <button className="button button-secondary">Cancel</button>
            <Button onClick={() => setSaved(true)}>
              {saved ? "Saved" : "Save changes"}
            </Button>
          </div>
        </div>
      </div>
    </Page>
  );
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
                element={<ListPage type="customers" />}
              />
              <Route
                path="/customers"
                element={<ListPage type="customers" />}
              />
              <Route
                path="/invoices/create"
                element={
                  <Page
                    title="Create invoice"
                    description="Send a professional invoice to your customer."
                  >
                    <div className="empty-panel form-placeholder">
                      <FileText size={28} />
                      <h3>Invoice builder</h3>
                      <p>
                        Add customers, line items, taxes, and due dates here.
                      </p>
                      <Button>Start building</Button>
                    </div>
                  </Page>
                }
              />
              <Route
                path="/invoices/:id"
                element={<ListPage type="invoices" />}
              />
              <Route path="/invoices" element={<ListPage type="invoices" />} />
              <Route
                path="/transactions"
                element={<ListPage type="transactions" />}
              />
              <Route path="/payments" element={<ListPage type="payments" />} />
              <Route
                path="/subscriptions"
                element={<ListPage type="subscriptions" />}
              />
              <Route
                path="/analytics"
                element={<ListPage type="analytics" />}
              />
              <Route path="/settings/*" element={<Settings />} />
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
