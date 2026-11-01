#!/bin/sh
# Assemble a patch written in Intel syntax at fixed addresses inside Mayday.exe.
#   tools/asm.sh file.s .section=0xADDR [.section2=0xADDR ...]
# Single-section files use .text:   tools/asm.sh fps.s .text=0x42a336
# Produces file.elf (symbols: nm file.elf) and one file.<section>.bin per section.
set -e
src=$1; shift
base=${src%.s}
args=""
for s in "$@"; do args="$args --section-start=${s}"; done
as --32 "$src" -o "$base.o"
ld -m elf_i386 -e _start $args -o "$base.elf" "$base.o"
for s in "$@"; do
  name=${s%%=*}
  objcopy -O binary -j "$name" "$base.elf" "$base$name.bin"
done
# file offset in Mayday.exe = VA - 0x401000 + 0x600 (CODE section)
# disassemble the result in place:
#   objdump -D -b binary -mi386 -M intel --adjust-vma=$((0x401000-0x600)) \
#     --start-address=0xADDR --stop-address=0xADDR2 Mayday.exe
