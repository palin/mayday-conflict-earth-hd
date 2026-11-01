#!/usr/bin/env python3
"""Mayday: Conflict Earth (PL, Mayday.exe 684 032 bytes) - minimap world outline.

Draws a 1-pixel steel blue line (palette colour 94) around the world diamond on the
minimap, so the "end of the world" is visible. Works with or without the
selection fix.

Changes:
  0047CE79 (file 0x7C479, 7 bytes): end of the minimap build function jumps to the new code
  0048B840 (file 0x8AE40, 106 bytes): new code, placed in unused padding of the code section

Usage:
  python minimap_outline.py "D:\\Retro\\Mayday\\Mayday.exe"
  python minimap_outline.py "D:\\Retro\\Mayday\\Mayday.exe" --colour 40     (other palette colour)
  python minimap_outline.py "D:\\Retro\\Mayday\\Mayday.exe" --undo
"""
import argparse
import sys

SIZE = 684032
HOOK_OFF = 0x7C479
HOOK_ORIG = bytes.fromhex("5f5e5b8be55dc3")
HOOK_NEW = bytes.fromhex("e9c2e900009090")      # jmp 0048B840 ; nop ; nop
CAVE_OFF = 0x8AE40
CAVE = bytes.fromhex("608b351c70630085f674578b0d2470630085c97e4d8d1c0931d239da7d4489d70faffb01f731c039d87d34803c07ff752b85c0740e8a4c07ff80f9ff740580f95e75158d480139d97d128a4c070180f9ff740980f95e7404c604075e40ebc842ebb8615f5e5b89ec5dc3")
COLOUR_POS = (0x40, 0x55, 0x5B)                 # the three places that hold the colour (94)


def main():
    ap = argparse.ArgumentParser(description="Mayday minimap outline")
    ap.add_argument("exe")
    ap.add_argument("--colour", type=int, default=94)
    ap.add_argument("--undo", action="store_true")
    a = ap.parse_args()
    if not 0 <= a.colour <= 254:
        sys.exit("--colour must be 0-254 (255 is transparent).")

    data = bytearray(open(a.exe, "rb").read())
    if len(data) != SIZE:
        sys.exit(f"Unexpected size {len(data)} - this patch is only for the Polish Mayday.exe.")
    hook = bytes(data[HOOK_OFF:HOOK_OFF + 7])
    cave = bytes(data[CAVE_OFF:CAVE_OFF + len(CAVE)])

    if a.undo:
        if hook == HOOK_ORIG:
            sys.exit("Outline is not installed - nothing to undo.")
        if hook != HOOK_NEW:
            sys.exit("Unexpected bytes - not touching the file.")
        data[HOOK_OFF:HOOK_OFF + 7] = HOOK_ORIG
        data[CAVE_OFF:CAVE_OFF + len(CAVE)] = bytes(len(CAVE))
        open(a.exe, "wb").write(data)
        print("Minimap outline removed.")
        return

    installed = hook == HOOK_NEW
    if not installed and (hook != HOOK_ORIG or any(cave)):
        sys.exit("This Mayday.exe is not the expected version - not touching the file.")
    new = bytearray(CAVE)
    for p in COLOUR_POS:
        assert new[p] == 94
        new[p] = a.colour
    data[HOOK_OFF:HOOK_OFF + 7] = HOOK_NEW
    data[CAVE_OFF:CAVE_OFF + len(CAVE)] = new
    open(a.exe, "wb").write(data)
    print(f"Minimap outline {'updated' if installed else 'installed'} (colour {a.colour}).")


if __name__ == "__main__":
    main()
