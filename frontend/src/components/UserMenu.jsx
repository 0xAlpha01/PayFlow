import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { LogOut, Settings as SettingsIcon, User } from "lucide-react";
import { useAuth } from "../context/AuthContext";

export default function UserMenu({ trigger, triggerClassName = "", direction = "down", align = "left" }) {
  const [open, setOpen] = useState(false);
  const [loggingOut, setLoggingOut] = useState(false);
  const ref = useRef(null);
  const { logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!open) return undefined;
    function onOutside(event) {
      if (ref.current && !ref.current.contains(event.target)) setOpen(false);
    }
    function onEscape(event) {
      if (event.key === "Escape") setOpen(false);
    }
    document.addEventListener("mousedown", onOutside);
    document.addEventListener("keydown", onEscape);
    return () => {
      document.removeEventListener("mousedown", onOutside);
      document.removeEventListener("keydown", onEscape);
    };
  }, [open]);

  async function handleLogout() {
    if (loggingOut) return;
    setOpen(false);
    setLoggingOut(true);
    try {
      await logout();
    } catch {
      /* local state is cleared by AuthContext regardless */
    }
  }
  function go(path) {
    setOpen(false);
    navigate(path);
  }

  return (
    <div className="user-menu" ref={ref}>
      <button
        type="button"
        className={triggerClassName}
        onClick={() => setOpen((value) => !value)}
        aria-haspopup="menu"
        aria-expanded={open}
        disabled={loggingOut}
      >
        {trigger}
      </button>
      {open && (
        <div className={`user-menu-dropdown ${direction === "up" ? "drop-up" : "drop-down"} align-${align}`} role="menu">
          <button type="button" role="menuitem" onClick={() => go("/settings/profile")}>
            <User size={15} /> Profile
          </button>
          <button type="button" role="menuitem" onClick={() => go("/settings/profile")}>
            <SettingsIcon size={15} /> Settings
          </button>
          <button type="button" role="menuitem" className="user-menu-logout" onClick={handleLogout} disabled={loggingOut}>
            <LogOut size={15} /> {loggingOut ? "Logging out..." : "Logout"}
          </button>
        </div>
      )}
    </div>
  );
}
