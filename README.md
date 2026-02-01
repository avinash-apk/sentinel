# Sentinel
> The Universal Command Center for Developers.

Sentinel is a CLI-based "Mission Control" that aggregates notifications from **Discord** and **Slack** into a single terminal dashboard, allowing you to reply instantly without context switching.

## Why Sentinel?
Devs waste hours every day alt-tabbing between communication apps. Sentinel brings the noise to where you already are: **The Terminal.**

* **Unified Feed:** See Discord and Slack messages in one TUI.
* **Instant Reply:** Reply to messages directly from the CLI.
* **ChatOps:** Execute workflows across platforms without leaving your keyboard.

## Tech Stack
* **Language:** Go (Golang)
* **TUI Framework:** Bubble Tea (Charm)
* **APIs:** DiscordGo, Slack-Go
* **Architecture:** Event-Driven (Pub/Sub Bus)

## Quick Start

### 1. Installation
```bash
git clone [https://github.com/YOUR_USERNAME/sentinel.git](https://github.com/YOUR_USERNAME/sentinel.git)
cd sentinel
go mod tidy