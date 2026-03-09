import { useState, useEffect, useRef } from "react";
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
      Authorization: `Bearer ${token}`,
      ...options.headers,
    },
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || "REQUEST FAILED");
  return data;
}

const apiGetMonitors = () => apiRequest("/monitors");
const apiAddMonitor = (url, timeInterval) =>
  apiRequest("/addmonitor", {
    method: "POST",
    body: JSON.stringify({ Url: url, Time_interval: timeInterval, Is_active: true }),
  });
const apiDeleteMonitor = (id) =>
  apiRequest("/deletemonitor", { method: "POST", body: JSON.stringify({ Id: id }) });
const apiGetChecks = (monitorId) =>
  apiRequest(`/checks/${monitorId}`, { method: "POST" });

// ─── Tiny glitch effect hook ───────────────────────────────────────────────
function useGlitch() {
  const ref = useRef(null);
  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const go = () => {
      el.classList.add("glitch-active");
      setTimeout(() => el.classList.remove("glitch-active"), 200);
    };
    const t = setInterval(go, 4000 + Math.random() * 6000);
    return () => clearInterval(t);
  }, []);
  return ref;
}

// ─── Auth Panel ────────────────────────────────────────────────────────────
function LoginForm({ onSuccess }) {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const submit = async () => {
    if (!login || !password) return setError("ALL FIELDS REQUIRED");
    setError(""); setLoading(true);
    try { await apiLogin(login, password); onSuccess(login); }
    catch (e) { setError(e.message); }
    setLoading(false);
  };

  return (
    <div className="form-body">
      <Field label="IDENTIFIER" value={login} onChange={setLogin} placeholder="user@domain" />
      <Field label="AUTH KEY" type="password" value={password} onChange={setPassword}
        placeholder="••••••••" onKeyDown={e => e.key === "Enter" && submit()} />
      {error && <div className="form-error">⚠ {error}</div>}
      <button className="btn-primary" onClick={submit} disabled={loading}>
        {loading ? <span className="spinner" /> : "ENTER SYSTEM"}
      </button>
    </div>
  );
}

function RegisterForm({ onRegistered }) {
  const [username, setUsername] = useState("");
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const submit = async () => {
    if (!username || !login || !password) return setError("ALL FIELDS REQUIRED");
    setError(""); setLoading(true);
    try { await apiRegister(username, login, password); onRegistered(); }
    catch (e) { setError(e.message); }
    setLoading(false);
  };

  return (
    <div className="form-body">
      <Field label="CALLSIGN" value={username} onChange={setUsername} placeholder="handle" />
      <Field label="IDENTIFIER" value={login} onChange={setLogin} placeholder="user@domain" />
      <Field label="AUTH KEY" type="password" value={password} onChange={setPassword}
        placeholder="••••••••" onKeyDown={e => e.key === "Enter" && submit()} />
      {error && <div className="form-error">⚠ {error}</div>}
      <button className="btn-primary" onClick={submit} disabled={loading}>
        {loading ? <span className="spinner" /> : "CREATE IDENTITY"}
      </button>
    </div>
  );
}

function Field({ label, value, onChange, placeholder, type = "text", onKeyDown }) {
  return (
    <div className="field">
      <span className="field-label">{label}</span>
      <input className="field-input" type={type} value={value}
        onChange={e => onChange(e.target.value)} placeholder={placeholder}
        onKeyDown={onKeyDown} autoComplete="off" />
    </div>
  );
}

function AuthScreen({ onSuccess }) {
  const [tab, setTab] = useState("login");
  const [ok, setOk] = useState("");
  const glitchRef = useGlitch();

  return (
    <div className="auth-root">
      <div className="auth-bg" />
      <div className="auth-card">
        <div className="auth-header">
          <div className="auth-logo" ref={glitchRef} data-text="NETWATCH">NETWATCH</div>
          <div className="auth-sub">// DISTRIBUTED MONITORING SYSTEM //</div>
        </div>
        <div className="auth-tabs">
          {["login", "register"].map(t => (
            <button key={t} className={`auth-tab ${tab === t ? "active" : ""}`}
              onClick={() => { setTab(t); setOk(""); }}>
              {t === "login" ? "AUTHENTICATE" : "REGISTER"}
            </button>
          ))}
        </div>
        {tab === "login"
          ? <LoginForm onSuccess={onSuccess} />
          : <RegisterForm onRegistered={() => { setOk("IDENTITY CREATED — AUTHENTICATE NOW"); setTab("login"); }} />
        }
        {ok && <div className="form-ok">✓ {ok}</div>}
      </div>
    </div>
  );
}

// ─── Monitor card with checks panel ───────────────────────────────────────
function MonitorCard({ monitor, onDeleted }) {
  const [deleting, setDeleting] = useState(false);
  const [checks, setChecks] = useState(null);
  const [loadingChecks, setLoadingChecks] = useState(false);
  const [expanded, setExpanded] = useState(false);

  const handleDelete = async () => {
    setDeleting(true);
    try { await apiDeleteMonitor(monitor.Id); onDeleted(); }
    catch (e) { console.error(e); setDeleting(false); }
  };

  const toggleChecks = async () => {
    if (expanded) { setExpanded(false); return; }
    setExpanded(true);
    setLoadingChecks(true);
    try {
      const data = await apiGetChecks(monitor.Id);
      setChecks(data.monitors || []);
    } catch { setChecks([]); }
    setLoadingChecks(false);
  };

  const lastCheck = checks && checks.length > 0 ? checks[0] : null;
  const avgPing = checks && checks.length > 0
    ? Math.round(checks.reduce((s, c) => s + c.Responce_time, 0) / checks.length)
    : null;

  return (
    <div className={`monitor-card ${expanded ? "expanded" : ""}`}>
      <div className="mc-top">
        <div className="mc-indicator" style={{ background: monitor.Is_active ? "var(--green)" : "var(--dim)" }} />
        <div className="mc-url">{monitor.Url}</div>
        <div className="mc-badge">{monitor.Is_active ? "LIVE" : "IDLE"}</div>
      </div>

      <div className="mc-meta">
        <span>INTERVAL <b>{monitor.Time_interval}s</b></span>
        {avgPing !== null && <span>AVG <b>{avgPing}ms</b></span>}
        {lastCheck && (
          <span className={lastCheck.Status_ok ? "ok" : "fail"}>
            {lastCheck.Status_ok ? "● ONLINE" : "● OFFLINE"}
          </span>
        )}
      </div>

      <div className="mc-actions">
        <button className="btn-ghost" onClick={toggleChecks}>
          {expanded ? "▲ HIDE HISTORY" : "▼ SHOW HISTORY"}
        </button>
        <button className="btn-danger" onClick={handleDelete} disabled={deleting}>
          {deleting ? "REMOVING…" : "REMOVE"}
        </button>
      </div>

      {expanded && (
        <div className="mc-checks">
          {loadingChecks
            ? <div className="checks-loading">LOADING TELEMETRY…</div>
            : checks && checks.length === 0
              ? <div className="checks-empty">// NO DATA YET — WORKER PENDING</div>
              : checks && (
                <table className="checks-table">
                  <thead>
                    <tr><th>TIME</th><th>PING</th><th>STATUS</th></tr>
                  </thead>
                  <tbody>
                    {checks.map((c, i) => (
                      <tr key={i}>
                        <td>{new Date(c.Checked_at).toLocaleTimeString()}</td>
                        <td>{Math.round(c.Responce_time)}ms</td>
                        <td className={c.Status_ok ? "ok" : "fail"}>
                          {c.Status_ok ? "OK" : "FAIL"}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )
          }
        </div>
      )}
    </div>
  );
}

// ─── Add monitor form ──────────────────────────────────────────────────────
function AddMonitorPanel({ onAdded }) {
  const [url, setUrl] = useState("");
  const [interval, setInterval] = useState(60);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [open, setOpen] = useState(false);

  const submit = async () => {
    if (!url) return setError("URL REQUIRED");
    setError(""); setLoading(true);
    try { await apiAddMonitor(url, interval); setUrl(""); setInterval(60); setOpen(false); onAdded(); }
    catch (e) { setError(e.message); }
    setLoading(false);
  };

  return (
    <div className="add-panel">
      <button className="btn-add" onClick={() => setOpen(o => !o)}>
        {open ? "✕ CANCEL" : "+ ADD MONITOR"}
      </button>
      {open && (
        <div className="add-form">
          <Field label="TARGET URL" value={url} onChange={setUrl} placeholder="https://example.com" />
          <Field label="INTERVAL (SEC)" type="number" value={interval} onChange={v => setInterval(Number(v))} placeholder="60" />
          {error && <div className="form-error">⚠ {error}</div>}
          <button className="btn-primary" onClick={submit} disabled={loading}>
            {loading ? <span className="spinner" /> : "DEPLOY MONITOR"}
          </button>
        </div>
      )}
    </div>
  );
}

// ─── Dashboard ─────────────────────────────────────────────────────────────
function Dashboard({ username, onLogout }) {
  const [monitors, setMonitors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [clock, setClock] = useState(new Date());
  const glitchRef = useGlitch();

  const fetchMonitors = async () => {
    setLoading(true);
    try { const d = await apiGetMonitors(); setMonitors(d.monitors || []); }
    catch { setMonitors([]); }
    setLoading(false);
  };

  useEffect(() => {
    fetchMonitors();
    const t = setInterval(() => setClock(new Date()), 1000);
    return () => clearInterval(t);
  }, []);

  const logout = () => { localStorage.removeItem("token"); onLogout(); };

  const active = monitors.filter(m => m.Is_active).length;

  return (
    <div className="dash-root">
      <div className="dash-bg" />

      <header className="dash-header">
        <div className="dh-left">
          <div className="dh-logo" ref={glitchRef} data-text="NETWATCH">NETWATCH</div>
          <div className="dh-status">
            <span className="pulse" />
            <span>SYSTEMS NOMINAL</span>
          </div>
        </div>
        <div className="dh-right">
          <div className="dh-clock">{clock.toLocaleTimeString()}</div>
          <div className="dh-user">// {username.toUpperCase()}</div>
          <button className="btn-logout" onClick={logout}>DISCONNECT</button>
        </div>
      </header>

      <div className="dash-stats">
        <div className="stat-box">
          <div className="stat-val">{monitors.length}</div>
          <div className="stat-label">TOTAL MONITORS</div>
        </div>
        <div className="stat-box">
          <div className="stat-val" style={{ color: "var(--green)" }}>{active}</div>
          <div className="stat-label">ACTIVE</div>
        </div>
        <div className="stat-box">
          <div className="stat-val" style={{ color: "var(--dim)" }}>{monitors.length - active}</div>
          <div className="stat-label">IDLE</div>
        </div>
      </div>

      <div className="dash-body">
        <div className="section-head">
          <span>// MONITOR NODES</span>
          <AddMonitorPanel onAdded={fetchMonitors} />
        </div>

        {loading
          ? <div className="dash-loading">SCANNING NODES…</div>
          : monitors.length === 0
            ? <div className="dash-empty">// NO MONITORS DEPLOYED</div>
            : <div className="monitor-grid">
                {monitors.map((m, i) => (
                  <MonitorCard key={m.Id || i} monitor={m} onDeleted={fetchMonitors} />
                ))}
              </div>
        }
      </div>
    </div>
  );
}

// ─── Root ──────────────────────────────────────────────────────────────────
export default function App() {
  const [user, setUser] = useState(null);
  return user
    ? <Dashboard username={user} onLogout={() => setUser(null)} />
    : <AuthScreen onSuccess={setUser} />;
}