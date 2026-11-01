#!/usr/bin/env python3
"""Mayday: Conflict Earth (PL) - opaque minimap panel.

Edits the "MAPA TAKTYCZNA" panel picture inside mayday.art (164 x 164 pixels,
8-bit, game palette DSpiel.pal) at file offset 0x15CD831. Only the inner area
(rows 33-160, columns 8-159) changes: the half-transparent checkerboard becomes a
solid colour (dark slate). Optionally a rectangular border can be drawn around it
(off by default - the world outline is drawn by the Mayday.exe patch instead).
Nothing else in the file changes.

Usage:
  python minimap_panel.py "D:\\Retro\\Mayday\\mayday.art"                  (solid slate, no border)
  python minimap_panel.py "D:\\Retro\\Mayday\\mayday.art" --fill 83          (other background colour)
  python minimap_panel.py "D:\\Retro\\Mayday\\mayday.art" --style yellow --width 2
  python minimap_panel.py "D:\\Retro\\Mayday\\mayday.art" --undo
Styles (only matter with --width > 0): light, steel, yellow.
--fill / --border are palette indices (0-254) of DSpiel.pal.
A backup (mayday.art.bak) of the original is written before the first change.
"""
import argparse
import hashlib
import os
import shutil
import sys

SIZE = 23185090
OFF = 0x15CD831
W = H = 164
LEN = W * H
ORIG_BLOCK_SHA = "7d6f869fe4a9da66124d24583a0a289a6f85f4c533141c76c25a98447c11c74d"
R0, R1, C0, C1 = 33, 160, 8, 159
STYLES = {"light": (96, 25), "steel": (96, 94), "yellow": (96, 40)}
# 96 = dark slate (42,47,55); 25 = light grey (216,227,230);
# 94 = steel blue-grey (181,190,213); 40 = yellow (255,219,79)


def main():
    ap = argparse.ArgumentParser(description="Mayday minimap panel")
    ap.add_argument("art")
    ap.add_argument("--style", choices=STYLES, default="light")
    ap.add_argument("--fill", type=int)
    ap.add_argument("--border", type=int)
    ap.add_argument("--width", type=int, default=0)
    ap.add_argument("--undo", action="store_true")
    a = ap.parse_args()

    backup = a.art + ".bak"
    if a.undo:
        if not os.path.exists(backup):
            sys.exit("No mayday.art.bak found.")
        shutil.copy2(backup, a.art)
        print("Original mayday.art restored.")
        return

    data = bytearray(open(a.art, "rb").read())
    if len(data) != SIZE:
        sys.exit(f"Unexpected size {len(data)} - this is not the expected mayday.art.")
    block = bytes(data[OFF:OFF + LEN])
    is_orig = hashlib.sha256(block).hexdigest() == ORIG_BLOCK_SHA
    # accept re-runs: the frame outside the inner area must be unchanged
    if not is_orig:
        if not os.path.exists(backup):
            sys.exit("The minimap panel is not the original and there is no backup - not touching the file.")
        orig = open(backup, "rb").read()[OFF:OFF + LEN]
        if hashlib.sha256(orig).hexdigest() != ORIG_BLOCK_SHA:
            sys.exit("Backup does not contain the original panel - not touching the file.")
        block = orig
    elif not os.path.exists(backup):
        shutil.copy2(a.art, backup)
        print(f"Backup written: {backup}")

    fill, border = STYLES[a.style]
    fill = a.fill if a.fill is not None else fill
    border = a.border if a.border is not None else border
    if not (0 <= fill <= 254 and 0 <= border <= 254) or not 0 <= a.width <= 10:
        sys.exit("fill/border must be 0-254 (255 is the transparent colour), width 0-10.")

    px = bytearray(block)
    bw = a.width
    for r in range(R0, R1 + 1):
        for c in range(C0, C1 + 1):
            edge = r < R0 + bw or r > R1 - bw or c < C0 + bw or c > C1 - bw
            px[r * W + c] = border if edge else fill
    data[OFF:OFF + LEN] = px
    open(a.art, "wb").write(data)
    print(f"Minimap panel updated (fill {fill}, border {border}, width {bw}).")


if __name__ == "__main__":
    main()
