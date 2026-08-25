import { Link, useLocation, useNavigate } from "react-router-dom";
import { LockKeyhole } from "lucide-react";
import {
  AuthButton,
  AuthError,
  AuthInput,
  AuthLayout,
  PasswordInput,
} from "../../components/auth";
import { useAuth } from "../../context/AuthContext";
import { useState } from "react";

export default function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [form, setForm] = useState({ email: "", password: "" });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const update = (event) =>
    setForm({ ...form, [event.target.name]: event.target.value });
  async function submit(event) {
    event.preventDefault();
    setError("");
    if (!form.email || !form.password)
      return setError("Enter your email and password.");
    setLoading(true);
    try {
      await login(form);
      navigate(location.state?.from || "/dashboard", { replace: true });
    } catch (requestError) {
      setError(
        requestError.response?.data?.message ||
          "Unable to sign in. Check your details and try again.",
      );
    } finally {
      setLoading(false);
    }
  }
  return (
    <AuthLayout title="Welcome back" subtitle="Sign in to your PayFlow account">
      <form className="auth-form" onSubmit={submit}>
        <AuthError message={error} />
        <AuthInput
          label="Email address"
          name="email"
          type="email"
          autoComplete="email"
          value={form.email}
          onChange={update}
          placeholder="you@business.com"
        />
        <PasswordInput
          label="Password"
          name="password"
          autoComplete="current-password"
          value={form.password}
          onChange={update}
          placeholder="Enter your password"
        />
        <div className="auth-row">
          <label className="check-label">
            <input type="checkbox" /> Remember me
          </label>
          <Link to="/forgot-password">Forgot password?</Link>
        </div>
        <AuthButton loading={loading} type="submit">
          <LockKeyhole size={16} /> Sign in
        </AuthButton>
        <p className="auth-switch">
          Don't have an account? <Link to="/signup">Create an account</Link>
        </p>
      </form>
    </AuthLayout>
  );
}
