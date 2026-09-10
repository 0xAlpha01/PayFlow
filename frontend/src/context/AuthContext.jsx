/* eslint-disable react-refresh/only-export-components, react-hooks/set-state-in-effect */
import { createContext, useContext, useEffect, useState } from "react";
import { authService } from "../services/authService";
import { configureAuth } from "../services/api";
import { connectSocket, disconnectSocket } from "../services/socket";

const AuthContext = createContext(null);
const refreshKey = "payflow_refresh_token";

function readRefreshToken() {
  return sessionStorage.getItem(refreshKey);
}
function saveSession(session) {
  sessionStorage.setItem(refreshKey, session.refreshToken);
  configureAuth(session.accessToken, session.refreshToken);
  return session;
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [accessToken, setAccessToken] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    configureAuth(
      null,
      readRefreshToken(),
      (session) => {
        setAccessToken(session.accessToken);
        saveSession(session);
      },
      clearSession,
    );
    const token = readRefreshToken();
    if (!token) {
      setLoading(false);
      return;
    }
    authService
      .refreshToken(token)
      .then(({ data }) => {
        const session = saveSession({
          accessToken: data.data.accessToken,
          refreshToken: token,
        });
        setAccessToken(session.accessToken);
        return authService.getCurrentUser();
      })
      .then(({ data }) => setUser(data.data.user))
      .catch(clearSession)
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (accessToken) connectSocket(accessToken);
    else disconnectSocket();
  }, [accessToken]);

  function clearSession() {
    sessionStorage.removeItem(refreshKey);
    setAccessToken(null);
    setUser(null);
    configureAuth(null, null);
  }
  async function login(credentials) {
    const { data } = await authService.login(credentials);
    const session = saveSession(data.data);
    setAccessToken(session.accessToken);
    setUser(session.user);
    return session.user;
  }
  function completeEmailVerification(session) {
    const saved = saveSession(session);
    setAccessToken(saved.accessToken);
    setUser(saved.user);
  }
  function markPhoneVerified() {
    setUser((current) => (current ? { ...current, phoneVerified: true } : current));
  }
  async function refreshUser() {
    const { data } = await authService.getCurrentUser();
    setUser(data.data.user);
    return data.data.user;
  }
  async function signup(details) {
    return authService.signup(details);
  }
  async function logout() {
    try {
      if (readRefreshToken()) await authService.logout(readRefreshToken());
    } finally {
      clearSession();
    }
  }
  async function refreshToken() {
    const { data } = await authService.refreshToken(readRefreshToken());
    setAccessToken(data.data.accessToken);
    configureAuth(data.data.accessToken, readRefreshToken());
    return data.data.accessToken;
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        accessToken,
        isAuthenticated: Boolean(user && accessToken),
        loading,
        login,
        completeEmailVerification,
        markPhoneVerified,
        refreshUser,
        signup,
        logout,
        refreshToken,
        clearSession,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
export function useAuth() {
  return useContext(AuthContext);
}
