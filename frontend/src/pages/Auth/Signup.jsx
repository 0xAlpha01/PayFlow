import { Link, useNavigate } from "react-router-dom";
import { useState } from "react";
import {
  AuthButton,
  AuthError,
  AuthInput,
  AuthLayout,
  PasswordInput,
  PasswordStrength,
} from "../../components/auth";
import { useAuth } from "../../context/AuthContext";

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const phonePattern = /^\+[1-9]\d{7,14}$/;
export default function Signup() {
  const { signup } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({
    firstName: "",
    lastName: "",
    businessName: "",
    email: "",
    phone: "",
    password: "",
    confirmPassword: "",
    terms: false,
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const update = (event) =>
    setForm({
      ...form,
      [event.target.name]:
        event.target.type === "checkbox"
          ? event.target.checked
          : event.target.value,
    });
  function validate() {
    if (
      !form.firstName ||
      !form.lastName ||
      !form.businessName ||
      !form.email ||
      !form.phone ||
      !form.password
    )
      return "Complete all required fields.";
    if (!emailPattern.test(form.email)) return "Enter a valid email address.";
    if (!phonePattern.test(form.phone))
      return "Use an international phone number, for example +2348012345678.";
    if (
      !/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z\d]).{8,}$/.test(
        form.password,
      )
    )
      return "Password must include uppercase, lowercase, number, and special character.";
    if (form.password !== form.confirmPassword)
      return "Passwords do not match.";
    if (!form.terms) return "Accept the Terms and Privacy Policy to continue.";
    return "";
  }
  async function submit(event) {
    event.preventDefault();
    const validation = validate();
    if (validation) return setError(validation);
    setError("");
    setLoading(true);
    try {
      await signup({
        firstName: form.firstName,
        lastName: form.lastName,
        businessName: form.businessName,
        email: form.email,
        phone: form.phone,
        password: form.password,
      });
      sessionStorage.setItem("payflow_pending_email", form.email);
      sessionStorage.setItem("payflow_pending_phone", form.phone);
      navigate("/verify-email");
    } catch (requestError) {
      setError(
        requestError.response?.data?.message ||
          "Unable to create your account. Please try again.",
      );
    } finally {
      setLoading(false);
    }
  }
  return (
    <AuthLayout
      title="Create your account"
      subtitle="Start managing your payments with PayFlow"
    >
      <form className="auth-form signup-form" onSubmit={submit}>
        <AuthError message={error} />
        <div className="auth-two-col">
          <AuthInput
            label="First name"
            name="firstName"
            value={form.firstName}
            onChange={update}
            placeholder="John"
          />
          <AuthInput
            label="Last name"
            name="lastName"
            value={form.lastName}
            onChange={update}
            placeholder="Doe"
          />
        </div>
        <AuthInput
          label="Business name"
          name="businessName"
          value={form.businessName}
          onChange={update}
          placeholder="Your company"
        />
        <AuthInput
          label="Email address"
          name="email"
          type="email"
          value={form.email}
          onChange={update}
          placeholder="you@business.com"
        />
        <AuthInput
          label="Phone number"
          name="phone"
          type="tel"
          value={form.phone}
          onChange={update}
          placeholder="+2348012345678"
        />
        <PasswordInput
          label="Password"
          name="password"
          value={form.password}
          onChange={update}
          placeholder="Create a strong password"
        />
        <PasswordStrength password={form.password} />
        <PasswordInput
          label="Confirm password"
          name="confirmPassword"
          value={form.confirmPassword}
          onChange={update}
          placeholder="Repeat your password"
        />
        <label className="check-label terms">
          <input
            type="checkbox"
            name="terms"
            checked={form.terms}
            onChange={update}
          />{" "}
          I agree to the <a href="/terms">Terms</a> and{" "}
          <a href="/privacy">Privacy Policy</a>
        </label>
        <AuthButton loading={loading} type="submit">
          Create account
        </AuthButton>
        <p className="auth-switch">
          Already have an account? <Link to="/login">Sign in</Link>
        </p>
      </form>
    </AuthLayout>
  );
}
