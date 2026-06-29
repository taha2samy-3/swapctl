package config

var (
	FstabPath      = "/etc/fstab"
	SysctlPath     = "/etc/sysctl.conf"
	DefaultSize    = "1G"
	DefaultSwap    = "/swapfile"
	SwappinessPath = "/proc/sys/vm/swappiness"
	CachePressure  = "/proc/sys/vm/vfs_cache_pressure"
)
