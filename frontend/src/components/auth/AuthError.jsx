export default function AuthError({ message }) {
  return message ? (
    <div className="auth-alert error" role="alert">
      {message}
    </div>
  ) : null;
}
