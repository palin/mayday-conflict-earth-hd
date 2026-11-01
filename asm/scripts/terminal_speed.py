#!/usr/bin/env python3
"""Mayday: Conflict Earth (PL, Mayday.exe 684 032 bytes) - terminal scroll speed fix.

The command terminal (KAMPFTERMINAL) scrolls left/right when the mouse is near the
screen edge. The game moved it once per drawn frame, so on modern PCs (60-144+ frames
per second) it flies from one end to the other. This patch runs the same movement at a
fixed rate - by default 30 steps per second - no matter how fast the PC draws.

Changes: the scroll routine at 0041F2E7 (file 0x1E8E7, 392 bytes) is rewritten.

Usage:
  python terminal_speed.py "D:\\Gry\\Mayday\\Mayday.exe"            (30 steps/s)
  python terminal_speed.py "D:\\Gry\\Mayday\\Mayday.exe" --rate 20  (slower)
  python terminal_speed.py "D:\\Gry\\Mayday\\Mayday.exe" --undo
--rate: 10..60 steps per second. Re-run any time with another value.
"""
import argparse
import sys

SIZE = 684032
OFF = 0x41F2E7 - 0x401000 + 0x600
LEN = 392
NEW = bytes.fromhex("e8abc0060089c22b15004f640083fa210f828c00000081fa84000000730c8105004f640021000000eb05a3004f6400a1ec2349008b0d1424490083f8407d146bc00f83c020c1f80683e81039c87d1f89c1eb1b3d400200007e142d400200006bc00f83c020c1f80639c87e0289c185c9742aa11024490001c883f83c7d05b83c0000003d780500007e05b878050000a31024490085c97f0341eb0149890d14244900c3")
TUNE = (0x0F, 0x18, 0x24)  # step ms (imm8), 4*step (imm32), step (imm32)


def build(rate):
    step = round(1000 / rate)
    b = bytearray(NEW)
    b[0x0F] = step
    b[0x18:0x1C] = (4 * step).to_bytes(4, "little")
    b[0x24:0x28] = step.to_bytes(4, "little")
    return bytes(b) + b"\xcc" * (LEN - len(b))


def is_new(cur):
    a, b = bytearray(cur[:len(NEW)]), bytearray(NEW)
    for p, n in ((0x0F, 1), (0x18, 4), (0x24, 4)):
        a[p:p + n] = b[p:p + n] = bytes(n)
    return a == b


def main():
    ap = argparse.ArgumentParser(description="Mayday terminal scroll speed")
    ap.add_argument("exe")
    ap.add_argument("--rate", type=int, default=30)
    ap.add_argument("--undo", action="store_true")
    a = ap.parse_args()
    data = bytearray(open(a.exe, "rb").read())
    if len(data) != SIZE:
        sys.exit(f"Unexpected size {len(data)} - this patch is only for the Polish Mayday.exe.")
    cur = bytes(data[OFF:OFF + LEN])
    backup = a.exe + ".terminal.bak"
    if a.undo:
        try:
            orig = open(backup, "rb").read()[OFF:OFF + LEN]
        except OSError:
            sys.exit("No Mayday.exe.terminal.bak found - cannot undo.")
        data[OFF:OFF + LEN] = orig
        open(a.exe, "wb").write(data)
        print("Original terminal scrolling restored.")
        return
    if not 10 <= a.rate <= 60:
        sys.exit("--rate must be between 10 and 60")
    if not is_new(cur):
        if cur[:6] != bytes.fromhex("558bec51833d"):
            sys.exit("This Mayday.exe is not the expected version - not touching the file.")
        import shutil, os
        if not os.path.exists(backup):
            shutil.copy2(a.exe, backup)
            print(f"Backup written: {backup}")
    data[OFF:OFF + LEN] = build(a.rate)
    open(a.exe, "wb").write(data)
    print(f"Terminal scroll fixed: {a.rate} steps per second.")


if __name__ == "__main__":
    main()
