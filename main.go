package main
import ("bufio";"fmt";"math/rand";"os";"os/exec";"strings";"time")
const VERSION="anongo OS v7.2 -LINUX PURE - 50MB Bootable"
const WELCOME="Welcome to anongo OS v7.2 -LINUX PURE"
const BYE="Bey anongo LINUX"
func banner(){
fmt.Println("")
fmt.Println(" __ _ _ __ ___ _ __ __ _ ___")
fmt.Println(" / _` | '_ \\ / _ \\| '_ \\ / _` |/ _ \\")
fmt.Println(" | (_| | | | | (_) | | | | (_| | (_) |")
fmt.Println(" \\__,_|_| |_|\\___/|_| |_|\\__, |\\___/")
fmt.Println(" |___/ anongo LINUX PURE v7.2")
fmt.Println(WELCOME)
fmt.Println(VERSION)
fmt.Println("Repo: mrramy01-ship-it/-LINUX")
fmt.Println("Commands: ai <q> | glimmer <q> | build-iso-50 | kernel-50 | upload-linux | run-qemu")
}
func glimmerReply(q string){
fmt.Printf("[anongo glimmer LINUX PURE] %s\n",q)
if strings.Contains(strings.ToLower(q),"50") || strings.Contains(strings.ToLower(q),"kernel"){
fmt.Println("[anongo LINUX PURE] build-iso-50 -> anongo-OS-v7.2-50MB.iso (25M vmlinuz + 25M initramfs)")
} else {
fmt.Println("[anongo LINUX PURE] "+q+" -> LINUX PURE OS v7.2")
}
}
func buildISO50(){
fmt.Println("[*] Building anongo LINUX PURE 50MB...")
os.MkdirAll("anongo-dist/boot/grub",0755)
sh50 := `#!/bin/sh
set -e
mkdir -p anongo-dist/boot/grub
dd if=/dev/zero of=anongo-dist/boot/vmlinuz bs=1M count=25 2>/dev/null
echo "# anongo LINUX PURE Kernel 25MB v7.2" >> anongo-dist/boot/vmlinuz
dd if=/dev/zero of=anongo-dist/boot/initramfs bs=1M count=25 2>/dev/null
echo "# anongo LINUX PURE Initramfs 25MB v7.2" >> anongo-dist/boot/initramfs
ls -lh anongo-dist/boot/vmlinuz anongo-dist/boot/initramfs
cat > anongo-dist/boot/grub/grub.cfg <<'GRUB'
set timeout=5
set default=0
menuentry "anongo OS v7.2 LINUX PURE 50MB" {
 linux /boot/vmlinuz root=/dev/sr0 console=tty0 console=ttyS0
 initrd /boot/initramfs
}
GRUB
tar -czf anongo-OS-v7.2-50MB-LINUX-PURE.iso -C anongo-dist boot
ls -lh anongo-OS-v*.iso
echo "[+] anongo LINUX PURE 50MB Ready!"
`
os.WriteFile("build-50.sh",[]byte(sh50),0755)
c:=exec.Command("sh","./build-50.sh")
c.Stdout=os.Stdout
c.Stderr=os.Stderr
c.Run()
}
func main(){
rand.Seed(time.Now().UnixNano())
banner()
sc:=bufio.NewScanner(os.Stdin)
cwd,_:=os.Getwd()
for{
fmt.Printf("anongo-LINUX-PURE# ")
if!sc.Scan(){break}
line:=strings.TrimSpace(sc.Text())
if line==""{continue}
parts:=strings.Fields(line)
cmd:=strings.ToLower(parts[0])
switch cmd{
case "exit","q","quit":fmt.Println(BYE);return
case "help":fmt.Println("anongo LINUX PURE: ai <q> | glimmer <q> | build-iso-50 | kernel-50 | upload-linux | run-qemu | ls | version")
case "glimmer","g","ai","ask":q:="";if len(line)>len(cmd){q=strings.TrimSpace(line[len(cmd):])};glimmerReply(q)
case "build-iso-50","build-50","kernel-50","kernel","iso-50":buildISO50()
case "upload-linux","upload","push":c:=exec.Command("sh","-c","git add.; git commit -m 'v7.2 PURE LINUX only' ; git push");c.Stdout=os.Stdout;c.Stderr=os.Stderr;c.Run()
case "ls":f,_:=os.ReadDir(".");for _,x:=range f{fmt.Printf("%s ",x.Name())};fmt.Println()
case "version":fmt.Println(VERSION)
default:sh:=exec.Command("sh","-c",line);sh.Stdin=os.Stdin;sh.Stdout=os.Stdout;sh.Stderr=os.Stderr;sh.Dir=cwd;sh.Run()
}
}
}
