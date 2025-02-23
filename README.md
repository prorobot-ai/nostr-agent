# 🤖 NOSTR-AGENT

A scalable, extensible bot system built in **Go (Golang)** for handling direct messages (DMs), group chats, and automated event handling on the **Nostr** network.

### 🌍 Features

- ✅ **Direct Message Bot** (`bot`) – Handles user queries and support via encrypted direct messages.
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
BOT_1_RELAY_URL=wss://relay.example.com
BOT_1_NSEC=your-secret-key
BOT_1_CHANNEL_ID=your-channel-id
```

###  🔨 Building and Running

#### ✅ 1. Running Locally (Without Docker)

Install dependencies
```
go mod download
```

Run the bot
```
go run main.go basic_bot
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
#### 🔌 Plugins System

The bot supports two types of plugins:
1.	Global Plugins – Triggered on every event.
2.	Handler-Specific Plugins – Triggered on specific message events.

#### 📄 Example Plugin Usage:

Logging Plugin (Global):
```
loggingPlugin := &plugins.LoggingPlugin{}
```

Channel Notifier Plugin (Handler-Specific):

```
channelNotifier := &plugins.ChannelNotifierPlugin{
	ChannelID: "channel-123",
}
```

Attach plugins when creating a bot:

```
supportHandler := &handlers.SupportHandler{
	Plugins: []bot.HandlerPlugin{channelNotifier},
}

supportBot := bot.NewBasicBot(
	relayURL,
	nsec,
	supportHandler,
	[]bot.BotPlugin{loggingPlugin},
)
```
#### 🔥 Commands

```
Command	Description
agent basic_bot	- Starts the Bot
```

```
Generate project tree
tree --prune -I "$(paste -sd'|' .treeignore)" > tree.txt
```
#### 🔍 Monitoring and Logs

Monitor agent logs:
```
docker logs -f nostr-agent
```
Restart agent manually:
```
docker-compose restart nostr-agent
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
