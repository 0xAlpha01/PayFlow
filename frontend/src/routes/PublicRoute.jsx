import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function PublicRoute() {
  const { isAuthenticated, loading } = useAuth();
  const location = useLocation();
  if (loading)
    return (
      <div className="auth-loading">
        <div className="brand-mark">P</div>
        <span>Loading your workspace...</span>
      </div>
    );
  if (location.pathname === "/verify-phone") {
    return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
  }
  if (!isAuthenticated) return <Outlet />;
  return <Navigate to="/dashboard" replace />;
}
