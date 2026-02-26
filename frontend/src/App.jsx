import { useState, useEffect } from "react";
import "./App.css";

const API_BASE = "http://localhost:8080";

// ── API layer ──────────────────────────────────────────────
async function apiLogin(login, password) {
  const res = await fetch(`${API_BASE}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Login: login, Password: password }),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || "ACCESS DENIED");
  return data;
}

async function apiRegister(username, login, password) {
  const res = await fetch(`${API_BASE}/registration`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Username: username, Login: login, Password: password }),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || "REGISTRATION FAILED");
  return data;
}

// ── Corner decorations ────────────────────────────────────
function Corners() {
  return (
    <>
      <div className="corner corner-tl" />
      <div className="corner corner-tr" />
      <div className="corner corner-bl" />
      <div className="corner corner-br" />
    </>
  );
}

// ── Login form ────────────────────────────────────────────
function LoginForm({ onSuccess }) {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async () => {
    if (!login || !password) { setError("ALL FIELDS REQUIRED"); return; }
    setError(""); setLoading(true);
    try {
      await apiLogin(login, password);
      onSuccess(login);
    } catch (e) {
      setError(e.message);
    }
    setLoading(false);
  };

  return (
    <>
      <div className="field-group">
        <label className="field-label">// IDENTIFIER</label>
        <input
          className="field-input"
          value={login}
          onChange={e => setLogin(e.target.value)}
          placeholder="user@domain.net"
          autoComplete="off"
        />
      </div>
      <div className="field-group">
        <label className="field-label">// AUTH KEY</label>
        <input
          className="field-input"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          placeholder="••••••••••••"
          onKeyDown={e => e.key === "Enter" && handleSubmit()}
        />
      </div>
      <button className="btn-primary" onClick={handleSubmit} disabled={loading}>
        {loading ? "AUTHENTICATING..." : "ENTER SYSTEM"}
      </button>
      {error && <div className="alert-error">⚠ {error}</div>}
    </>
  );
}

// ── Register form ─────────────────────────────────────────
function RegisterForm({ onRegistered }) {
  const [username, setUsername] = useState("");
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async () => {
    if (!username || !login || !password) { setError("ALL FIELDS REQUIRED"); return; }
    setError(""); setLoading(true);
    try {
      await apiRegister(username, login, password);
      onRegistered(login);
    } catch (e) {
      setError(e.message);
    }
    setLoading(false);
  };

  return (
    <>
      <div className="field-group">
        <label className="field-label">// ALIAS</label>
        <input
          className="field-input"
          value={username}
          onChange={e => setUsername(e.target.value)}
          placeholder="your_handle"
        />
      </div>
      <div className="field-group">
        <label className="field-label">// IDENTIFIER</label>
        <input
          className="field-input"
          value={login}
          onChange={e => setLogin(e.target.value)}
          placeholder="user@domain.net"
          autoComplete="off"
        />
      </div>
      <div className="field-group">
        <label className="field-label">// AUTH KEY</label>
        <input
          className="field-input"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          placeholder="••••••••••••"
          onKeyDown={e => e.key === "Enter" && handleSubmit()}
        />
      </div>
      <button className="btn-primary" onClick={handleSubmit} disabled={loading}>
        {loading ? "REGISTERING..." : "CREATE IDENTITY"}
      </button>
      {error && <div className="alert-error">⚠ {error}</div>}
    </>
  );
}

// ── Auth panel ────────────────────────────────────────────
function AuthPanel({ onSuccess }) {
  const [mode, setMode] = useState("login");
  const [success, setSuccess] = useState("");

  const handleRegistered = () => {
    setSuccess("IDENTITY REGISTERED // PROCEED TO AUTH");
    setMode("login");
  };

  return (
    <div className="panel">
      <Corners />
      <div className="panel-logo">SYS // v2.4.1</div>
      <div className="panel-title">NETWATCH</div>
      <div className="panel-subtitle">// DISTRIBUTED MONITORING SYSTEM</div>

      <div className="tabs">
        <button
          className={`tab ${mode === "login" ? "active" : ""}`}
          onClick={() => { setMode("login"); setSuccess(""); }}
        >
          AUTH
        </button>
        <button
          className={`tab ${mode === "register" ? "active" : ""}`}
          onClick={() => { setMode("register"); setSuccess(""); }}
        >
          REGISTER
        </button>
      </div>

      {mode === "login"
        ? <LoginForm onSuccess={onSuccess} />
        : <RegisterForm onRegistered={handleRegistered} />
      }

      {success && <div className="alert-success">✓ {success}</div>}
    </div>
  );
}

// ── Dashboard ─────────────────────────────────────────────
function Dashboard({ username, onLogout }) {
  const [time, setTime] = useState(new Date());

  useEffect(() => {
    const t = setInterval(() => setTime(new Date()), 1000);
    return () => clearInterval(t);
  }, []);

  const cards = [
    { label: "ACTIVE NODES", value: "—" },
    { label: "UPTIME",       value: "—" },
    { label: "ALERTS",       value: "—" },
  ];

  return (
    <div className="dashboard">
      <div className="bg-grid" />
      <div className="scanline" />

      <div className="topbar">
        <div className="topbar-logo">NETWATCH</div>
        <div className="topbar-right">
          <span><span className="status-dot" />ONLINE</span>
          <span className="topbar-username">{username}</span>
          <span>{time.toLocaleTimeString()}</span>
          <button className="btn-logout" onClick={onLogout}>DISCONNECT</button>
        </div>
      </div>

      <div className="dash-content">
        <div className="dash-title">SYSTEM ONLINE</div>
        <div className="dash-subtitle">// MONITORING DASHBOARD — AWAITING DATA STREAMS</div>
        <div className="cards">
          {cards.map((c, i) => (
            <div key={i} className="card">
              <Corners />
              <div className="card-label">{c.label}</div>
              <div className="card-value">{c.value}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

// ── Root ──────────────────────────────────────────────────
export default function App() {
  const [user, setUser] = useState(null);

  if (user) {
    return <Dashboard username={user} onLogout={() => setUser(null)} />;
  }

  return (
    <div className="app-root">
      <div className="bg-grid" />
      <div className="scanline" />
      <AuthPanel onSuccess={setUser} />
    </div>
  );
}