#!/usr/bin/env python3
"""Mayday: Conflict Earth (PL, Mayday.exe 684 032 bytes) - box selection fix, v2.

Two changes:
 1. 004319D5 (file 0x30FD5, 5 bytes): the drag box always uses the per-unit
    test (SelectVisibleUnits -> IsObjectInBox) instead of the tile walk.
 2. 0042DB0C (file 0x2D10C, 265 bytes): IsObjectInBox is rewritten. A unit is
    selected as soon as the box TOUCHES a rectangle around the unit, instead of
    only when the unit's ground point is inside the box.

The rectangle is measured from the unit's ground point (where it stands):
    LEFT / RIGHT  pixels to each side
    UP            pixels above the ground point (unit height)
    DOWN          pixels below the ground point
Defaults: left 16, right 16, up 40, down 6. Each must be between -96 and 127.

Usage:
  python selection_fix.py "D:\\Retro\\Mayday\\Mayday.exe"
  python selection_fix.py "D:\\Retro\\Mayday\\Mayday.exe" --up 48 --down 4
  python selection_fix.py "D:\\Retro\\Mayday\\Mayday.exe" --undo
You can run it again with different numbers to re-tune; no need to undo first.
A backup (Mayday.exe.bak) of the original is written before the first change.
"""
import argparse
import hashlib
import os
import shutil
import sys

SIZE = 684032
SHA_ORIGINAL = "9f9ab3f11a25bb2241537e984e1e487eb8e97282dbeb20f6ab0d137f359ebfb4"

JMP_OFF = 0x30FD5
JMP_ORIG = bytes.fromhex("ff7510ff75")
JMP_NEW = bytes.fromhex("e92d010000")  # jmp 00431B07

FUNC_OFF = 0x2D10C
FUNC_LEN = 265
FUNC_ORIG_SHA = "ad3d7dc043934be2eaf59002eea84d2aa3e74539f2620ec7897e01aebbcf8855"
# new IsObjectInBox; the 4 tuning bytes are at 0x42, 0x4C, 0x56, 0x60 (see TUNE)
FUNC_NEW = bytes.fromhex(
    "5589e553568b4508c1e0058b1d8cc0600001c38b4d0c8b5318035135"
    "8b7104c1e6040335307063000fb6760e8b431c29f00341398b731c2b"
    "7304d1fe29f089d683c6303b75107c2589d683c6103b75187f1b89c6"
    "83c6263b75147c1189c683c6f83b751c7f07b801000000eb0231c05e"
    "5b5dc3"
)
TUNE = {"right": 0x42, "left": 0x4C, "down": 0x56, "up": 0x60}
# byte value = R+32, 32-L, D+32, 32-U   (signed 8-bit)


def tune_bytes(left, right, up, down):
    vals = {"right": right + 32, "left": 32 - left, "down": down + 32, "up": 32 - up}
    for k, v in vals.items():
        if not -128 <= v <= 127:
            sys.exit(f"--{k} is out of range")
    return {k: v & 0xFF for k, v in vals.items()}


def build_func(left, right, up, down):
    f = bytearray(FUNC_NEW)
    for k, v in tune_bytes(left, right, up, down).items():
        f[TUNE[k]] = v
    return bytes(f) + b"\xcc" * (FUNC_LEN - len(f))


def is_v2(func):
    a, b = bytearray(func[: len(FUNC_NEW)]), bytearray(FUNC_NEW)
    for i in TUNE.values():
        a[i] = b[i] = 0
    return a == b


def main():
    ap = argparse.ArgumentParser(description="Mayday box selection fix v2")
    ap.add_argument("exe")
    ap.add_argument("--left", type=int, default=16)
    ap.add_argument("--right", type=int, default=16)
    ap.add_argument("--up", type=int, default=40)
    ap.add_argument("--down", type=int, default=6)
    ap.add_argument("--undo", action="store_true")
    a = ap.parse_args()

    data = bytearray(open(a.exe, "rb").read())
    if len(data) != SIZE:
        sys.exit(f"Unexpected size {len(data)} - this patch is only for the Polish Mayday.exe.")
    jmp = bytes(data[JMP_OFF:JMP_OFF + 5])
    func = bytes(data[FUNC_OFF:FUNC_OFF + FUNC_LEN])
    func_orig = hashlib.sha256(func).hexdigest() == FUNC_ORIG_SHA
    backup = a.exe + ".bak"

    if a.undo:
        if not os.path.exists(backup):
            sys.exit("No Mayday.exe.bak found - restore the original file manually.")
        if hashlib.sha256(open(backup, "rb").read()).hexdigest() != SHA_ORIGINAL:
            sys.exit("Mayday.exe.bak is not the original version - not restoring.")
        shutil.copy2(backup, a.exe)
        print("Original Mayday.exe restored from backup.")
        return

    if jmp not in (JMP_ORIG, JMP_NEW) or not (func_orig or is_v2(func)):
        sys.exit("This Mayday.exe is not the expected version - not touching the file.")
    if hashlib.sha256(data).hexdigest() == SHA_ORIGINAL and not os.path.exists(backup):
        shutil.copy2(a.exe, backup)
        print(f"Backup written: {backup}")

    data[JMP_OFF:JMP_OFF + 5] = JMP_NEW
    data[FUNC_OFF:FUNC_OFF + FUNC_LEN] = build_func(a.left, a.right, a.up, a.down)
    open(a.exe, "wb").write(data)
    print(f"Selection fix v2 applied (left {a.left}, right {a.right}, up {a.up}, down {a.down}).")


if __name__ == "__main__":
    main()
