# 🧩 Sentinel

**Minimal programmable network intrusion detector — a quiet guardian for your open ports.**

Sentinel is a lightweight, local-first network intrusion detector (NID) designed for developers and DevSecOps engineers.  
It watches your network traffic, alerts on suspicious activity, and optionally blocks unsafe connections — all configured from a single file.

---

## ✨ Features
- 🧱 Lightweight packet inspection via `pcap`/`BPF`
- ⚙️ One-file configuration (`sentinel.toml`)
- 🧩 Preset modes: `dev`, `strict`, `passive`, `honeypot`
- 🔔 Real-time alerts on new listeners or unknown IPs
- 🚫 Optional auto-blocking for denied ports or domains

---

## ⚙️ Usage
```bash
# Start Sentinel with default config
sentinel run

# Or specify a config file
sentinel run --config sentinel.toml
````

Example `sentinel.toml`:

```toml
mode = "strict"

[allow]
domains = ["github.com", "docker.io"]
ports = [22, 443]

[deny]
ports = [23, 3389]

[thresholds]
max_connections_per_minute = 100
alert_on_new_listener = true
```

---

## 🧠 Philosophy

Sentinel focuses on **practical local defense** — no enterprise bloat, no cloud analytics, no hidden data flow.
Just clean, minimal, programmable security for your machine.

---

## 📜 License

MIT — open source and always free.
