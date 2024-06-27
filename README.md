# pv-reverse-proxy

Simple UDP reverse proxy for [Plasmo Voice 2.1.0](https://github.com/plasmoapp/plasmo-voice/releases/tag/2.1.0-SNAPSHOT) written in Go.

This was written for testing purposes, but maybe someone will find it helpful.

## Prerequisites
- Go 1.22 installed on your machine

## Installation
### Step 1: Build the Proxy
1. Clone the repository:
   ```sh
   git clone https://github.com/Apehum/pv-reverse-proxy.git
   cd pv-reverse-proxy
   ```

2. Build the proxy executable:
   ```sh
   go build -o pv-reverse-proxy pv-reverse-proxy/cmd/proxy
   ```

### Step 2: Create Configuration File
1. Create a file named `servers.toml` in the root directory of the project.
2. Add your server configurations to `servers.toml`:
   ```toml
   [servers]
   "pv-test-domain1.com" = "localhost:25565"
   # "pv-test-domain2.com" = "localhost:25566"
   ```

### Step 3: Configure Plasmo Voice
Configure Plasmo Voice to connect through the reverse proxy. Update your Plasmo Voice configuration file as follows:
   ```toml
   [host]
   ip = "0.0.0.0"
   port = 25565

   [host.public]
   ip = "pv-test-domain1.com"
   port = 30000
   ```

## Usage
Start the proxy with the desired listening port:
   ```sh
   ./pv-reverse-proxy -p <listen port>
   ```
