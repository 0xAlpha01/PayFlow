import { useEffect, useRef, useState } from "react";
import { Bell } from "lucide-react";
import {
  getNotifications,
  getUnreadNotificationCount,
  markAllNotificationsRead,
  markNotificationRead,
} from "../services/notificationService";
import { onSocketEvent } from "../services/socket";

function timeAgo(value) {
  const seconds = Math.floor((Date.now() - new Date(value).getTime()) / 1000);
  if (seconds < 60) return "just now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d ago`;
  return new Date(value).toLocaleDateString();
}

export default function NotificationBell() {
  const [open, setOpen] = useState(false);
  const [items, setItems] = useState([]);
  const [unread, setUnread] = useState(0);
  const [loading, setLoading] = useState(false);
  const ref = useRef(null);

  useEffect(() => {
    getUnreadNotificationCount()
      .then(({ data }) => setUnread(data.data.count))
      .catch(() => {});
    return onSocketEvent("notification", (notification) => {
      setUnread((count) => count + 1);
      setItems((current) => [notification, ...current].slice(0, 30));
    });
  }, []);

  useEffect(() => {
    if (!open) return undefined;
    function onOutside(event) {
      if (ref.current && !ref.current.contains(event.target)) setOpen(false);
    }
    document.addEventListener("mousedown", onOutside);
    return () => document.removeEventListener("mousedown", onOutside);
  }, [open]);

  function toggle() {
    setOpen((value) => {
      const next = !value;
      if (next) {
        setLoading(true);
        getNotifications({ page: 1, limit: 20 })
          .then(({ data }) => setItems(data.data || []))
          .catch(() => {})
          .finally(() => setLoading(false));
      }
      return next;
    });
  }
  async function markRead(notification) {
    if (notification.read) return;
    setItems((current) => current.map((item) => (item.id === notification.id ? { ...item, read: true } : item)));
    setUnread((count) => Math.max(0, count - 1));
    try {
      await markNotificationRead(notification.id);
    } catch {
      /* local state already reflects intent; a background refresh will reconcile */
    }
  }
  async function markAllRead() {
    setItems((current) => current.map((item) => ({ ...item, read: true })));
    setUnread(0);
    try {
      await markAllNotificationsRead();
    } catch {
      /* local state already reflects intent */
    }
  }

  return (
    <div className="notification-menu" ref={ref}>
      <button type="button" className="icon-button notification" onClick={toggle} aria-haspopup="menu" aria-expanded={open}>
        <Bell size={19} />
        {unread > 0 && <i />}
      </button>
      {open && (
        <div className="notification-panel" role="menu">
          <div className="notification-panel-head">
            <strong>Notifications</strong>
            {unread > 0 && (
              <button type="button" onClick={markAllRead}>
                Mark all as read
              </button>
            )}
          </div>
          <div className="notification-panel-list">
            {loading ? (
              <div className="data-loading">Loading notifications...</div>
            ) : items.length === 0 ? (
              <div className="data-empty">You're all caught up.</div>
            ) : (
              items.map((item) => (
                <button
                  type="button"
                  key={item.id}
                  className={`notification-item ${item.read ? "" : "unread"}`}
                  onClick={() => markRead(item)}
                >
                  <span className="notification-item-title">{item.title}</span>
                  <span className="notification-item-message">{item.message}</span>
                  <span className="notification-item-time">{timeAgo(item.createdAt)}</span>
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}
