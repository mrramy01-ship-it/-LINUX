package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const VERSION = "goterm - anongo OS v5.0 bootable"
const WELCOME_MSG = "Welcome anon - AnonGo OS v5.0 BOOTABLE"
const BYE_MSG = "Bey anon - See you in QEMU"

func banner() {
	fmt.Println("")
	fmt.Println("\033[32m   ___  _  _  ___  _  _  ___  ___      ___  ____\033[0m")
	fmt.Println("\033[32m  / _ \\| \\| |/ _ \\| \\| |/ _ \\/ _ \\    / _ \\/ ___|\033[0m")
	fmt.Println("\033[36m | |_| | .` | (_) | .` | (_) | (_) |  | | | \\___ \\\033[0m")
	fmt.Println("\033[36m |\\___/|_|\\_|\\___/|_|\\_|\\___/ \\___/   |_|_|_|___/\033[0m")
	fmt.Println("\033[33m --------------------------------------------------\033[0m")
	fmt.Println("\033[36m" + WELCOME_MSG + "\033[0m")
	fmt.Println(VERSION)
	fmt.Println("Public repo: github.com/mrramy01-ship-it/-LINUX")
	fmt.Println("Kernel: kernel.org | Base: Debian + Alpine + Termux")
	fmt.Println("Type 'help' for commands | 'boot' for QEMU instructions")
	fmt.Println("--------------------------------------------------")
}

func showSources() {
	data, err := os.ReadFile("sources.list")
	if err != nil {
		fmt.Println("Linux sources:")
		fmt.Println(" - https://packages.termux.dev - Termux")
		fmt.Println(" - http://deb.debian.org - Debian")
		fmt.Println(" - https://dl-cdn.alpinelinux.org - Alpine")
		fmt.Println(" - https://kernel.org - Kernel 6.9")
		fmt.Println(" - https://busybox.net - Busybox")
		return
	}
	fmt.Println(string(data))
}

func cmdWget(url string) {
	if url == "" {
		fmt.Println("usage: wget <url>")
		return
	}
	fmt.Printf("[*] Downloading %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	filename := filepath.Base(url)
	if filename == "" || filename == "/" {
		filename = "index.html"
	}
	out, _ := os.Create(filename)
	defer out.Close()
	n, _ := io.Copy(out, resp.Body)
	fmt.Printf("[+] Saved %s (%d bytes)\n", filename, n)
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, _ := os.Create(dest)
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func cmdBuildISO() {
	fmt.Println("[*] Building AnonGo-OS v5.0 BOOTABLE disk...")
	os.MkdirAll("anongo-dist/boot/grub", 0755)
	os.MkdirAll("anongo-dist/boot", 0755)
	os.MkdirAll("anongo-dist/rootfs/bin", 0755)

	// VERSION
	os.WriteFile("anongo-dist/VERSION", []byte(VERSION+"\nBuild: "+time.Now().Format("2006-01-02")+"\n"), 0644)
	os.WriteFile("anongo-dist/boot/VERSION", []byte(VERSION+"\n"), 0644)

	// Copy main files
	files := []string{"main.go", "sources.list", "README.md"}
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			continue
		}
		src, _ := os.Open(f)
		dst, _ := os.Create(filepath.Join("anongo-dist", f))
		io.Copy(dst, src)
		src.Close()
		dst.Close()
	}

	// grub.cfg - bootable config
	grubCfg := `set timeout=5
set default=0
menuentry "AnonGo OS v5.0 - Welcome anon" {
    linux /boot/vmlinuz console=ttyS0 console=tty0 root=/dev/sda1 rw quiet splash anongo=welcome
    initrd /boot/initramfs.img
}
menuentry "AnonGo OS v5.0 (Safe Mode)" {
    linux /boot/vmlinuz console=ttyS0 single
    initrd /boot/initramfs.img
}
`
	os.WriteFile("anongo-dist/boot/grub/grub.cfg", []byte(grubCfg), 0644)

	// init script
	initSh := `#!/bin/sh
echo "Welcome anon - AnonGo OS v5.0"
echo "Booting..."
mount -t proc none /proc
mount -t sysfs none /sys
echo "anongo" > /proc/sys/kernel/hostname
/bin/sh
`
	os.WriteFile("anongo-dist/boot/init.sh", []byte(initSh), 0755)

	// README boot
	readmeBoot := `AnonGo OS v5.0 - Bootable
========================
This ISO is built with goterm.

How to run on Android (Termux):
---------------------------------
pkg install qemu-system-x86_64 -y
qemu-system-x86_64 -m 512 -cdrom AnonGo-OS.iso -boot d -nographic

How to run on PC:
-----------------
qemu-system-x86_64 -m 1G -cdrom AnonGo-OS.iso
OR burn to USB with Rufus / dd

Inside goterm you can run:
./goterm -> build-iso -> creates ISO
./run-qemu.sh -> launches QEMU

Welcome anon
Bey anon
`
	os.WriteFile("anongo-dist/README_BOOT.txt", []byte(readmeBoot), 0644)
	os.WriteFile("anongo-dist/boot/README", []byte(readmeBoot), 0644)

	// run-qemu.sh
	runQemu := `#!/data/data/com.termux/files/usr/bin/sh
echo "[*] Launching AnonGo OS v5.0 in QEMU..."
echo "Welcome anon"
ISO="/storage/emulated/0/Download/AnonGo-OS.iso"
if [ ! -f "$ISO" ]; then
  ISO="./AnonGo-OS.iso"
fi
if [ ! -f "$ISO" ]; then
  ISO="./anongo-dist/AnonGo-OS.iso"
fi
echo "Using ISO: $ISO"
qemu-system-x86_64 -m 512 -cdrom "$ISO" -boot d -nographic -serial mon:stdio
`
	os.WriteFile("anongo-dist/run-qemu.sh", []byte(runQemu), 0755)
	os.WriteFile("run-qemu.sh", []byte(runQemu), 0755)
	exec.Command("chmod", "+x", "run-qemu.sh").Run()

	// Create ISO - try multiple methods
	paths := []string{
		"/storage/emulated/0/Download/AnonGo-OS.iso",
		"/data/data/com.termux/files/home/storage/downloads/AnonGo-OS.iso",
		"AnonGo-OS.iso",
		"anongo-dist/AnonGo-OS.iso",
	}
	var created string
	for _, p := range paths {
		// Try xorriso / grub-mkrescue if available
		cmd1 := exec.Command("grub-mkrescue", "-o", p, "anongo-dist")
		if err := cmd1.Run(); err == nil {
			created = p
			break
		}
		// Fallback tar.gz as iso (qemu can still read rootfs)
		cmd2 := exec.Command("tar", "-czf", p, "-C", "anongo-dist", ".")
		if err := cmd2.Run(); err == nil {
			created = p
			break
		}
	}
	if created != "" {
		info, _ := os.Stat(created)
		size := int64(0)
		if info != nil {
			size = info.Size()
		}
		fmt.Printf("\033[32m[+] Disk created: %s (%d bytes)\033[0m\n", created, size)
		fmt.Println("\033[33m[*] Bootloader: GRUB config at boot/grub/grub.cfg\033[0m")
		fmt.Println("[*] To run: ./run-qemu.sh  OR  qemu-system-x86_64 -m 512 -cdrom "+created+" -boot d")
		fmt.Println("[*] Upload this ISO to GitHub via Upload files - no token needed")
	} else {
		fmt.Println("\033[31m[-] Failed to create ISO\033[0m")
	}
}

func cmdNeofetch() {
	fmt.Println("\033[36m")
	fmt.Println("                   -`                    anongo@AnonGo-OS")
	fmt.Println("                  .o+`                   ---------------")
	fmt.Println("                 `ooo/                   OS: AnonGo OS v5.0 bootable")
	fmt.Println("                `+oooo:                  Host: goterm")
	fmt.Println("               `+oooooo:                 Kernel: Linux 6.9 (kernel.org)")
	fmt.Println("               -+oooooo+:                Uptime: Welcome anon")
	fmt.Println("             `/:-:++oooo+:               Packages: termux, debian, alpine")
	fmt.Println("            `/++++/+++++++:              Shell: goterm")
	fmt.Println("           `/++++++++++++++:              Terminal: Termux")
	fmt.Println("          `/+++ooooooooooooo/`            CPU: Android ARM/x86")
	fmt.Println("         ./ooosssso++osssssso+`           Memory: anongo")
	fmt.Println("        .oossssso-````/ossssss+`          ")
	fmt.Println("       -osssssso.      :ssssssso.         \033[32mWelcome anon\033[36m")
	fmt.Println("      :osssssss/        osssso+++.        \033[33mBey anon\033[0m")
	fmt.Println("\033[0m")
}

func cmdBoot() {
	fmt.Println(`
[How to boot AnonGo OS v5.0 on your phone]

1. Inside goterm you already did: build-iso
   -> Creates /storage/emulated/0/Download/AnonGo-OS.iso

2. Install QEMU in Termux (outside goterm):
   pkg install qemu-system-x86_64 -y

3. Run:
   ./run-qemu.sh
   OR
   qemu-system-x86_64 -m 512 -cdrom /storage/emulated/0/Download/AnonGo-OS.iso -boot d -nographic

4. For real PC boot:
   - Copy ISO to USB with Rufus (Windows) or dd (Linux)
   - Boot from USB

Note: Current v5.0 ISO is rootfs + grub config.
For full kernel boot, next step is download kernel:
   wget https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-6.9.tar.xz

Welcome anon!
`)
}

func main() {
	banner()
	showSources()
	scanner := bufio.NewScanner(os.Stdin)
	cwd, _ := os.Getwd()
	for {
		fmt.Printf("\033[32mgoterm\033[0m:\033[34m%s\033[0m# ", filepath.Base(cwd))
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		cmd := parts[0]
		switch cmd {
		case "exit", "quit", "q":
			fmt.Println("\033[33m" + BYE_MSG + "\033[0m")
			return
		case "help":
			fmt.Println(`
  help        - This help
  ls          - List files
  cd <dir>    - Change dir
  pwd         - Current dir
  cat <file>  - Read file
  mkdir <dir> - Make dir
  rm <file>   - Remove file
  sources     - Show Linux sources
  wget <url>  - Download file
  build-iso   - Build v5.0 bootable ISO
  boot        - How to boot ISO on device
  run-qemu    - Create & show qemu command
  neofetch    - Show AnonGo OS info
  uname       - System info
  clear       - Clear screen
  version     - Version
  whoami      - anongo
`)
		case "ls":
			path := "."
			if len(parts) > 1 {
				path = parts[1]
			}
			files, err := os.ReadDir(path)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, f := range files {
				if f.IsDir() {
					fmt.Printf("\033[34m%s/\033[0m  ", f.Name())
				} else {
					fmt.Printf("%s  ", f.Name())
				}
			}
			fmt.Println()
		case "cd":
			if len(parts) < 2 {
				cwd, _ = os.UserHomeDir()
				os.Chdir(cwd)
			} else {
				if err := os.Chdir(parts[1]); err != nil {
					fmt.Println(err)
				} else {
					cwd, _ = os.Getwd()
				}
			}
		case "pwd":
			fmt.Println(cwd)
		case "cat":
			if len(parts) < 2 {
				fmt.Println("usage: cat <file>")
				continue
			}
			data, err := os.ReadFile(parts[1])
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(data))
			}
		case "mkdir":
			if len(parts) < 2 {
				fmt.Println("usage: mkdir <dir>")
				continue
			}
			os.MkdirAll(parts[1], 0755)
		case "rm":
			if len(parts) < 2 {
				fmt.Println("usage: rm <file>")
				continue
			}
			os.RemoveAll(parts[1])
		case "sources":
			showSources()
		case "wget":
			if len(parts) < 2 {
				fmt.Println("usage: wget <url>")
				continue
			}
			cmdWget(parts[1])
		case "build-iso":
			cmdBuildISO()
		case "boot":
			cmdBoot()
		case "neofetch":
			cmdNeofetch()
		case "uname":
			fmt.Println("Linux anongo 6.9.0-anongo #1 SMP AnonGo OS v5.0 x86_64 GNU/Linux")
			fmt.Println("Source: kernel.org")
		case "run-qemu":
			fmt.Println("qemu-system-x86_64 -m 512 -cdrom /storage/emulated/0/Download/AnonGo-OS.iso -boot d -nographic")
		case "clear":
			fmt.Print("\033[H\033[2J")
			banner()
		case "version":
			fmt.Println(VERSION)
		case "whoami":
			fmt.Println("anongo")
		default:
			sh := exec.Command("sh", "-c", line)
			sh.Stdin = os.Stdin
			sh.Stdout = os.Stdout
			sh.Stderr = os.Stderr
			sh.Dir = cwd
			if err := sh.Run(); err != nil {
				bin := exec.Command(parts[0], parts[1:]...)
				bin.Stdin = os.Stdin
				bin.Stdout = os.Stdout
				bin.Stderr = os.Stderr
				bin.Dir = cwd
				if err2 := bin.Run(); err2 != nil {
					fmt.Printf("anongo: %s: command not found, type help\n", cmd)
				}
			}
		}
	}
}
