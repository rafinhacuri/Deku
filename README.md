<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=venom&height=300&color=gradient&text=Deku"/>
</p>

---

## 💡 Why Choose Deku

> *A next‑generation personal finance dashboard that combines clarity, intelligence, and a beautiful user experience.*

Deku is designed for anyone who wants **full visibility and control** over their financial life — from tracking income sources to understanding spending habits with precision.

---

### 🚀 Key Highlights

#### 📊 Unified Finance Overview
Manage **all incomes, expenses, categories, and payment sources** in a single elegant interface, with real‑time monthly analysis.

#### 🔒 Secure & Reliable Architecture
- JWT‑based authentication  
- AES‑256 encryption support for sensitive data  
- Designed for privacy and safe local computation  

#### 💻 Built with Modern Web Tech
Crafted using **Nuxt 4** and **Nuxt UI 4**, ensuring a smooth, responsive, and visually appealing experience with dark‑mode support.

#### 🤖 Smart Financial Insights
Automatically analyzes:
- Spending per category  
- Remaining balance per income source  
- Month‑to‑month variations  
- Expense distribution and patterns  

#### 🌐 Flexible & Extensible
Perfect for personal use or as a foundation for more advanced budgeting systems.

---

## 🔍 Technical Overview

### ⚙️ Backend
- **Gin (Go)** — ultra‑fast REST API framework  
- **SQLite** — lightweight and efficient local storage  
- Includes: income management, expense tracking, analytics, JWT auth, pagination, and monthly filters  

### 🖥️ Frontend
- **Nuxt 4 + Nuxt UI 4** — next‑gen Vue ecosystem  
- Beautiful card‑based UI with gradient highlights  
- Interactive charts and animated transitions  

### 🧠 Finance Architecture
- Central analysis layer for income and expenses  
- Automatic calculation of:
  - Remaining balance per income source  
  - Category totals  
  - Percentage distribution  
  - Monthly summaries  

---

## 🧭 Application Pages Overview

Deku provides a clean, modular interface — each page designed for clarity and simplicity in financial tracking.

### 💼 **Income**
Manage all your income sources:
- Add, edit, and delete income entries  
- Categorized types (salary, freelance, investments, benefits, etc.)  
- Per‑source summary including **remaining balance** and **total spent**  
- Visual proportional bars  

<p align="center">
  <img src="https://github.com/rafinhacuri/deku/blob/main/public/system.png"/>
</p>

### 🛒 **Expenses**
Full control over spending:
- Category‑based tracking (food, shopping, bills, entertainment, etc.)
- Each category includes:
  - Total spent  
  - Percentage of the month  
  - Colored visual identity  
  - Breakdown per income source  
- Detailed history with pagination, edit, and delete actions  

### 📈 **Dashboard**
Real‑time overview of your financial health:
- Total income, total expenses, remaining monthly balance  
- Category distribution  
- Insights highlighting major spending  

### ⚙️ **Configuration**
Adjust system behavior:
- Language selection (i18n ready: English & Portuguese)  
- Theme preferences  
- API base configuration  

### 👥 **Users**
(Administrator only — if enabled)  
Manage system users:
- Add or remove users  
- Control authentication and access  

---

## ✅ Quick Start Checklist

### 🔧 Step 1 — Install Prerequisites
- [ ] Install **Docker**  
- [ ] Install **Docker Compose**  
👉 Official guide: https://docs.docker.com/get-started/get-docker/

---

### 📦 Step 2 — Get the `docker-compose.yaml`

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

Or copy it directly from the repository.

---

### 📝 Step 3 — Configure Environment

Place your `.env` file in the project root directory:

```
JWT_SECRET="exmp"
CRYPT_KEY="exp"
```

> **Note:**  
> - `JWT_SECRET` must be **at least 32 characters**.  
> - `CRYPT_KEY` must be a **base64 string that decodes to 32 bytes (AES‑256)**.  
>   
> Generate a valid key:
> ```bash
> openssl rand -base64 32
> ```

---

### 🚀 Step 4 — Launch Services

```bash
docker compose pull
docker compose up -d --force-recreate
```

Frontend and backend automatically boot via Docker.

---

### 🔍 Step 5 — Verify Installation

```bash
docker compose ps
```

If everything shows as `Up`, you’re ready to use **Deku** 🎉  

---

## 🤝 Contribution

Contributions, issues, and feature requests are welcome!  
Visit the repository issues page to participate.

---

## 📜 License

> Licensed under the MIT License  
> © 2025 Rafael Curi Leonardo

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)