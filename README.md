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

### Variables

| Workload / Service | `vm.swappiness` | `vm.vfs_cache_pressure` | `vm.overcommit_memory` | Technical Justification |
| :--- | :---: | :---: | :---: | :--- |
| **Redis / In-Memory** | `0` or `1` | `100` (Default) | `1` (Always) | **Swappiness**: Redis latency spikes severely if swapped. Docs recommend turning off swap entirely or setting to 0/1.<br>**Overcommit**: MUST be set to 1. Redis uses background saving (BGSAVE) via process forking. Setting to 1 prevents the OS from blocking the fork, avoiding fatal Out-Of-Memory errors. |
| **SQL Databases (MySQL, PostgreSQL)** | `1` to `10` | `50` | `0` or `2` | **Swappiness**: Prevents the OS from paging out the database's internal Buffer Pool/Shared Buffers to disk.<br>**Cache Pressure**: Lowering to 50 retains inode/dentry caches in RAM longer, significantly accelerating disk I/O metadata lookups.<br>**Overcommit**: Setting to 2 (Strict) restricts over-allocation, protecting the DB from being terminated by the Linux OOM Killer. |
| **General Web Servers (Nginx, Apache)**| `60` (Default) | `100` (Default) | `0` (Default) | Relies on default Linux kernel heuristics, providing an efficient, balanced distribution of memory between active web processes and idle connections. |
| **High Traffic Storage / CDN / NAS** | `10` | `200` to `1000` | `0` (Default) | **Swappiness**: Keeps core OS and networking processes in RAM.<br>**Cache Pressure**: Aggressively evicts metadata (directories/inodes) to clear RAM instantly, prioritizing the massive caching of actual file contents (Page Cache) to serve users faster. |

### License
Open Source - MIT. Created by [Taha Samy](https://github.com/taha2samy-3).
