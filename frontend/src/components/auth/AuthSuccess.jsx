export default function AuthSuccess({ message }) {
  return message ? (
    <div className="auth-alert success" role="status">
      {message}
    </div>
  ) : null;
}
