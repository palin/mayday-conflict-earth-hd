#!/usr/bin/env python3
"""Mayday: Conflict Earth (PL, Mayday.exe 684 032 bytes) - accurate minimap view frame.

1. The orange frame on the minimap gets the real size of the screen view
   (1 minimap pixel = 32 screen pixels across, 16 screen pixels down):
     640x480 -> 20x30   800x600 -> 25x37   1024x768 -> 32x48
     1280x1024 -> 40x64   1920x1080 -> 60x67
   (the original drew a fixed square: 24, 28 or 32 pixels).
2. Clicking on the minimap centres the view on the clicked point
   (the original moved the camera by a fixed 12 tiles, tuned for 640x480).

Changes: 0044EE2F (file 0x4E42F, 183 bytes, the old frame code, rewritten),
0044A704 (file 0x49D04, 44 bytes) and a code cave at 0048B980 (file 0x8AF80, 93 bytes).

Fine tuning (if the frame is off by a pixel or two in your game):
  --dx N   move the frame right by N minimap pixels (negative = left), default 0
  --dy N   move the frame down by N minimap pixels (negative = up),   default 2
Clicking on the minimap follows the same tuning.

Usage:
  python minimap_frame.py "D:\\Retro\\Mayday\\Mayday.exe"
  python minimap_frame.py "D:\\Retro\\Mayday\\Mayday.exe" --dx 1 --dy 1
  python minimap_frame.py "D:\\Retro\\Mayday\\Mayday.exe" --undo
Works on the original exe or one with the earlier patches (selection, minimap,
resolution). Re-run with other --dx/--dy any time. A backup
(Mayday.exe.frame.bak) is written before the first change.
"""
import argparse
import os
import shutil
import sys

SIZE = 684032


def off(va):
    return va - 0x401000 + 0x600


FRAME_VA, FRAME_LEN = 0x44EE2F, 183
FRAME_ORIG = bytes.fromhex(
    "a180c060003d8002000075"
    "07be18000000eb18a180c060003d2003000075"
    "07be1c000000eb05be200000002bfe8bc32bc63bf87c098d56018bfb2bfaeb06"
    "85ff7d0233ff8b45ec03c63bd87d0a8d56018bcb2bca894dec536a3056ff75ec"
    "57ff3520706300e82932ffff83c418536a30568b45ec03c65057ff3520706300"
    "e81032ffff83c418536a3056ff75ec57ff3520706300e89f31ffff83c418536a"
    "308d560152ff75ec03f756ff3520706300e88431ffff83c418"
)
FRAME_NEW = bytes.fromhex(
    "8b3580c0600089f0f7d883e03fc1e80529c7c1ee0583c7008d43ff39c67e0289"
    "c68b0d84c06000c1e90439c17e0289c129f789d829f04839c77e0289c785ff7d"
    "0231ff8b45f00345f483c00289da29ca4a39d07e0289d085c07d0231c08945ec"
    "8d1408536a30565057ff3520706300e84732ffff89542408e83e32ffff894424"
    "08894c240ce8ae31ffff41894c240c01fe89742404e89e31ffff83c418"
)
FRAME_DX, FRAME_DY = 0x17, 0x4B          # imm8: +dx, +dy

CLICK_VA = 0x44A704
CLICK_ORIG = bytes.fromhex(
    "8b152870630083c2f483faf67e07832d287063000c8b0d287063000faf0d2470"
    "6300a12c70630003c8890d18"
)
CLICK_NEW = bytes.fromhex(
    "8b3580c0600089f7f7df83e73fc1ef05c1ee058d347e83c6008b3d84c06000c1"
    "ef0483c704f7dfe950120400"
)
CLICK_DX, CLICK_DY = 0x18, 0x24          # imm8: -2*dx, +2*dy

CAVE_VA = 0x48B980
CAVE_NEW = bytes.fromhex(
    "8d543702c1fa0229f783c702c1ff0201d1a12870630001f88b157cc06000c1ea"
    "06f7da39d07d0289d08b1d2470630001da39d07e0289d0a32870630085c97d02"
    "31c939d97e0289d9890d2c7063000fafc301c8a318726300e95fedfbff"
)
CAVE_ORIG = bytes(len(CAVE_NEW))

PARTS = [  # (va, original, new)
    (FRAME_VA, FRAME_ORIG, FRAME_NEW + b"\x90" * (FRAME_LEN - len(FRAME_NEW))),
    (CLICK_VA, CLICK_ORIG, CLICK_NEW),
    (CAVE_VA, CAVE_ORIG, CAVE_NEW),
]


def same_but_tune(cur, new, tune_pos):
    a, b = bytearray(cur), bytearray(new)
    for p in tune_pos:
        a[p] = b[p] = 0
    return a == b


def main():
    ap = argparse.ArgumentParser(description="Mayday accurate minimap view frame")
    ap.add_argument("exe")
    ap.add_argument("--dx", type=int, default=0)
    ap.add_argument("--dy", type=int, default=2)
    ap.add_argument("--undo", action="store_true")
    a = ap.parse_args()
    if not (-20 <= a.dx <= 20 and -20 <= a.dy <= 20):
        sys.exit("--dx/--dy must be between -20 and 20")

    data = bytearray(open(a.exe, "rb").read())
    if len(data) != SIZE:
        sys.exit(f"Unexpected size {len(data)} - this patch is only for the Polish Mayday.exe.")

    tunes = [(FRAME_DX, FRAME_DY), (CLICK_DX, CLICK_DY), ()]
    state = []
    for (va, orig, new), t in zip(PARTS, tunes):
        cur = bytes(data[off(va):off(va) + len(new)])
        state.append("orig" if cur == orig else "new" if same_but_tune(cur, new, t) else "?")
    if "?" in state or len(set(state)) > 1:
        sys.exit("This Mayday.exe is not the expected version - not touching the file.")

    if a.undo:
        if state[0] == "orig":
            sys.exit("The frame patch is not applied - nothing to undo.")
        for va, orig, _ in PARTS:
            data[off(va):off(va) + len(orig)] = orig
        open(a.exe, "wb").write(data)
        print("Original minimap frame restored.")
        return

    backup = a.exe + ".frame.bak"
    if state[0] == "orig" and not os.path.exists(backup):
        shutil.copy2(a.exe, backup)
        print(f"Backup written: {backup}")
    for va, _, new in PARTS:
        data[off(va):off(va) + len(new)] = new
    f = off(FRAME_VA)
    c = off(CLICK_VA)
    data[f + FRAME_DX] = a.dx & 0xFF
    data[f + FRAME_DY] = a.dy & 0xFF
    data[c + CLICK_DX] = (-2 * a.dx) & 0xFF
    data[c + CLICK_DY] = (2 * a.dy) & 0xFF
    open(a.exe, "wb").write(data)
    print(f"Accurate minimap frame applied (dx {a.dx}, dy {a.dy}).")


if __name__ == "__main__":
    main()
