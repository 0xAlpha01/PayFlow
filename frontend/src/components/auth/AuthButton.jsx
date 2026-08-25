export default function AuthButton({ children, loading, ...props }) {
  return (
    <button
      className="auth-button"
      disabled={loading || props.disabled}
      {...props}
    >
      {loading ? (
        <>
          <span className="spinner" /> {children}...
        </>
      ) : (
        children
      )}
    </button>
  );
}
