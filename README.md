# 🤖 NOSTR-AGENT

A scalable, extensible bot system built in **Go (Golang)** for handling direct messages (DMs), group chats, and automated event handling on the **Nostr** network.

### 🌍 Features

- ✅ **Direct Message Bot** (`support_bot`) – Handles user queries and support via encrypted direct messages.
- ✅ **Group Chat Bot** (`weather_bot`) – Broadcasts weather updates to group channels.
- ✅ **Welcome Bot** (`welcome_bot`) – Greets new users in group channels.
- ✅ **Session Manager Bot** (`session`) – Manages multi-bot sessions (Yin and Yang) for interaction.
- ✅ **Extensible Plugin System** – Supports global and handler-specific plugins (e.g., logging, notifications).
- ✅ **Dockerized Deployment** – Easily deploy with Docker, including auto-restart capabilities.
- ✅ **Auto-Resilience** – Automatically restarts if the bot crashes or encounters errors.
- ✅ **Dynamic Configuration** – Load bot configurations via YAML files.

---

### 🚀 **Getting Started**

#### 📦 **Prerequisites**

- [Go 1.21+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

---

#### 🔑 **Environment Variables**

Create a `.env` file in the root directory with the following values:

```env
METEOMATICS_USERNAME=
METEOMATICS_PASSWORD=
```

---

### ⚙️ **Configuration**

Create YAML configuration files inside the `configs/` directory.

#### Example: `configs/support_bot.yaml`
```yaml
bots:
  - name: "Support Bot"
    relay_url: "wss://relay.example.com"
    nsec: "your-secret-key"
    channel_id: "your-channel-id"
    listener: "DMListener"
    publisher: "DMPublisher"
    handler: "SupportHandler"
    event_type: "DMResponseEvent"
```

#### Example: `configs/weather_bot.yaml`
```yaml
bots:
  - name: "Weather Bot"
    relay_url: "wss://relay.example.com"
    nsec: "your-secret-key"
    channel_id: "your-channel-id"
    listener: "GroupListener"
    publisher: "GroupPublisher"
    handler: "GroupHandler"
    event_type: "GroupResponseEvent"
```

#### Example: `configs/welcome_bot.yaml`
```yaml
bots:
  - name: "Welcome Bot"
    relay_url: "wss://relay.example.com"
    nsec: "your-secret-key"
    channel_id: "your-channel-id"
    listener: "DMListener"
    publisher: "GroupPublisher"
    handler: "WelcomeHandler"
    event_type: "GroupResponseEvent"
```

#### Example: `configs/session.yaml`
```yaml
bots:
  - name: "Yin Bot"
    relay_url: "wss://relay.example.com"
    nsec: "your-secret-key"
    channel_id: "your-channel-id"
    listener: "GroupListener"
    publisher: "GroupPublisher"
    handler: "ExchangeHandler"
    event_type: "GroupResponseEvent"

  - name: "Yang Bot"
    relay_url: "wss://relay.example.com"
    nsec: "your-secret-key"
    channel_id: "your-channel-id"
    listener: "GroupListener"
    publisher: "GroupPublisher"
    handler: "ExchangeHandler"
    event_type: "GroupResponseEvent"
```

---

### 🔨 **Building and Running**

#### ✅ 1. Running Locally (Without Docker)

**Install dependencies**:
```bash
go mod download
```

**Run the bot**:
```bash
go run main.go --config=configs/support_bot.yaml
```

#### 🐳 2. Running with Docker

**Step 1**: Build the Docker image:
```bash
docker-compose build
```

**Step 2**: Start the bot using Docker Compose:
```bash
docker-compose up -d
```

**Step 3**: View logs:
```bash
docker logs -f support-bot
```

**Step 4**: Stop the bot:
```bash
docker-compose down
```

---

### 🔥 **Commands**

| Command                           | Description                          |
|-----------------------------------|--------------------------------------|
| `agent --config=configs/support_bot.yaml` | Starts support bot                    |
| `agent --config=configs/weather_bot.yaml` | Starts weather bot                    |
| `agent --config=configs/welcome_bot.yaml` | Starts welcome bot                    |
| `agent --config=configs/session.yaml`    | Starts session with Yin and Yang bots |

**Generate project tree**:
```bash
tree --prune -I "$(paste -sd'|' .treeignore)" > tree.txt
```

---

### 🐳 **Docker Compose Configuration Example**

```yaml
version: "3.8"

services:
  support_bot:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: support-bot
    restart: unless-stopped
    env_file:
      - .env
    command: ["--config=configs/support_bot.yaml"]
    volumes:
      - ./logs/support:/app/logs

  weather_bot:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: weather-bot
    restart: unless-stopped
    env_file:
      - .env
    command: ["--config=configs/weather_bot.yaml"]
    volumes:
      - ./logs/weather:/app/logs

  welcome_bot:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: welcome-bot
    restart: unless-stopped
    env_file:
      - .env
    command: ["--config=configs/welcome_bot.yaml"]
    volumes:
      - ./logs/welcome:/app/logs
```

---

### 🤝 **Contributing**

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Create a new Pull Request.

---

### 📝 **License**

This project is licensed under the MIT License.

---

### 📫 **Contact**

For support or collaboration inquiries, reach out to:
- GitHub Issues
- Email: prorobot.ai.sales@gmail.com

---

### 🌟 **Acknowledgments**

- Built using **Go**
- Uses the **Nostr** protocol
- Dockerized for easy deployment 🚀