# Deployment Guide for Production/Web Access

This guide explains how to publish the hotel project so clients can access it from the browser over the internet.

## 1. Recommended production setup
Use the following stack:
- VPS or cloud server (Ubuntu/Debian)
- Docker + Docker Compose
- Nginx as reverse proxy
- HTTPS via Let’s Encrypt
- Domain name such as hotel.example.com

## 2. Server prerequisites
On the production server, install:

```bash
sudo apt update
sudo apt install docker.io docker-compose-plugin nginx certbot python3-certbot-nginx -y
```

## 3. Upload the project
Upload the project to the server, then go to the project root:

```bash
cd /path/to/HotelManagement
```

Create environment file:

```bash
cp .env.example .env
nano .env
```

Use production values like:

```env
APP_PORT=8089
DATABASE_URL=postgres://hotel:hotel@postgres:5432/hotel_management?sslmode=disable
API_TOKEN=ReplaceWithStrongSecretToken
HOTEL_TIMEZONE=Asia/Riyadh
```

Important:
- Keep the database host as `postgres` when running inside Docker Compose.
- Do not use the demo token in production.

## 4. Start the app
```bash
docker compose up --build -d
```

Check app health:

```bash
curl http://localhost:8089/health
```

Expected response:

```json
{"status":"ok"}
```

## 5. Configure Nginx
Create the site config:

```bash
sudo nano /etc/nginx/sites-available/hotel
```

Add:

```nginx
server {
    listen 80;
    server_name hotel.example.com;

    location / {
        proxy_pass http://127.0.0.1:8089;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable it:

```bash
sudo ln -s /etc/nginx/sites-available/hotel /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 6. Enable HTTPS
```bash
sudo certbot --nginx -d hotel.example.com
```

After that, the site can be accessed via:

```text
https://hotel.example.com
```

## 7. Secure production settings
Before handing it to clients:
- Replace demo credentials
- Use a real strong API token
- Put the app behind HTTPS
- Enable backups for PostgreSQL
- Keep Docker images updated
- Restrict access to private admin endpoints with proper auth strategy

## 8. If you only need internal access
If the app is only meant for local office or internal network access, you can skip DNS and HTTPS and open the site directly with:

```text
http://SERVER_IP:8089
```

## 9. Recommended client-facing URL pattern
- Public booking page: https://hotel.example.com/
- Admin dashboard: https://hotel.example.com/admin/

## 10. Useful commands
```bash
# Start application
docker compose up --build -d

# Stop application
docker compose down

# Restart
docker compose restart

# View logs
docker compose logs -f

# Check health
curl http://localhost:8089/health
```

## 11. Final note
This is the standard production model for a web-based hotel booking and management system. It is simple, stable, and suitable for client-facing deployment.
