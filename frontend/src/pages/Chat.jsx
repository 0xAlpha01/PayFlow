import { useEffect, useRef, useState } from "react";
import { ArrowLeft, MessageSquare, Search, Send } from "lucide-react";
import {
  getConversations,
  getMessages,
  markConversationRead,
  searchChatUsers,
  sendMessage,
  startConversation,
} from "../services/chatService";
import { onSocketEvent } from "../services/socket";
import { useAuth } from "../context/AuthContext";

function initials(user) {
  return `${user?.firstName?.[0] || ""}${user?.lastName?.[0] || ""}`.toUpperCase() || "?";
}
function timeLabel(value) {
  if (!value) return "";
  const date = new Date(value);
  const sameDay = date.toDateString() === new Date().toDateString();
  return sameDay ? date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }) : date.toLocaleDateString();
}

export default function ChatPage() {
  const { user } = useAuth();
  const [conversations, setConversations] = useState([]);
  const [loadingConversations, setLoadingConversations] = useState(true);
  const [active, setActive] = useState(null);
  const [messages, setMessages] = useState([]);
  const [loadingMessages, setLoadingMessages] = useState(false);
  const [draft, setDraft] = useState("");
  const [sending, setSending] = useState(false);
  const [query, setQuery] = useState("");
  const [searchResults, setSearchResults] = useState([]);
  const [searching, setSearching] = useState(false);
  const [mobileShowChat, setMobileShowChat] = useState(false);
  const [error, setError] = useState("");
  const searchTimer = useRef(null);
  const scrollRef = useRef(null);

  async function loadConversations() {
    setLoadingConversations(true);
    try {
      const { data } = await getConversations();
      setConversations(data.data || []);
    } catch {
      setError("Unable to load conversations.");
    } finally {
      setLoadingConversations(false);
    }
  }
  useEffect(() => {
    loadConversations();
  }, []);

  useEffect(() => {
    return onSocketEvent("message", (message) => {
      loadConversations();
      setActive((current) => {
        if (current && current.conversation.id === message.conversationId) {
          setMessages((existing) => [...existing, message]);
          markConversationRead(message.conversationId).catch(() => {});
        }
        return current;
      });
    });
  }, []);

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight, behavior: "smooth" });
  }, [messages]);

  function onQueryChange(value) {
    setQuery(value);
    clearTimeout(searchTimer.current);
    if (!value.trim()) {
      setSearchResults([]);
      return;
    }
    searchTimer.current = setTimeout(async () => {
      setSearching(true);
      try {
        const { data } = await searchChatUsers(value.trim());
        setSearchResults(data.data || []);
      } catch {
        setSearchResults([]);
      } finally {
        setSearching(false);
      }
    }, 300);
  }
  async function openWithUser(person) {
    try {
      const { data } = await startConversation(person.id);
      const conversation = data.data;
      setQuery("");
      setSearchResults([]);
      await loadConversations();
      await openConversation({ conversation, peer: person, unreadCount: 0 });
    } catch (requestError) {
      setError(requestError.response?.data?.message || "Unable to start this conversation.");
    }
  }
  async function openConversation(summary) {
    setActive(summary);
    setMobileShowChat(true);
    setLoadingMessages(true);
    setMessages([]);
    try {
      const { data } = await getMessages(summary.conversation.id, { page: 1, limit: 100 });
      setMessages(data.data || []);
      if (summary.unreadCount > 0) {
        await markConversationRead(summary.conversation.id);
        setConversations((current) =>
          current.map((item) => (item.conversation.id === summary.conversation.id ? { ...item, unreadCount: 0 } : item)),
        );
      }
    } catch {
      setError("Unable to load this conversation.");
    } finally {
      setLoadingMessages(false);
    }
  }
  async function submit(event) {
    event.preventDefault();
    const content = draft.trim();
    if (!content || !active) return;
    setSending(true);
    setDraft("");
    try {
      const { data } = await sendMessage(active.conversation.id, content);
      setMessages((current) => [...current, data.data]);
      loadConversations();
    } catch (requestError) {
      setError(requestError.response?.data?.message || "Unable to send that message.");
      setDraft(content);
    } finally {
      setSending(false);
    }
  }

  return (
    <>
      <div className="page-heading">
        <div>
          <h2>Messages</h2>
          <p>Chat privately with other PayFlow users.</p>
        </div>
      </div>
      {error && <div className="data-empty">{error}</div>}
      <div className={`chat-shell panel ${mobileShowChat ? "mobile-show-chat" : ""}`}>
        <aside className="chat-sidebar">
          <div className="chat-search">
            <Search size={16} />
            <input placeholder="Find people to message" value={query} onChange={(e) => onQueryChange(e.target.value)} />
          </div>
          {query && (
            <div className="chat-search-results">
              {searching ? (
                <div className="data-loading">Searching...</div>
              ) : searchResults.length === 0 ? (
                <div className="data-empty">No users found.</div>
              ) : (
                searchResults.map((person) => (
                  <button type="button" key={person.id} className="conversation-item" onClick={() => openWithUser(person)}>
                    <span className="person-avatar">{initials(person)}</span>
                    <span className="conversation-meta">
                      <strong>{person.firstName} {person.lastName}</strong>
                      <small>{person.email}</small>
                    </span>
                  </button>
                ))
              )}
            </div>
          )}
          <div className="conversation-list">
            {loadingConversations ? (
              <div className="data-loading">Loading conversations...</div>
            ) : conversations.length === 0 ? (
              <div className="data-empty">Search for someone to start chatting.</div>
            ) : (
              conversations.map((summary) => (
                <button
                  type="button"
                  key={summary.conversation.id}
                  className={`conversation-item ${active?.conversation.id === summary.conversation.id ? "active" : ""}`}
                  onClick={() => openConversation(summary)}
                >
                  <span className="person-avatar">{initials(summary.peer)}</span>
                  <span className="conversation-meta">
                    <strong>{summary.peer.firstName} {summary.peer.lastName}</strong>
                    <small>{summary.conversation.lastMessage || "No messages yet"}</small>
                  </span>
                  {summary.unreadCount > 0 && <span className="conversation-unread">{summary.unreadCount}</span>}
                </button>
              ))
            )}
          </div>
        </aside>
        <section className="chat-window">
          {!active ? (
            <div className="chat-empty">
              <MessageSquare size={28} />
              <h3>Select a conversation</h3>
              <p>Choose an existing conversation or search for a PayFlow user to start chatting.</p>
            </div>
          ) : (
            <>
              <div className="chat-header">
                <button type="button" className="chat-back" onClick={() => setMobileShowChat(false)}>
                  <ArrowLeft size={18} />
                </button>
                <span className="person-avatar">{initials(active.peer)}</span>
                <div>
                  <strong>{active.peer.firstName} {active.peer.lastName}</strong>
                  <small>{active.peer.email}</small>
                </div>
              </div>
              <div className="chat-messages" ref={scrollRef}>
                {loadingMessages ? (
                  <div className="data-loading">Loading messages...</div>
                ) : messages.length === 0 ? (
                  <div className="data-empty">No messages yet. Say hello!</div>
                ) : (
                  messages.map((message) => (
                    <div key={message.id} className={`chat-bubble ${message.senderId === user?.id ? "mine" : "theirs"}`}>
                      <p>{message.content}</p>
                      <span>{timeLabel(message.createdAt)}</span>
                    </div>
                  ))
                )}
              </div>
              <form className="chat-input-row" onSubmit={submit}>
                <input
                  placeholder="Type a message"
                  value={draft}
                  onChange={(e) => setDraft(e.target.value)}
                  disabled={sending}
                />
                <button className="button" disabled={sending || !draft.trim()}>
                  <Send size={16} />
                </button>
              </form>
            </>
          )}
        </section>
      </div>
    </>
  );
}
