import { useEffect, useState } from "react";
import { ShieldCheck } from "lucide-react";
import {
  changePassword,
  getBusiness,
  getNotifications,
  getProfile,
  updateBusiness,
  updateNotifications,
  updateProfile,
} from "../services/settingsService";
import { useAuth } from "../context/AuthContext";
import PasswordInput from "../components/auth/PasswordInput";
import PasswordStrength from "../components/auth/PasswordStrength";

const TABS = ["Profile", "Business", "Notifications", "Security"];

function Notice({ message, tone = "success" }) {
  if (!message) return null;
  return <div className={`auth-alert ${tone}`} style={{ margin: "0 25px 16px" }}>{message}</div>;
}

function ProfileTab() {
  const [form, setForm] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  useEffect(() => {
    getProfile()
      .then(({ data }) => setForm(data.data))
      .catch(() => setError("Unable to load your profile."))
      .finally(() => setLoading(false));
  }, []);
  async function save(event) {
    event.preventDefault();
    setSaving(true);
    setError("");
    setNotice("");
    try {
      await updateProfile({ firstName: form.firstName, lastName: form.lastName, phone: form.phone });
      setNotice("Profile updated successfully.");
    } catch (requestError) {
      setError(requestError.response?.data?.message || "Unable to update your profile.");
    } finally {
      setSaving(false);
    }
  }
  if (loading) return <div className="settings-section"><div className="data-loading">Loading profile...</div></div>;
  if (!form) return <div className="settings-section"><div className="data-empty">{error || "Unable to load your profile."}</div></div>;
  return (
    <form onSubmit={save}>
      <div className="settings-section">
        <div>
          <h3>Profile information</h3>
          <p>Update your personal details and contact information.</p>
        </div>
        <Notice message={notice} />
        <Notice message={error} tone="error" />
        <div className="form-grid">
          <label>
            First name
            <input value={form.firstName || ""} onChange={(e) => setForm({ ...form, firstName: e.target.value })} required />
          </label>
          <label>
            Last name
            <input value={form.lastName || ""} onChange={(e) => setForm({ ...form, lastName: e.target.value })} required />
          </label>
          <label>
            Email address
            <input value={form.email || ""} disabled title="Email address cannot be changed here." />
          </label>
          <label>
            Phone number
            <input value={form.phone || ""} onChange={(e) => setForm({ ...form, phone: e.target.value })} placeholder="+2348012345678" />
          </label>
        </div>
      </div>
      <div className="form-actions">
        <button className="button" disabled={saving}>{saving ? "Saving..." : "Save changes"}</button>
      </div>
    </form>
  );
}

function BusinessTab() {
  const [form, setForm] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  useEffect(() => {
    getBusiness()
      .then(({ data }) => setForm(data.data))
      .catch(() => setError("Unable to load business settings."))
      .finally(() => setLoading(false));
  }, []);
  async function save(event) {
    event.preventDefault();
    setSaving(true);
    setError("");
    setNotice("");
    try {
      await updateBusiness(form);
      setNotice("Business settings updated successfully.");
    } catch (requestError) {
      setError(requestError.response?.data?.message || "Unable to update business settings.");
    } finally {
      setSaving(false);
    }
  }
  if (loading) return <div className="settings-section"><div className="data-loading">Loading business settings...</div></div>;
  if (!form) return <div className="settings-section"><div className="data-empty">{error || "Unable to load business settings."}</div></div>;
  return (
    <form onSubmit={save}>
      <div className="settings-section">
        <div>
          <h3>Business details</h3>
          <p>These details appear on invoices and receipts sent to your customers.</p>
        </div>
        <Notice message={notice} />
        <Notice message={error} tone="error" />
        <div className="form-grid">
          <label>
            Business name
            <input value={form.businessName || ""} onChange={(e) => setForm({ ...form, businessName: e.target.value })} />
          </label>
          <label>
            Business email
            <input type="email" value={form.businessEmail || ""} onChange={(e) => setForm({ ...form, businessEmail: e.target.value })} />
          </label>
          <label>
            Business phone
            <input value={form.businessPhone || ""} onChange={(e) => setForm({ ...form, businessPhone: e.target.value })} />
          </label>
          <label>
            Website
            <input value={form.businessWebsite || ""} onChange={(e) => setForm({ ...form, businessWebsite: e.target.value })} placeholder="https://" />
          </label>
          <label>
            Address
            <input value={form.businessAddress || ""} onChange={(e) => setForm({ ...form, businessAddress: e.target.value })} />
          </label>
          <label>
            Description
            <input value={form.businessDescription || ""} onChange={(e) => setForm({ ...form, businessDescription: e.target.value })} />
          </label>
        </div>
      </div>
      <div className="form-actions">
        <button className="button" disabled={saving}>{saving ? "Saving..." : "Save changes"}</button>
      </div>
    </form>
  );
}

const NOTIFICATION_ITEMS = [
  { key: "emailNotifications", label: "Email notifications", description: "General account and product updates." },
  { key: "paymentNotifications", label: "Payment notifications", description: "Alerts when a payment succeeds or fails." },
  { key: "invoiceNotifications", label: "Invoice notifications", description: "Reminders about invoices you send and receive." },
  { key: "subscriptionNotifications", label: "Subscription notifications", description: "Renewals, cancellations, and billing changes." },
  { key: "securityNotifications", label: "Security notifications", description: "Sign-ins and changes to your account security." },
];

function NotificationsTab() {
  const [form, setForm] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  useEffect(() => {
    getNotifications()
      .then(({ data }) => setForm(data.data))
      .catch(() => setError("Unable to load notification preferences."))
      .finally(() => setLoading(false));
  }, []);
  async function save(event) {
    event.preventDefault();
    setSaving(true);
    setError("");
    setNotice("");
    try {
      await updateNotifications(form);
      setNotice("Notification preferences updated.");
    } catch (requestError) {
      setError(requestError.response?.data?.message || "Unable to update notification preferences.");
    } finally {
      setSaving(false);
    }
  }
  if (loading) return <div className="settings-section"><div className="data-loading">Loading notification preferences...</div></div>;
  if (!form) return <div className="settings-section"><div className="data-empty">{error || "Unable to load notification preferences."}</div></div>;
  return (
    <form onSubmit={save}>
      <div className="settings-section">
        <div>
          <h3>Notification preferences</h3>
          <p>Choose which updates PayFlow sends you.</p>
        </div>
        <Notice message={notice} />
        <Notice message={error} tone="error" />
        <div style={{ display: "flex", flexDirection: "column", gap: 14, marginTop: 20 }}>
          {NOTIFICATION_ITEMS.map((item) => (
            <label key={item.key} style={{ display: "flex", alignItems: "flex-start", gap: 10, fontSize: 12, color: "var(--muted)" }}>
              <input
                type="checkbox"
                checked={Boolean(form[item.key])}
                onChange={(e) => setForm({ ...form, [item.key]: e.target.checked })}
                style={{ marginTop: 3 }}
              />
              <span>
                <strong style={{ display: "block", color: "var(--ink)", fontSize: 13 }}>{item.label}</strong>
                {item.description}
              </span>
            </label>
          ))}
        </div>
      </div>
      <div className="form-actions">
        <button className="button" disabled={saving}>{saving ? "Saving..." : "Save changes"}</button>
      </div>
    </form>
  );
}

function SecurityTab() {
  const { user } = useAuth();
  const [form, setForm] = useState({ currentPassword: "", newPassword: "", confirmPassword: "" });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  async function save(event) {
    event.preventDefault();
    setError("");
    setNotice("");
    if (form.newPassword !== form.confirmPassword) {
      setError("New password and confirmation do not match.");
      return;
    }
    setSaving(true);
    try {
      await changePassword({ currentPassword: form.currentPassword, newPassword: form.newPassword });
      setNotice("Password changed successfully.");
      setForm({ currentPassword: "", newPassword: "", confirmPassword: "" });
    } catch (requestError) {
      setError(requestError.response?.data?.message || "Unable to change your password.");
    } finally {
      setSaving(false);
    }
  }
  return (
    <>
      <div className="settings-section security-row">
        <div className="setting-icon">
          <ShieldCheck size={21} />
        </div>
        <div>
          <h3>Account security</h3>
          <p>Signed in as {user?.email}. Email {user?.emailVerified ? "verified" : "not verified"}, phone {user?.phoneVerified ? "verified" : "not verified"}.</p>
        </div>
      </div>
      <form onSubmit={save}>
        <div className="settings-section">
          <div>
            <h3>Change password</h3>
            <p>Use a strong password you don't use anywhere else.</p>
          </div>
          <Notice message={notice} />
          <Notice message={error} tone="error" />
          <div className="form-grid">
            <PasswordInput
              label="Current password"
              value={form.currentPassword}
              onChange={(e) => setForm({ ...form, currentPassword: e.target.value })}
              required
            />
            <div>
              <PasswordInput
                label="New password"
                value={form.newPassword}
                onChange={(e) => setForm({ ...form, newPassword: e.target.value })}
                required
              />
              {form.newPassword && <PasswordStrength password={form.newPassword} />}
            </div>
            <PasswordInput
              label="Confirm new password"
              value={form.confirmPassword}
              onChange={(e) => setForm({ ...form, confirmPassword: e.target.value })}
              required
            />
          </div>
        </div>
        <div className="form-actions">
          <button className="button" disabled={saving}>{saving ? "Saving..." : "Update password"}</button>
        </div>
      </form>
    </>
  );
}

export default function Settings() {
  const [tab, setTab] = useState("Profile");
  return (
    <>
      <div className="page-heading">
        <div>
          <h2>Settings</h2>
          <p>Manage your account and business preferences.</p>
        </div>
      </div>
      <div className="settings-layout">
        <div className="settings-nav">
          {TABS.map((item) => (
            <button key={item} className={tab === item ? "selected" : ""} onClick={() => setTab(item)}>
              {item}
            </button>
          ))}
        </div>
        <div className="settings-content">
          {tab === "Profile" && <ProfileTab />}
          {tab === "Business" && <BusinessTab />}
          {tab === "Notifications" && <NotificationsTab />}
          {tab === "Security" && <SecurityTab />}
        </div>
      </div>
    </>
  );
}
