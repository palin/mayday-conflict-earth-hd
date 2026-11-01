#!/usr/bin/env python3
"""Mayday: Conflict Earth (PL, Mayday.exe 684 032 bytes) - "no sound at start" fix.

Cause: when the CD music is paused (ESC / pause menu), the game fades the CD volume
to zero through the old Windows mixer and leaves it there. On Windows 10/11 that
mixer setting is the program's own volume in the Windows Volume Mixer, and Windows
remembers it per output device. Quit while paused -> Mayday starts silent next time.

Second cause: when quitting, the game "restores" the CD volume it read when it
started. On Windows 10/11 that is Mayday's own Volume Mixer level - so once it is
0, the game reads 0 at start and writes 0 back at exit, every time.

Fix:
 1. On exit the game sets full volume instead of the value read at start.
 2. Right after the music is stopped by the pause, the previous volume is set
    again (the music is already stopped, so nothing is heard).

Changes: 0040289B (file 0x1E9B, 6 bytes), 0040284B (file 0x1E4B, 13 bytes) and a
code cave at 0048B9E0 (file 0x8AFE0, 28 bytes).

If Mayday is already at 0: start the patched game, set Mayday to 100 in the
Volume Mixer once, and quit the game normally (menu -> exit). From then on it stays.

Usage:
  python cd_volume_fix.py "D:\\Retro\\Mayday\\Mayday.exe"
  python cd_volume_fix.py "D:\\Retro\\Mayday\\Mayday.exe" --undo
Works on the original exe or one with the earlier patches.
"""
import sys

SIZE = 684032


def off(va):
    return va - 0x401000 + 0x600


PARTS = [  # (va, original, new)
    # on exit the game "restores" the volume it read at startup - on Win10/11 that is
    # Mayday's own Volume Mixer level, so a 0 keeps itself alive forever; restore full instead
    (0x40289B, bytes.fromhex("ff3518874900"), bytes.fromhex("68ffff000090")),
    (0x40284B, bytes.fromhex("6894854900e8cdf6ffff59eb29"),
     bytes.fromhex("e990910800" + "90" * 8)),
    (0x48B9E0, bytes(28),
     bytes.fromhex("6894854900e83865f7ff59ff351c874900e8956cf7ff59e9856ef7ff")),
]


def main():
    if len(sys.argv) < 2:
        sys.exit(__doc__)
    path, undo = sys.argv[1], "--undo" in sys.argv[2:]
    data = bytearray(open(path, "rb").read())
    if len(data) != SIZE:
        sys.exit(f"Unexpected size {len(data)} - this patch is only for the Polish Mayday.exe.")
    state = set()
    for va, orig, new in PARTS:
        cur = bytes(data[off(va):off(va) + len(new)])
        state.add("orig" if cur == orig else "new" if cur == new else "?")
    first = [bytes(data[off(va):off(va) + len(n)]) for va, _, n in PARTS]
    if (first[0] == PARTS[0][1] and first[1] == PARTS[1][2] and first[2] == PARTS[2][2]
            and not undo):
        data[off(PARTS[0][0]):off(PARTS[0][0]) + 6] = PARTS[0][2]
        open(path, "wb").write(data)
        print("Sound fix updated to version 2.")
        return
    if "?" in state or len(state) > 1:
        sys.exit("This Mayday.exe is not the expected version - not touching the file.")
    cur = state.pop()
    if undo:
        if cur == "orig":
            sys.exit("The sound fix is not applied - nothing to undo.")
        for va, orig, _ in PARTS:
            data[off(va):off(va) + len(orig)] = orig
        open(path, "wb").write(data)
        print("Sound fix removed.")
        return
    if cur == "new":
        sys.exit("The sound fix is already applied.")
    for va, _, new in PARTS:
        data[off(va):off(va) + len(new)] = new
    open(path, "wb").write(data)
    print("Sound fix applied.")


if __name__ == "__main__":
    main()
