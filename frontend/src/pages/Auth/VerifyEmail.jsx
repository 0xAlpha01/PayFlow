import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { MailCheck } from "lucide-react";
import { authService } from "../../services/authService";
import { useAuth } from "../../context/AuthContext";
import {
  AuthButton,
  AuthError,
  AuthLayout,
  AuthSuccess,
  OTPInput,
} from "../../components/auth";

export default function VerifyEmail() {
  const email = sessionStorage.getItem("payflow_pending_email") || "";
  const navigate = useNavigate();
  const { completeEmailVerification } = useAuth();
  const [code, setCode] = useState("");
  const [seconds, setSeconds] = useState(600);
  const [cooldown, setCooldown] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  useEffect(() => {
    const timer = setInterval(
      () => setSeconds((value) => Math.max(value - 1, 0)),
      1000,
    );
    return () => clearInterval(timer);
  }, []);
  async function verify(event) {
    event.preventDefault();
    if (code.length !== 6) return setError("Enter all 6 digits.");
    setError("");
    setLoading(true);
    try {
      const { data } = await authService.verifyEmail({ email, code });
      completeEmailVerification(data.data);
      sessionStorage.removeItem("payflow_pending_email");
      sessionStorage.removeItem("payflow_pending_phone");
      setSuccess("Email verified. Taking you to your dashboard.");
      setTimeout(() => navigate("/dashboard", { replace: true }), 600);
    } catch (requestError) {
      setError(
        requestError.response?.data?.message ||
          "That code is invalid or expired.",
      );
    } finally {
      setLoading(false);
    }
  }
  async function resend() {
    if (cooldown) return;
    setCooldown(60);
    setError("");
    try {
      await authService.resendEmailVerification(email);
      setSuccess("A new verification code has been sent.");
    } catch {
      setError("Unable to resend the code right now.");
    }
  }
  const time = `${String(Math.floor(seconds / 60)).padStart(2, "0")}:${String(seconds % 60).padStart(2, "0")}`;
  return (
    <AuthLayout
      title="Verify your email"
      subtitle={`We sent a 6-digit code to ${email || "your email address"}`}
    >
      <form className="auth-form otp-form" onSubmit={verify}>
        <div className="auth-icon">
          <MailCheck size={22} />
        </div>
        <AuthError message={error} />
        <AuthSuccess message={success} />
        <OTPInput value={code} onChange={setCode} />
        <span className="otp-timer">Code expires in {time}</span>
        <AuthButton loading={loading} type="submit">
          Verify email
        </AuthButton>
        <p className="resend-copy">
          Didn't receive the code?{" "}
          <button type="button" onClick={resend} disabled={cooldown}>
            {cooldown ? `Resend in ${cooldown}s` : "Resend code"}
          </button>
        </p>
        <Link className="back-link" to="/login">
          Back to sign in
        </Link>
      </form>
    </AuthLayout>
  );
}
