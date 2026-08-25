import { useState } from "react";
import { Link } from "react-router-dom";
import {
  AuthButton,
  AuthError,
  AuthInput,
  AuthLayout,
  AuthSuccess,
} from "../../components/auth";
import { authService } from "../../services/authService";
export default function ForgotPassword() {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  async function submit(event) {
    event.preventDefault();
    if (!email) return setError("Enter your email address.");
    setError("");
    setLoading(true);
    try {
      await authService.forgotPassword(email);
      setSuccess("If that email exists, reset instructions have been sent.");
      sessionStorage.setItem("payflow_reset_email", email);
    } catch {
      setError("Unable to send reset instructions right now.");
    } finally {
      setLoading(false);
    }
  }
  return (
    <AuthLayout
      title="Forgot your password?"
      subtitle="Enter your email and we’ll send reset instructions."
    >
      <form className="auth-form" onSubmit={submit}>
        <AuthError message={error} />
        <AuthSuccess message={success} />
        <AuthInput
          label="Email address"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          placeholder="you@business.com"
        />
        <AuthButton loading={loading} type="submit">
          Send reset link
        </AuthButton>
        <Link className="back-link" to="/login">
          Back to sign in
        </Link>
      </form>
    </AuthLayout>
  );
}
