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
)

const VERSION = "goterm - anongo OS v4.0"
const WELCOME_MSG = "Welcome anon"
const BYE_MSG = "Bey anon"

func banner() {
        fmt.Println("")
        fmt.Println("\033[32m goterm \033[0m - anongo OS v4.0")
        fmt.Println("\033[36m" + WELCOME_MSG + "\033[0m")
        fmt.Println("--------------------------------------------------")
        fmt.Println(VERSION)
        fmt.Println("Public repo: github.com/mrramy01-ship-it/-LINUX")
        fmt.Println("Source: Linux itself (kernel.org, debian, alpine, termux)")
        fmt.Println("Type 'help' for commands | 'exit' to quit")
        fmt.Println("--------------------------------------------------")
}

func showSources() {
        data, err := os.ReadFile("sources.list")
        if err != nil {
                fmt.Println("sources.list not found, using defaults:")
                fmt.Println(" - https://packages.termux.dev")
                fmt.Println(" - http://deb.debian.org")
                fmt.Println(" - https://dl-cdn.alpinelinux.org")
                fmt.Println(" - https://kernel.org")
                return
        }
        fmt.Println("Linux sources:")
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
        io.Copy(out, resp.Body)
        fmt.Printf("[+] Saved as %s\n", filename)
}

func cmdBuildISO() {
        fmt.Println("[*] Building AnonGo-OS disk...")
        os.MkdirAll("anongo-dist/boot", 0755)
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
        os.WriteFile("anongo-dist/VERSION", []byte(VERSION+"\n"), 0644)
        paths := []string{
                "/data/data/com.termux/files/home/storage/downloads/AnonGo-OS.iso",
                "/storage/emulated/0/Download/AnonGo-OS.iso",
                "AnonGo-OS.iso",
        }
        for _, p := range paths {
                cmd := exec.Command("tar", "-czf", p, "-C", "anongo-dist", ".")
                if err := cmd.Run(); err == nil {
                        fmt.Printf("\033[32m[+] Disk created: %s\033[0m\n", p)
                        return
                }
        }
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
  help       - This help
  ls         - List files
  cd <dir>   - Change dir
  pwd        - Current dir
  cat <file> - Read file
  mkdir <dir>- Make dir
  rm <file>  - Remove file
  sources    - Show Linux sources
  wget <url> - Download file
  build-iso  - Build ISO
  clear      - Clear screen
  whoami     - anongo
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
                                        fmt.Printf("anongo: %s: command not found\n", cmd)
                                }
                        }
                }
        }
}
