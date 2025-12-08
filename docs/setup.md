# ⚙️ Setup Guide

Welcome to **Deku**, your modern, clean, and intelligent personal finance dashboard.  
This guide shows how to install, configure, and run Deku in a simple step‑by‑step flow.

---

## 🚀 Quick Start (Recommended)

Follow these steps to get Deku running in minutes.

---

## 🔧 Step 1 — Install Requirements

Before continuing, install:

- **Docker**
- **Docker Compose**

Official guide:  
https://docs.docker.com/get-started/get-docker/

---

## 📦 Step 2 — Download `docker-compose.yaml`

Use one of the commands below:

<details>
<summary>🔽 Using curl</summary>

```bash
curl -L -o docker-compose.yaml https://raw.githubusercontent.com/rafinhacuri/deku/main/docker-compose.yaml
```
</details>

<details>
<summary>🔽 Using wget</summary>

```bash
wget -O docker-compose.yaml https://raw.githubusercontent.com/rafinhacuri/deku/main/docker-compose.yaml
```
</details>

Or get it directly from the repository.

---

## 📝 Step 3 — Create Your `.env` File

Create a file named `.env` in the project root:

```
JWT_SECRET="example_super_secret_key_here"
CRYPT_KEY="base64_encoded_32_bytes_key"
```

### Important Notes

- `JWT_SECRET` must be **32+ characters**
- `CRYPT_KEY` must decode to **32 bytes (AES‑256)**

Generate a valid key:

```bash
openssl rand -base64 32
```

Paste the generated output into `CRYPT_KEY`.

---

## ▶️ Step 4 — Start Deku

Run:

```bash
docker compose pull
docker compose up -d --force-recreate
```

Deku automatically starts:

- Backend API (Go + Gin + SQLite)
- Frontend UI (Nuxt 4 + Nuxt UI 4)

No extra steps needed.

---

## 🔍 Step 5 — Verify Everything Is Running

```bash
docker compose ps
```

If all services show **Up**, the system is live and ready 🎉

---

## 🧠 What Deku Sets Up for You

Once running, Deku provides:

- Secure JWT authentication  
- Encrypted financial storage (AES‑256 support)  
- Income and expense APIs  
- Category calculations and analytics  
- Monthly summaries and insights  
- A beautiful interactive UI  

All running locally inside your Docker environment.

---

## 🤝 Contributing

Feedback, improvements, and feature ideas are always welcome!  
Visit the issues page on GitHub.

---

## 📜 License

Licensed under the **MIT License**  
© 2025 Rafael Curi Leonardo

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)