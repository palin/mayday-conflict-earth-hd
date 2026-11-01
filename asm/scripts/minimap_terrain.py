#!/usr/bin/env python3
"""Mayday: Conflict Earth (PL, Mayday.exe 684 032 bytes) - brighter minimap terrain.

The minimap shades land by height in 24 levels. Originally the lowest ground is
almost black (palette 0xF5 = RGB 18,16,14), so flat maps (snow plains) look dark.
This patch changes the 24 colour bytes (one 'mov bl, colour' per level in the
minimap build function at 0047CDB6-0047CE13) to a lighter grey scale.
Water (dark blue) and everything else stay unchanged. Works with or without the
other patches.

Usage:
  python minimap_terrain.py "D:\\Retro\\Mayday\\Mayday.exe"                 (default: light)
  python minimap_terrain.py "D:\\Retro\\Mayday\\Mayday.exe" --scale lighter
  python minimap_terrain.py "D:\\Retro\\Mayday\\Mayday.exe" --undo
"""
import argparse
import sys

SIZE = 684032
# file offsets of the colour byte for height levels 0..23 (lowest .. highest)
OFFS = [0x7C3B7 + 4 * i for i in range(24)]
ORIG = [0xF5, 0xF4, 0xF3, 0xF2, 0xF1, 0xF1, 0xF0, 0xF0, 0xEF, 0xEF, 0xEE, 0xEE,
        0xED, 0xED, 0xEB, 0xEB, 0xEC, 0xEC, 0xEA, 0xEA, 0xE9, 0xE8, 0xE7, 0xE6]
# palette greys, dark -> light: F5 18 | F4 27 | F3 37 | F2 55 | F1 72 | F0 89 | EF 106 |
# EE 123 | ED 140 | EB 152 | EC 157 | EA 174 | E9 191 | E8 208 | E7 225 | E6 242
SCALES = {
    "light":   [0xEE] * 3 + [0xED] * 3 + [0xEB] * 3 + [0xEC] * 3 + [0xEA] * 3 + [0xE9] * 3
               + [0xE8] * 2 + [0xE7] * 2 + [0xE6] * 2,               # 123 .. 242
    "lighter": [0xEA] * 4 + [0xE9] * 5 + [0xE8] * 5 + [0xE7] * 5 + [0xE6] * 5,   # 174 .. 242
    "soft":    [0xF0] * 2 + [0xEF] * 2 + [0xEE] * 3 + [0xED] * 3 + [0xEB] * 2 + [0xEC] * 2
               + [0xEA] * 3 + [0xE9] * 3 + [0xE8] * 2 + [0xE6] * 2,  # 89 .. 242
}


def main():
    ap = argparse.ArgumentParser(description="Mayday minimap terrain brightness")
    ap.add_argument("exe")
    ap.add_argument("--scale", choices=SCALES, default="light")
    ap.add_argument("--undo", action="store_true")
    a = ap.parse_args()
    data = bytearray(open(a.exe, "rb").read())
    if len(data) != SIZE:
        sys.exit(f"Unexpected size {len(data)} - this patch is only for the Polish Mayday.exe.")
    for o in OFFS:
        if data[o - 1] != 0xB3 or not 0xE6 <= data[o] <= 0xF5:
            sys.exit("Unexpected bytes - not touching the file.")
    new = ORIG if a.undo else SCALES[a.scale]
    for o, v in zip(OFFS, new):
        data[o] = v
    open(a.exe, "wb").write(data)
    print("Original minimap colours restored." if a.undo else f"Minimap terrain: {a.scale} scale applied.")


if __name__ == "__main__":
    main()
