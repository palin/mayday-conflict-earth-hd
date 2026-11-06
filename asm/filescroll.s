.intel_syntax noprefix
# Save/load list (FILE_SCROLLUP / FILE_SCROLLDOWN in the monitor's LOAD/SAVE screen):
# the original moves the list by one line on every pass of the menu loop while the mouse
# button is held - thousands of lines per second on a modern PC. Now: one line at once,
# then (still held) after 400 ms one line every 125 ms.
.set TOP,   0x5ea430                 # first visible line of the file list
.set COUNT_A, 0x5ea428               # folders + files in the list (two counters)
.set COUNT_B, 0x5ea42c
.set LAST,  0x646DA0                 # time of the previous gate call
.set NEXT,  0x646DA4                 # earliest time of the next repeated step
.set GetTickCount, 0x48b397

.globl _start
_start:
# ---- 0x41181d: "dec dword ptr [TOP]" ----
scroll_up:
  call gate
  jnz 1f
  dec dword ptr [TOP]
1:
  ret

# ---- 0x41186a: "inc dword ptr [TOP]" ----
scroll_down:
  call gate
  jnz 1f
  inc dword ptr [TOP]
1:
  ret

# ---- the file list is (re)read: start at the bottom, where the newest saves are ----
# 0x40c79a / 0x40ccc8 / 0x40cd77: "mov [COUNT_A], eax"; 0x41fa5d: "mov [COUNT_A], ecx"
# 0x4115c2 / 0x411b4d: "mov [TOP], edx" (edx = 0, list re-read after a folder change / save)
open_eax:
  mov dword ptr [COUNT_A], eax
  jmp set_bottom
open_ecx:
  mov dword ptr [COUNT_A], ecx
set_bottom:
  push eax
  mov eax, dword ptr [COUNT_A]
  add eax, dword ptr [COUNT_B]
  sub eax, 10                        # 10 lines are visible
  jge 1f
  xor eax, eax
1:
  mov dword ptr [TOP], eax
  pop eax
  ret

# ZF=1: do the step now, ZF=0: skip (registers preserved)
gate:
  push eax
  push ecx
  push edx
  call GetTickCount
  mov edx, eax
  sub edx, dword ptr [LAST]
  mov dword ptr [LAST], eax
  cmp edx, 100                       # no call for 100 ms: a new press
  ja 2f
  mov edx, eax                       # still held: repeat when due
  sub edx, dword ptr [NEXT]
  js 3f
  add eax, 125
  jmp 4f
2:
  add eax, 400
4:
  mov dword ptr [NEXT], eax
  xor eax, eax                       # ZF=1: step
  jmp 5f
3:
  or eax, 1                          # ZF=0: skip
5:
  pop edx
  pop ecx
  pop eax
  ret
