import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";
export default function PasswordInput({ label = "Password", error, ...props }) {
  const [visible, setVisible] = useState(false);
  return (
    <label className="auth-field">
      {label}
      <span className="password-wrap">
        <input {...props} type={visible ? "text" : "password"} />
        <button
          type="button"
          aria-label={visible ? "Hide password" : "Show password"}
          onClick={() => setVisible(!visible)}
        >
          {visible ? <EyeOff size={17} /> : <Eye size={17} />}
        </button>
      </span>
      {error && <small className="auth-field-error">{error}</small>}
    </label>
  );
}
