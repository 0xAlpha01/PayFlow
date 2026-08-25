export default function AuthInput({ label, error, ...props }) {
  return (
    <label className="auth-field">
      {label}
      <input {...props} />
      {error && <small className="auth-field-error">{error}</small>}
    </label>
  );
}
