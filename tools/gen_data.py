#!/usr/bin/env python3
"""Regenerate patcher/data_gen.go from the chain of verified Mayday.exe builds.

Each fix ("group") is the byte difference between two consecutive builds; the
"Orig" bytes always come from the untouched original exe, so the patcher can
rebuild any combination from the original.

Usage:
  python3 tools/gen_data.py BUILD_DIR ORIGINAL_EXE DATA_GEN_GO
BUILD_DIR holds one folder per build (v2, v4, ... see STAGES), each with Mayday.exe.
Groups for other files (e.g. "panel" for mayday.art) are kept from the existing
DATA_GEN_GO unchanged.

To add a fix: build the new exe from the last stage (e.g. v16 from v15), append
("newid", "v15", "v16") to STAGES, add the group to Groups in core.go (same id,
order = layering order), run this script, run `go test ./...`.
"""
import os, re, sys

STAGES = [  # (group id, previous build, new build); "orig" = the original exe
    ("selection", "orig", "v2"),
    ("outline", "v2", "v4"),
    ("terrain", "v4", "v5"),
    ("rings", "v5", "v6"),
    ("fullhd", "v6", "v9"),
    ("frame", "v9", "v10"),
    ("sound", "v10", "v11"),
    ("terminal", "v11", "v12"),
    ("fps60", "v12", "v13"),
    ("smooth", "v13", "v14h"),
    ("hz", "v14h", "v16"),
]
# Byte changes of older releases that the current groups no longer produce. The patcher only
# uses them to recognise such files and rebuild the original (update without a backup).
LEGACY = [
    ("hz-1.0", "v14h", "v15"),   # 1.0: 144 Hz via a multimedia timer (replaced in 1.0.1)
]
GAP = 16  # merge runs separated by <= GAP bytes that are equal in orig, old and new


def spans(orig, a, b):
    runs, i = [], 0
    while i < len(a):
        if a[i] != b[i]:
            j = i
            while j < len(a) and a[j] != b[j]:
                j += 1
            runs.append([i, j]); i = j
        else:
            i += 1
    m = []
    for r in runs:
        if m and r[0] - m[-1][1] <= GAP and all(orig[k] == a[k] == b[k] for k in range(m[-1][1], r[0])):
            m[-1][1] = r[1]
        else:
            m.append(r)
    return m


def group_go(name, orig, b, sp, var="groupSpans"):
    body = "".join(f'\t\t{{Off: 0x{o:x}, Orig: hx("{orig[o:e].hex()}"), New: hx("{b[o:e].hex()}")}},\n' for o, e in sp)
    return f'\t{var}["{name}"] = []Span{{\n{body}\t}}\n'


def main():
    bdir, orig_path, out = sys.argv[1:4]
    orig = open(orig_path, "rb").read()
    load = lambda v: orig if v == "orig" else open(os.path.join(bdir, v, "Mayday.exe"), "rb").read()
    old = open(out, encoding="utf-8").read()
    kept = {name: whole for whole, name in re.findall(r'(\tgroupSpans\["(\w+)"\] = \[\]Span\{.*?\n\t\}\n)', old, re.S)}
    exe_ids = {s[0] for s in STAGES}
    parts = [group_go(g, orig, load(n), spans(orig, load(p), load(n))) for g, p, n in STAGES]
    parts += [txt for name, txt in kept.items() if name not in exe_ids]
    parts += [group_go(g, orig, load(n), spans(orig, load(p), load(n)), "legacySpans") for g, p, n in LEGACY]
    head = f"// Code generated from the verified builds (orig, {STAGES[0][2]} ... {STAGES[-1][2]}). DO NOT EDIT.\n\npackage main\n\nfunc init() {{\n"
    open(out, "w", encoding="utf-8").write(head + "".join(parts) + "}\n")
    print("written", out)


if __name__ == "__main__":
    main()
