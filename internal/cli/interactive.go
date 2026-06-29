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

	if _, err := os.Stat(swapPath); err == nil {
		fmt.Printf("Existing swap found at %s. Overwrite? (y/n): ", swapPath)
		conf, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(conf)) != "y" {
			fmt.Println("Aborted.")
			return
		}
		system.DisableSwap(swapPath)
		os.Remove(swapPath)
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

	system.AppendToFstab(config.FstabPath, swapPath)

	fmt.Print("Enter swappiness level (0-100) [60]: ")
	swpStr, _ := reader.ReadString('\n')
	swpStr = strings.TrimSpace(swpStr)
	if swpStr == "" {
		swpStr = "60"
	}
	swpVal, _ := strconv.Atoi(swpStr)
	system.SetKernelParam(config.SwappinessPath, swpVal)
	system.UpdateSysctl(config.SysctlPath, "vm.swappiness", swpStr)

	fmt.Print("Enter vfs_cache_pressure (0-100) [50]: ")
	vfsStr, _ := reader.ReadString('\n')
	vfsStr = strings.TrimSpace(vfsStr)
	if vfsStr == "" {
		vfsStr = "50"
	}
	vfsVal, _ := strconv.Atoi(vfsStr)
	system.SetKernelParam(config.CachePressure, vfsVal)
	system.UpdateSysctl(config.SysctlPath, "vm.vfs_cache_pressure", vfsStr)

	fmt.Println("=======================================")
	fmt.Println("FINAL STATUS:")
	final, _ := system.GetActiveSwap()
	fmt.Println(final)
	fmt.Println("Process completed successfully!")
}
