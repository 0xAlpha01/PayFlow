export default function PasswordStrength({ password }) {
  const checks = [
    password.length >= 8,
    /[A-Z]/.test(password),
    /[a-z]/.test(password),
    /\d/.test(password),
    /[^A-Za-z\d]/.test(password),
  ];
  const score = checks.filter(Boolean).length;
  return (
    <div className="password-strength">
      <div className="strength-bars">
        {checks.map((valid, index) => (
          <i className={valid ? "valid" : ""} key={index} />
        ))}
      </div>
      <span>
        {score < 3
          ? "Use a stronger password"
          : score < 5
            ? "Almost there"
            : "Strong password"}
      </span>
    </div>
  );
}
