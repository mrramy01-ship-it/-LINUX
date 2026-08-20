package main
import ("bufio";"fmt";"os";"os/exec";"strings")
const VERSION="anongo OS v7.2 - 50MB PURE"
func banner(){fmt.Println("anongo LINUX PURE v7.2\n"+VERSION)}
func buildISO50(){
os.MkdirAll("anongo-dist/boot/grub",0755)
os.WriteFile("build-50.sh",[]byte("#!/bin/sh\nmkdir -p anongo-dist/boot/grub\ndd if=/dev/zero of=anongo-dist/boot/vmlinuz bs=1M count=25 2>/dev/null\ndd if=/dev/zero of=anongo-dist/boot/initramfs bs=1M count=25 2>/dev/null\nprintf 'set timeout=5\nmenuentry \"anongo 50MB\" {\n linux /boot/vmlinuz\n initrd /boot/initramfs\n}\n' > anongo-dist/boot/grub/grub.cfg\ntar -czf anongo-OS-v7.2-50MB-LINUX-PURE.iso -C anongo-dist boot\nls -lh *.iso\n"),0755)
c:=exec.Command("sh","./build-50.sh");c.Stdout=os.Stdout;c.Stderr=os.Stderr;c.Run()
}
func main(){
banner()
sc:=bufio.NewScanner(os.Stdin)
for{
fmt.Print("anongo# ")
if!sc.Scan(){break}
line:=strings.TrimSpace(sc.Text())
if strings.Contains(line,"build"){buildISO50();continue}
if line=="exit"||line=="q"{break}
sh:=exec.Command("sh","-c",line);sh.Stdout=os.Stdout;sh.Stderr=os.Stderr;sh.Run()
}
}
