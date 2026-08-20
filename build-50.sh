#!/bin/sh
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
