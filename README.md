# swapctl

A simple Go tool to manage Linux Swap space. It helps you create a swap file, configure partitions, and tune kernel settings like swappiness and cache pressure.

### One-Line Installation

Run the command for your system architecture:

**For AMD64 (Most VPS and Servers):**

```bash
curl -L https://github.com/taha2samy-3/swapctl/releases/download/v1.1.1/swapctl-linux-amd64 -o swapctl && chmod +x swapctl && sudo ./swapctl
```

**For ARM64 (Raspberry Pi, AWS Graviton):**
```bash
curl -L https://github.com/taha2samy-3/swapctl/releases/download/v1.1.1/swapctl-linux-arm64 -o swapctl && chmod +x swapctl && sudo ./swapctl
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


### Workload Tuning Recommendations

| Workload / Service | `vm.swappiness` | `vm.vfs_cache_pressure` | `vm.overcommit_memory` |
| :--- | :---: | :---: | :---: |
| **Redis / In-Memory** | `1` | `100` | `1` |
| **SQL Databases** | `10` | `50` | `0` |
| **Web Servers** | `60` | `100` | `0` |
| **High I/O Storage**| `10` | `100` | `0` |

---

### Understanding Kernel Parameters

#### 1. `vm.swappiness` (0 - 100)
Controls how aggressively the Linux kernel moves memory pages from RAM to the swap space.
*   **Low Value (1-10):** Instructs the kernel to avoid swapping as much as possible. It keeps application data in RAM until memory is critically low. Best for databases to prevent disk I/O latency.
*   **High Value (60-100):** The kernel will swap idle processes out of RAM more frequently to free up space for file system caching. Good for general-purpose desktop or balanced web servers.

#### 2. `vm.vfs_cache_pressure` (0 - 100+)
Controls the tendency of the kernel to reclaim the memory which is used for caching directory and inode objects (metadata about where files are on the disk).
*   **Low Value (e.g., 50):** The kernel retains metadata in RAM longer. This consumes more memory but significantly speeds up file lookups and disk queries. Highly recommended for SQL databases.
*   **High Value (>100):** The kernel aggressively drops metadata from RAM. This frees up memory quickly but slows down file system operations. (Default is 100).

#### 3. `vm.overcommit_memory` (0, 1, or 2)
Defines how the kernel handles large memory allocation requests from applications.
*   **0 (Heuristic - Default):** The kernel estimates if enough memory is available. If the request is absurdly large, it denies it.
*   **1 (Always Overcommit):** The kernel pretends there is always enough memory and grants all requests. This is **mandatory for Redis** because it uses a background process (BGSAVE) that temporarily forks memory. Without this, Redis will crash with Out-Of-Memory errors.
*   **2 (Strict):** The kernel strictly checks RAM + Swap limits and denies any request that exceeds them. Excellent for system stability but can cause strict applications to fail.

---

### Technical Justifications

*   **Redis / In-Memory:** Swap latency destroys Redis performance, so `swappiness` is set to `1`. `overcommit_memory` MUST be `1` to allow background snapshot saving without crashing.
*   **SQL Databases (MySQL/PostgreSQL):** Lowering `vfs_cache_pressure` to `50` ensures fast disk metadata lookups. `swappiness` at `10` prevents the database's internal memory buffers from being paged to the slower disk.
*   **Web Servers (Nginx/Apache):** The default values (`60`, `100`, `0`) provide the best balance between serving active connections and keeping the server responsive.
*   **High I/O Storage:** A low `swappiness` (`10`) ensures core networking processes stay in RAM, preventing transfer drops, while keeping other values at standard defaults to manage large file blocks efficiently.

### License
Open Source - MIT. Created by [Taha Samy](https://github.com/taha2samy-3).
