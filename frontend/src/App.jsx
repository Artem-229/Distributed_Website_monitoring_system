import { useState, useEffect } from "react";
import "./App.css";

const API_BASE = "http://localhost:8080";

async function apiLogin(login, password) {
  const res = await fetch(`${API_BASE}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Login: login, Password: password }),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || "ACCESS DENIED");
  localStorage.setItem("token", data.token);
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

async function apiRequest(url, options = {}) {
  const token = localStorage.getItem("token");
  const res = await fetch(`${API_BASE}/api${url}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`,
      ...options.headers,
    },
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || "REQUEST FAILED");
  return data;
}

async function apiGetMonitors() {
  return apiRequest("/monitors");
}

async function apiAddMonitor(url, timeInterval) {
  return apiRequest("/addmonitor", {
    method: "POST",
    body: JSON.stringify({ Url: url, Time_interval: timeInterval, Is_active: true }),
  });
}

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
      <div className="panel-logo">VERSION // v1.0</div>
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

function AddMonitorForm({ onAdded }) {
  const [url, setUrl] = useState("");
  const [interval, setInterval] = useState(60);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async () => {
    if (!url) { setError("URL REQUIRED"); return; }
    setError(""); setLoading(true);
    try {
      await apiAddMonitor(url, interval);
      setUrl("");
      setInterval(60);
      onAdded();
    } catch (e) {
      setError(e.message);
    }
    setLoading(false);
  };

  return (
    <div className="card" style={{ marginBottom: "32px" }}>
      <Corners />
      <div className="card-label">ADD NEW MONITOR</div>
      <div style={{ marginTop: "16px" }}>
        <div className="field-group">
          <label className="field-label">// TARGET URL</label>
          <input
            className="field-input"
            value={url}
            onChange={e => setUrl(e.target.value)}
            placeholder="https://example.com"
          />
        </div>
        <div className="field-group">
          <label className="field-label">// CHECK INTERVAL (SEC)</label>
          <input
            className="field-input"
            type="number"
            value={interval}
            onChange={e => setInterval(Number(e.target.value))}
            placeholder="60"
          />
        </div>
        <button className="btn-primary" onClick={handleSubmit} disabled={loading} style={{ marginTop: "12px" }}>
          {loading ? "ADDING..." : "ADD MONITOR"}
        </button>
        {error && <div className="alert-error">⚠ {error}</div>}
      </div>
    </div>
  );
}

function MonitorList({ monitors }) {
  if (monitors.length === 0) {
    return (
      <div style={{ color: "var(--text-muted)", fontSize: "12px", letterSpacing: "2px" }}>
        // NO MONITORS CONFIGURED
      </div>
    );
  }

  return (
    <div className="cards">
      {monitors.map((m, i) => (
        <div key={i} className="card">
          <Corners />
          <div className="card-label">MONITOR</div>
          <div className="card-value" style={{ fontSize: "14px", wordBreak: "break-all" }}>{m.Url}</div>
          <div style={{ marginTop: "8px", fontSize: "10px", color: "var(--text-muted)", letterSpacing: "2px" }}>
            INTERVAL: {m.Time_interval}s &nbsp;|&nbsp; {m.Is_active ? "ACTIVE" : "INACTIVE"}
          </div>
        </div>
      ))}
    </div>
  );
}

function Dashboard({ username, onLogout }) {
  const [time, setTime] = useState(new Date());
  const [monitors, setMonitors] = useState([]);
  const [loadingMonitors, setLoadingMonitors] = useState(true);

  const fetchMonitors = async () => {
    setLoadingMonitors(true);
    try {
      const data = await apiGetMonitors();
      setMonitors(data.monitors || []);
    } catch (e) {
      setMonitors([]);
    }
    setLoadingMonitors(false);
  };

  useEffect(() => {
    const t = setInterval(() => setTime(new Date()), 1000);
    fetchMonitors();
    return () => clearInterval(t);
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    onLogout();
  };

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
          <button className="btn-logout" onClick={handleLogout}>DISCONNECT</button>
        </div>
      </div>

      <div className="dash-content">
        <div className="dash-title">SYSTEM ONLINE</div>
        <div className="dash-subtitle">// MONITORING DASHBOARD</div>

        <AddMonitorForm onAdded={fetchMonitors} />

        <div className="dash-subtitle" style={{ marginBottom: "16px" }}>// ACTIVE MONITORS</div>
        {loadingMonitors
          ? <div style={{ color: "var(--text-muted)", fontSize: "12px", letterSpacing: "2px" }}>LOADING...</div>
          : <MonitorList monitors={monitors} />
        }
      </div>
    </div>
  );
}

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