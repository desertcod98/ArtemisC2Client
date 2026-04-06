# ArtemisC2Client

> ⚠️ **Proof of Concept Project**
>
> ArtemisC2 is primarily a **learning and research project**. It is a **proof of concept (PoC)** designed to explore DNS-based C2 communication with a modular design in mind.
>
> It is **not a production-ready C2 framework**. While it demonstrates core concepts, it **lacks many features expected in mature tools**, such as evasion (the agent has no obfuscation of strings and messages exchanged with the server), many features, and encryption of messages.
>
> The project should be considered a **modular and extensible base**, intended for experimentation, study, and further development—not for real-world operational use.

ArtemisC2 is a Command and Control (C2) that targets Windows systems and that uses DNS as its communication channel with agents. The project is designed for red teaming, research, and educational purposes, providing an interactive web interface for managing agents, jobs, and results.

*This project is the C2 client (implant).*  

The server component can be found in the [ArtemisC2Server](https://github.com/desertcod98/ArtemisC2Server) repo.

## State of the project

The project is stale as I am busy, but it may occasionally see some commits to extend it.

## Main features

- **DNS-based C2 communication**: all communication with the server is performed via DNS TXT queries encoded with [base64url](https://base64.guru/standards/base64url).
- **Remote command execution**, available commands as of now:
  - **whoami**  
    - *Args*: none  
    - *Returns*: `DOMAIN\Username` or `Username`

  - **shell**  
    - *Args*: `command` (string)  
    - *Returns*: stdout + stderr as string  

  - **download**  
    - *Args*: `filepath` (string)  
    - *Returns*: file content (chunked, base64 over DNS)  

  - **setbeaconinterval**  
    - *Args*: `interval` (int, seconds)  
    - *Returns*: confirmation string  
- **Beaconing system**: periodically contacts the server to fetch jobs and send results.
- **Chunked data transfer**: large outputs are split into chunks and reliably transmitted using a TCP-like protocol over DNS.
- **Persistence mechanisms (Windows)**:
  - Registry RunKey.
  - WMI Event Subscription. <span style="color:red">(It MIGHT be bugged, I do not remember the last time I tested)</span>
- **Configuration management**: stores persistent configuration (AgentId, beacon interval, etc.) in `%APPDATA%\ArtemisC2\cfg`.
- **Single instance enforcement**: prevents multiple instances via system mutex.
- **Extensible architecture**: modular dispatcher system for easily adding new commands.

## Technologies used

- Go (>= 1.25.5)
- miekg/dns (DNS communication) <span style="color:red">(It is very important to know that I used this library to force the client to make the DNS queries to localhost for development, in production no library is to be used so that queries are sent through the system DNS resolver)</span>
- go-ole/go-ole (WMI integration)
- Native Windows APIs (persistence and mutex handling)

## Getting started

- Read [Build](#build) for build instructions <span style="color:red">(including how to safely run in debug mode without enabling persistence)</span>.
- When running the executable it will run without a terminal window.

### Prerequisites

- Go >= 1.25.5  
- As of now works  on windows systems, persistence only supported there and there is no system as of now to opt out of it.

## Build

This project supports two build modes:

- **Debug build** for development  
- **Release build**

### Debug build

Enables verbose logging and **disables persistence mechanisms**, meaning:

- No Registry RunKey is created  
- No WMI Event Subscription is installed  
- No permanent artifacts are written to the system  

This prevents the agent from persisting on your machine during testing.

Build manually:

```sh
go build -tags=debug ./cmd/client
```

### Release build

Builds the client without debug protections and with optimizations:

- Strips debug symbols (`-s -w`)  
- Removes build metadata  
- Uses `windowsgui` subsystem (no console window)  
- Enables full behavior, including persistence mechanisms  

```sh
go build -trimpath -ldflags="-s -w -buildid= -H=windowsgui" -buildvcs=false ./cmd/client
```

### VSCode tasks
If you are using VSCode, tasks for both builds are ready to use.
