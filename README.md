<div>
  <img style="100%" src="https://capsule-render.vercel.app/api?type=waving&height=100&section=footer&reversal=false&text=%F0%9F%9A%80%20Event%20manager%20and%20Certificate%20Issuance&fontSize=20&fontColor=FFFFFF&fontAlign=50&fontAlignY=50&stroke=-&animation=twinkling&descSize=20&descAlign=50&descAlignY=50&textBg=false&color=gradient"  />
</div>

###

<p align="left">A production-ready backend system for managing events and issuing digital certificates.<br>Built with Go, designed with clean, scalability, observability, and DevOps-friendly architecture in mind.<br>This service allows organizations to create and manage events, and distribute certificates automatically with PDF generation and structured logging.</p>

###

<h2 align="left">🏗️ System Architecture</h2>

###

<div align="center">
  <img height="200" src="https://i.imgflip.com/65efzo.gif"  />
</div>

###

<h2 align="left">🗄️ Database</h2>

###

<div align="center">
  <img height="200" src="https://i.imgflip.com/65efzo.gif"  />
</div>

###

<h2 align="left">✨ Features</h2>

###

<p align="left">🎟️ Event & certificate management<br>📄 Automatic PDF certificate generation (Gotenberg)<br>🔐 JWT authentication & authorization<br>⚡ High-performance REST API (Gin)<br>🗄️ PostgreSQL main database<br>🧠 Redis caching layer<br>📊 Centralized logging with ELK stack<br>🐳 Fully dockerized environment<br>📦 Production-ready architecture<br>🧪 Input validation<br>⚙️ Configurable via environment variables<br>🔎 Observability ready</p>

###

<h2 align="left">📦 Docker Hub Image</h2>

###

<p align="left">If you pushed your image to Docker Hub:</p>

###

<h2 align="left">🚀 Run with Docker Compose (Recommended)</h2>

###

<p align="left">1️⃣ Clone project<br><br>git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git<br>cd YOUR_REPO<br><br>2️⃣ Run all services<br><br>docker compose up -d<br><br>This will start:<br>API server<br>PostgreSQL<br>Elasticsearch<br>Kibana<br>Filebeat<br>Gotenberg<br><br>3️⃣ Access services<br>Service	URL<br><br>API	http://localhost:8080<br><br>Kibana	http://localhost:5601<br><br>PostgreSQL	localhost:5432</p>

###

<h2 align="left">⚙️ Configuration</h2>

###

<p align="left">Project uses Viper for configuration.<br><br>Example structure:<br><br>config/<br> ├── config.yaml<br> ├── config-docker.yaml</p>

###

<h2 align="left">📄 API Documentation</h2>

###

<p align="left">Postman collection:<br><br>docs/postman_collection.json<br><br>Import into Postman to test endpoints.</p>

###

<h2 align="left">📊 Logging Architecture</h2>

###

<p align="left">Application logs → Filebeat<br><br>Filebeat → Elasticsearch<br><br>Elasticsearch → Kibana<br><br>Structured logging via:<br><br>Zap<br><br>Zerolog</p>

###

<h2 align="left">🛠️ Run Locally (Without Docker)</h2>

###

<p align="left">go mod tidy<br>go run main.go</p>

###

<h2 align="left">📌 TODO</h2>

###

<p align="left">Kubernetes deployment<br><br> CI/CD pipeline<br><br> Rate limiting<br><br> Email sending<br><br> Admin panel</p>

###

<h2 align="left">⭐ Support</h2>

###

<p align="left">If you like this project, give it a star ⭐ on GitHub.</p>

###
