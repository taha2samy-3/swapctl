# swapctl

A simple Go tool to manage Linux Swap space. It helps you create a swap file, configure partitions, and tune kernel settings like swappiness and cache pressure.

### One-Line Installation

Run the command for your system architecture:

**For AMD64 (Most VPS and Servers):**
```bash
curl -L https://github.com/taha2samy-3/swapctl/releases/download/v1.0.1/swapctl-linux-amd64 -o swapctl && chmod +x swapctl && sudo ./swapctl
```

**For ARM64 (Raspberry Pi, AWS Graviton):**
```bash
curl -L https://github.com/taha2samy-3/swapctl/releases/download/v1.0.1/swapctl-linux-arm64 -o swapctl && chmod +x swapctl && sudo ./swapctl
```

---

### Features
* Detects available partitions and free space.
* Creates and enables swap files automatically.
* Updates `/etc/fstab` for persistence after reboot.
* Tunes `vm.swappiness` and `vm.vfs_cache_pressure`.
* Safely checks for existing swap files before overwriting.

### Usage
When you run the tool:
1. Select the partition number.
2. Enter the swap size (e.g., `2G` or `4G`).
3. Set the swappiness level (Recommended: `10`).
4. Set the cache pressure (Recommended: `50`).

### Requirements
* Linux OS (Ubuntu, Debian, CentOS, etc.)
* Root/Sudo privileges.
* **Note:** Swap cannot be enabled inside Docker containers or GitHub Codespaces due to kernel restrictions.

### License
Open Source - MIT. Created by [Taha Samy](https://github.com/taha2samy-3).
