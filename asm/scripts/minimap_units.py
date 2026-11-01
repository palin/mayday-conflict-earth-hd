#!/usr/bin/env python3
"""Mayday: Conflict Earth (PL, Mayday.exe 684 032 bytes) - outlined units on the minimap.

Every frame, after the minimap terrain, units, buildings and fog are drawn, each
unit/building pixel gets a 1-pixel ring (all 8 surrounding pixels) in black, so
the coloured dots stand out from the terrain. The ring is never drawn over other
units, the world outline, or outside the world. Units under fog stay hidden.
Works with or without the other patches.

Changes:
  0044EDEE (file 0x4E3EE, 6 bytes): jump to the new code
  0048B8B0 (file 0x8AEB0, 198 bytes): new code, in unused padding of the code section

Usage:
  python minimap_units.py "D:\\Retro\\Mayday\\Mayday.exe"
  python minimap_units.py "D:\\Retro\\Mayday\\Mayday.exe" --colour 97   (other ring colour, palette index)
  python minimap_units.py "D:\\Retro\\Mayday\\Mayday.exe" --undo
"""
import argparse
import sys

SIZE = 684032
HOOK_OFF = 0x4E3EE
HOOK_ORIG = bytes.fromhex("8b1d24706300")      # mov ebx,[00637024]
HOOK_NEW = bytes.fromhex("e9bdca030090")       # jmp 0048B8B0 ; nop
CAVE_OFF = 0x8AEB0
CAVE = bytes.fromhex("608b352070630085f6747b8b0d2470630085c97e718d1c09ba010000008d43ff39c27d6289d70faffb01f7b9010000008d43ff39c17d4c8a040fe870000000753f8d2c0f89e829d848e84400000040e83e00000040e8380000008d45ffe8300000008d4501e82800000089e801d848e81e00000040e81800000040e81200000041ebad42eb97618b1d24706300e9b234fcff528a1080faff741380fa5e740e5088d0e808000000587403c600005ac33c1674123c20740e3c25740a3c2e74063c3574023ce3c3")
COLOUR_POS = 0xAC                              # 'mov byte [eax], colour'


def main():
    ap = argparse.ArgumentParser(description="Mayday minimap unit outlines")
    ap.add_argument("exe")
    ap.add_argument("--colour", type=int, default=0)
    ap.add_argument("--undo", action="store_true")
    a = ap.parse_args()
    if not 0 <= a.colour <= 254 or a.colour in (0x16, 0x20, 0x25, 0x2E, 0x35, 0xE3, 94):
        sys.exit("--colour must be 0-254 and not an object or outline colour.")
    data = bytearray(open(a.exe, "rb").read())
    if len(data) != SIZE:
        sys.exit(f"Unexpected size {len(data)} - this patch is only for the Polish Mayday.exe.")
    hook = bytes(data[HOOK_OFF:HOOK_OFF + 6])
    cave = bytes(data[CAVE_OFF:CAVE_OFF + len(CAVE)])
    if a.undo:
        if hook == HOOK_ORIG:
            sys.exit("Unit outlines are not installed - nothing to undo.")
        if hook != HOOK_NEW:
            sys.exit("Unexpected bytes - not touching the file.")
        data[HOOK_OFF:HOOK_OFF + 6] = HOOK_ORIG
        data[CAVE_OFF:CAVE_OFF + len(CAVE)] = bytes(len(CAVE))
        open(a.exe, "wb").write(data)
        print("Unit outlines removed.")
        return
    installed = hook == HOOK_NEW
    if not installed and (hook != HOOK_ORIG or any(cave)):
        sys.exit("This Mayday.exe is not the expected version - not touching the file.")
    new = bytearray(CAVE)
    assert new[COLOUR_POS - 2:COLOUR_POS] == b"\xc6\x00"
    new[COLOUR_POS] = a.colour
    data[HOOK_OFF:HOOK_OFF + 6] = HOOK_NEW
    data[CAVE_OFF:CAVE_OFF + len(CAVE)] = new
    open(a.exe, "wb").write(data)
    print(f"Unit outlines {'updated' if installed else 'installed'} (colour {a.colour}).")


if __name__ == "__main__":
    main()
