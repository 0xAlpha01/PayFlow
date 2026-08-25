import { useState } from "react";
import { Link } from "react-router-dom";
import {
  AuthButton,
  AuthError,
  AuthLayout,
  AuthSuccess,
  PasswordInput,
  PasswordStrength,
} from "../../components/auth";
import { authService } from "../../services/authService";
export default function ResetPassword() {
  const [form, setForm] = useState({ password: "", confirm: "" });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const token =
    new URLSearchParams(window.location.search).get("token") ||
    sessionStorage.getItem("payflow_reset_token") ||
    "";
  async function submit(event) {
    event.preventDefault();
    if (
      form.password.length < 8 ||
      !/[A-Z]/.test(form.password) ||
      !/[a-z]/.test(form.password) ||
      !/\d/.test(form.password) ||
      !/[^A-Za-z\d]/.test(form.password)
    )
      return setError(
        "Use at least 8 characters with uppercase, lowercase, number, and special character.",
      );
    if (form.password !== form.confirm)
      return setError("Passwords do not match.");
    if (!token) return setError("This reset link is missing a valid token.");
    setError("");
    setLoading(true);
    try {
      await authService.resetPassword({ token, password: form.password });
      setSuccess("Password reset successfully.");
      sessionStorage.removeItem("payflow_reset_token");
    } catch (requestError) {
      setError(
        requestError.response?.data?.message ||
          "This reset link is invalid or expired.",
      );
    } finally {
      setLoading(false);
    }
  }
  return (
    <AuthLayout
      title="Create a new password"
      subtitle="Choose a strong password for your PayFlow account."
    >
      <form className="auth-form" onSubmit={submit}>
        <AuthError message={error} />
        <AuthSuccess message={success} />
        <PasswordInput
          label="New password"
          value={form.password}
          onChange={(event) =>
            setForm({ ...form, password: event.target.value })
          }
          placeholder="Create a strong password"
        />
        <PasswordStrength password={form.password} />
        <PasswordInput
          label="Confirm password"
          value={form.confirm}
          onChange={(event) =>
            setForm({ ...form, confirm: event.target.value })
          }
          placeholder="Repeat your password"
        />
        <AuthButton loading={loading} type="submit">
          Reset password
        </AuthButton>
        {success ? (
          <Link className="back-link" to="/login">
            Return to sign in
          </Link>
        ) : (
          <Link className="back-link" to="/login">
            Back to sign in
          </Link>
        )}
      </form>
    </AuthLayout>
  );
}
