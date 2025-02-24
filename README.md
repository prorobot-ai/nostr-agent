# 🤖 NOSTR-AGENT

A scalable, extensible bot system built in **Go (Golang)** for handling direct messages (DMs), group chats, and automated event handling on the **Nostr** network.

### 🌍 Features

- ✅ **Direct Message Bot** (`support_bot`) – Handles user queries and support via encrypted direct messages.
- ✅ **Extensible Plugin System** – Supports global and handler-specific plugins (e.g., logging, notifications).
- ✅ **Dockerized Deployment** – Easily deploy with Docker, including auto-restart capabilities.
- ✅ **Auto-Resilience** – Automatically restarts if the bot crashes or encounters errors.

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
DISPATCH_RELAY_URL=wss://relay.example.com
DISPATCH_NSEC=your-secret-key
DISPATCH_CHANNEL_ID=your-channel-id
```

###  🔨 Building and Running

#### ✅ 1. Running Locally (Without Docker)

Install dependencies
```
go mod download
```

Run the bot
```
go run main.go welcome_bot
```

#### 🐳 2. Running with Docker

Step 1: Build the Docker image:
```
docker-compose build
```
Step 2: Start the bot using Docker Compose:
```
docker-compose up -d
```
Step 3: View logs:
```
docker logs -f nostr-agent
```
Step 4: Stop the bot:
```
docker-compose down
```
#### 🔥 Commands

```
Command	Description
agent support_bot - Starts support bot
agent weather_bot - Starts weather bot
agent welcome_bot - Starts welcome bot
```

```
Generate project tree
tree --prune -I "$(paste -sd'|' .treeignore)" > tree.txt
```

#### 🤝 Contributing
1.	Fork the repository.
2.	Create a new branch (git checkout -b feature/your-feature).
3.	Commit your changes (git commit -am 'Add new feature').
4.	Push to the branch (git push origin feature/your-feature).
5.	Create a new Pull Request.

#### 📝 License

This project is licensed under the MIT License.

#### 📫 Contact

For support or collaboration inquiries, reach out to:
•	GitHub Issues
•	Email: prorobot.ai.sales@gmail.com

#### 🌟 Acknowledgments
•	Built using Go
•	Uses the Nostr protocol
•	Dockerized for easy deployment 🚀
