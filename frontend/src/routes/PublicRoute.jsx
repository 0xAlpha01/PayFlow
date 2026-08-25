import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function PublicRoute() {
  const { isAuthenticated, loading } = useAuth();
  if (loading)
    return (
      <div className="auth-loading">
        <div className="brand-mark">P</div>
        <span>Loading your workspace...</span>
      </div>
    );
  return isAuthenticated ? <Navigate to="/dashboard" replace /> : <Outlet />;
}
