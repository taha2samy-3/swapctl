package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/taha2samy/swapctl/internal/config"
	"github.com/taha2samy/swapctl/internal/system"
)

func StartInteractiveSession() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Checking current swap status...")
	current, _ := system.GetActiveSwap()
	if current != "" {
		fmt.Println(current)
	} else {
		fmt.Println("No active swap found.")
	}
	fmt.Println("---------------------------------------")

	// 1. تجميع البيانات (Data Collection)
	mounts, _ := system.GetMountPoints()
	fmt.Println("Available partitions:")
	for i, m := range mounts {
		fmt.Printf("%d) %s (Free: %s)\n", i+1, m.Path, m.Free)
	}

	fmt.Print("Select partition number [1]: ")
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice := 1
	if choiceStr != "" {
		choice, _ = strconv.Atoi(choiceStr)
	}

	selectedMount := mounts[choice-1].Path
	swapPath := config.DefaultSwap
	if selectedMount != "/" {
		swapPath = strings.TrimSuffix(selectedMount, "/") + config.DefaultSwap
	}

	fmt.Print("Enter swap size (e.g., 2 or 2G) [1G]: ")
	size, _ := reader.ReadString('\n')
	size = strings.TrimSpace(size)
	if size == "" {
		size = config.DefaultSize
	}
	if _, err := strconv.Atoi(size); err == nil {
		size = size + "G"
	}

	fmt.Print("Enter swappiness level (0-100) [60]: ")
	swpStr, _ := reader.ReadString('\n')
	swpStr = strings.TrimSpace(swpStr)
	if swpStr == "" {
		swpStr = "60"
	}

	fmt.Print("Enter vfs_cache_pressure (0-100) [50]: ")
	vfsStr, _ := reader.ReadString('\n')
	vfsStr = strings.TrimSpace(vfsStr)
	if vfsStr == "" {
		vfsStr = "50"
	}

	fmt.Println("\n--- Overcommit Memory Modes ---")
	fmt.Println("0: Heuristic (Default)")
	fmt.Println("1: Always Overcommit (Recommended for Redis/Databases)")
	fmt.Println("2: Strict (Don't Overcommit)")
	fmt.Print("Select overcommit mode [0]: ")
	ovcStr, _ := reader.ReadString('\n')
	ovcStr = strings.TrimSpace(ovcStr)
	if ovcStr == "" {
		ovcStr = "0"
	}

	fmt.Println("\n=======================================")
	fmt.Println("   SUMMARY OF PROPOSED CHANGES")
	fmt.Println("=======================================")
	fmt.Printf("Swap Path:      %s\n", swapPath)
	fmt.Printf("Swap Size:      %s\n", size)
	fmt.Printf("Swappiness:     %s\n", swpStr)
	fmt.Printf("Cache Pressure: %s\n", vfsStr)
	fmt.Printf("Overcommit:     %s\n", ovcStr)
	fmt.Println("---------------------------------------")
	fmt.Print("Do you want to apply these changes? (y/n): ")

	confirm, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		fmt.Println("Aborted. No changes were made.")
		return
	}

	fmt.Println("\nStarting execution...")

	if _, err := os.Stat(swapPath); err == nil {
		fmt.Printf("Existing swap found at %s. Removing old file...\n", swapPath)
		system.DisableSwap(swapPath)
		os.Remove(swapPath)
	}

	fmt.Printf("Creating %s swap file at %s...\n", size, swapPath)
	if err := system.CreateSwapFile(swapPath, size); err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}

	system.SetPermissions(swapPath)

	fmt.Println("Formatting swap...")
	if err := system.FormatSwap(swapPath); err != nil {
		fmt.Printf("Error formatting swap: %v\n", err)
		return
	}

	fmt.Println("Enabling swap...")
	if err := system.EnableSwap(swapPath); err != nil {
		fmt.Printf("CRITICAL ERROR: %v\n", err)
		return
	}

	fmt.Println("Updating /etc/fstab for persistence...")
	system.AppendToFstab(config.FstabPath, swapPath)

	fmt.Println("Tuning kernel parameters...")

	swpVal, _ := strconv.Atoi(swpStr)
	system.SetKernelParam(config.SwappinessPath, swpVal)
	system.UpdateSysctl(config.SysctlPath, "vm.swappiness", swpStr)

	vfsVal, _ := strconv.Atoi(vfsStr)
	system.SetKernelParam(config.CachePressure, vfsVal)
	system.UpdateSysctl(config.SysctlPath, "vm.vfs_cache_pressure", vfsStr)

	ovcVal, _ := strconv.Atoi(ovcStr)
	system.SetKernelParam(config.OvercommitPath, ovcVal)
	system.UpdateSysctl(config.SysctlPath, "vm.overcommit_memory", ovcStr)

	fmt.Println("=======================================")
	fmt.Println("FINAL STATUS:")
	final, _ := system.GetActiveSwap()
	fmt.Println(final)
	fmt.Println("Process completed successfully!")
}
