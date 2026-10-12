import { apiBaseUrl } from "./apiBase";

const listeners = new Map();
let socket = null;
let reconnectTimer = null;
let currentToken = null;

function wsUrl(token) {
  return `${apiBaseUrl.replace(/^http/, "ws")}/ws?token=${encodeURIComponent(token)}`;
}

export function connectSocket(token) {
  if (!token) return;
  currentToken = token;
  if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) return;
  clearTimeout(reconnectTimer);
  const next = new WebSocket(wsUrl(token));
  next.onmessage = (event) => {
    let parsed;
    try {
      parsed = JSON.parse(event.data);
    } catch {
      return;
    }
    listeners.get(parsed.type)?.forEach((callback) => callback(parsed.data));
  };
  next.onclose = () => {
    if (socket === next) socket = null;
    if (currentToken) reconnectTimer = setTimeout(() => connectSocket(currentToken), 3000);
  };
  next.onerror = () => next.close();
  socket = next;
}
export function disconnectSocket() {
  currentToken = null;
  clearTimeout(reconnectTimer);
  socket?.close();
  socket = null;
}
export function onSocketEvent(type, callback) {
  if (!listeners.has(type)) listeners.set(type, new Set());
  listeners.get(type).add(callback);
  return () => listeners.get(type)?.delete(callback);
}
