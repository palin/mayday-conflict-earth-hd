#!/usr/bin/env python3
"""Mayday: Conflict Earth (PL, Mayday.exe 684 032 bytes) - choose the HIRES resolution.

The game's "1024 X 768" (HIRES) setting is turned into another resolution.
640 x 480 and 800 x 600 stay as they are. Works with or without the other patches,
and can switch between modes at any time.

    python res_mode.py "D:\\Retro\\Mayday\\Mayday.exe" --mode 1920x1080
    python res_mode.py "D:\\Retro\\Mayday\\Mayday.exe" --mode 1280x1024
    python res_mode.py "D:\\Retro\\Mayday\\Mayday.exe" --mode 1024x768     (original)

How the values are found: the game builds its on-screen grid of 64 x 32 map tiles
from a few numbers per resolution. The rule below reproduces the game's own numbers
exactly for 640x480, 800x600 and 1024x768, and is used for the new sizes.
"""
import argparse
import sys

SIZE = 684032
MODES = {"1024x768": (1024, 768), "1280x1024": (1280, 1024), "1920x1080": (1920, 1080)}


def tiles(XL, rowsA, YL, rowsB):
    n, x0 = 0, XL
    for _ in range(rowsA):
        x = x0
        while True:
            n += 1
            x += 32
            if x > XL:
                break
        x0 -= 64
    yst = 32
    for _ in range(rowsB):
        x, y = 32, yst
        while True:
            n += 1
            x += 32
            y += 16
            if y > YL or x > XL:
                break
        yst += 32
    return n


def layout(W, H):
    XL = W + 32
    while XL % 64 != 32:
        XL += 1
    YL = H + 224
    while YL % 32:
        YL += 1
    rowsA, rowsB = (XL + 32) // 64, YL // 32
    return dict(W=W, H=H, XL=XL, YL=YL, rowsA=rowsA, rowsB=rowsB,
                bufW=XL + 32 - W, bufH=YL + 32 - H, cull=YL - H - 192,
                tiles=tiles(XL, rowsA, YL, rowsB))


def d32(v): return v.to_bytes(4, "little")


def label(W, H):
    s = f"{W}X{H}".encode()
    return (s + b"\x00" * 10)[:10]


def patch_table(m):
    """(address, bytes) for one mode"""
    L = layout(*MODES[m])
    W, H = L["W"], L["H"]
    assert L["bufW"] < 128 and L["cull"] < 128 and L["rowsA"] < 128 and L["rowsB"] < 128
    lab = b"1024 X 768" if m == "1024x768" else label(W, H)
    return [
        (0x42C654, b"\x68" + d32(H)),                       # push height
        (0x42C659, b"\x68" + d32(W)),                       # push width
        (0x42D86E, b"\x81\xf9" + d32(W)),                   # screen setup: cmp width
        (0x42D87B, b"\x3d" + d32(H)),                       # screen setup: cmp height
        (0x42D882, b"\xc7\x05\xd4\x47\x49\x00" + d32(L["tiles"])),
        (0x42D892, b"\x83\xc2" + bytes([L["bufW"]])),      # off-screen buffer width  = W + bufW
        (0x42D8A1, b"\x81\xc1" + d32(L["bufH"])),           # off-screen buffer height = H + bufH
        (0x42D8C3, b"\x83\xc2" + bytes([L["cull"]])),      # bottom drawing limit     = H + cull
        (0x42C736, b"\x3d" + d32(W)),                       # use the high-res tile grid
        (0x44E104, b"\xb9" + d32(L["XL"])),                 # grid A start x
        (0x44E1E0, b"\x81\xf9" + d32(L["XL"])),             # grid A right edge
        (0x44E209, b"\x83\xfe" + bytes([L["rowsA"]])),     # grid A rows
        (0x44E2ED, b"\x81\xfb" + d32(L["YL"])),             # grid B bottom edge
        (0x44E2F5, b"\x81\xf9" + d32(L["XL"])),             # grid B right edge
        (0x44E321, b"\x83\xfe" + bytes([L["rowsB"]])),     # grid B rows
        (0x44E332, b"\xc7\x05\xd4\x47\x49\x00" + d32(L["tiles"])),
        (0x444EC0, b"\x3d" + d32(W)),                       # display mode: cmp width
        (0x444EC7, b"\xbb" + d32(W)),                       # display mode: width
        (0x444ECC, b"\xbe" + d32(H)),                       # display mode: height
        (0x49431B, lab),                                    # menu label (DATA section)
    ]


def off(va):
    return va - 0x401000 + 0x600 if va < 0x48C000 else va - 0x401000


def main():
    ap = argparse.ArgumentParser(description="Mayday HIRES resolution")
    ap.add_argument("exe")
    ap.add_argument("--mode", choices=MODES, required=True)
    a = ap.parse_args()
    data = bytearray(open(a.exe, "rb").read())
    if len(data) != SIZE:
        sys.exit(f"Unexpected size {len(data)} - this patch is only for the Polish Mayday.exe.")
    tables = {m: patch_table(m) for m in MODES}
    # every location must currently hold the bytes of SOME known mode
    for i, (va, _) in enumerate(tables["1024x768"]):
        n = len(tables["1024x768"][i][1])
        cur = bytes(data[off(va):off(va) + n])
        if not any(cur == tables[m][i][1] for m in MODES):
            sys.exit(f"Unexpected bytes at {va:08X} - not touching the file.")
    for va, b in tables[a.mode]:
        data[off(va):off(va) + len(b)] = b
    open(a.exe, "wb").write(data)
    L = layout(*MODES[a.mode])
    print(f"HIRES setting is now {a.mode}: tile area {L['XL']} x {L['YL']} px, "
          f"{L['tiles']} visible tiles, off-screen buffer {L['W'] + L['bufW']} x {L['H'] + L['bufH']}.")


if __name__ == "__main__":
    main()
