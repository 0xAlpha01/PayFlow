import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Smartphone } from "lucide-react";
import { authService } from "../../services/authService";
import {
  AuthButton,
  AuthError,
  AuthLayout,
  AuthSuccess,
  OTPInput,
} from "../../components/auth";

function maskPhone(phone) {
  return phone
    ? `${phone.slice(0, 4)}******${phone.slice(-3)}`
    : "your phone number";
}
export default function VerifyPhone() {
  const phone = sessionStorage.getItem("payflow_pending_phone") || "";
  const navigate = useNavigate();
  const [code, setCode] = useState("");
  const [seconds, setSeconds] = useState(600);
  const [cooldown, setCooldown] = useState(0);
  const [loading, setLoading] = useState(false);
  const [sending, setSending] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  useEffect(() => {
    const timer = setInterval(() => {
      setSeconds((value) => Math.max(value - 1, 0));
      setCooldown((value) => Math.max(value - 1, 0));
    }, 1000);
    authService
      .sendPhoneCode(phone)
      .then(() => setSuccess("Verification code sent."))
      .catch((requestError) =>
        setError(
          requestError.response?.data?.message ||
            "Unable to send a verification code.",
        ),
      )
      .finally(() => setSending(false));
    return () => clearInterval(timer);
  }, [phone]);
  async function verify(event) {
    event.preventDefault();
    if (code.length !== 6) return setError("Enter all 6 digits.");
    setError("");
    setLoading(true);
    try {
      await authService.verifyPhone({ phone, code });
      sessionStorage.removeItem("payflow_pending_email");
      sessionStorage.removeItem("payflow_pending_phone");
      navigate("/login", {
        state: { message: "Phone verified. Sign in to continue." },
      });
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
    try {
      await authService.sendPhoneCode(phone);
      setSuccess("A new verification code has been sent.");
    } catch {
      setError("Unable to resend the code right now.");
    }
  }
  const time = `${String(Math.floor(seconds / 60)).padStart(2, "0")}:${String(seconds % 60).padStart(2, "0")}`;
  return (
    <AuthLayout
      title="Verify your phone number"
      subtitle={`We sent a 6-digit code to ${maskPhone(phone)}`}
    >
      <form className="auth-form otp-form" onSubmit={verify}>
        <div className="auth-icon">
          <Smartphone size={22} />
        </div>
        <AuthError message={error} />
        <AuthSuccess message={success} />
        <OTPInput value={code} onChange={setCode} />
        <span className="otp-timer">Code expires in {time}</span>
        <AuthButton loading={loading || sending} type="submit">
          {sending ? "Sending code" : "Verify phone"}
        </AuthButton>
        <p className="resend-copy">
          Didn't receive the code?{" "}
          <button type="button" onClick={resend} disabled={cooldown}>
            {cooldown ? `Resend in ${cooldown}s` : "Resend code"}
          </button>
        </p>
      </form>
    </AuthLayout>
  );
}
