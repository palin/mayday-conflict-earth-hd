.intel_syntax noprefix
.globl _start
# Code at 0x42A336 (old, never-called frame function). Variables in unused .idata space.
.set SAVEBUF, 0x646800
.set LAST,    0x646C00
.set SAVED,   0x646C04
.set SAVEPTR, 0x646C08
.set RONLY,   0x646C0C
.set SAVEBASE,0x646C10

_start:
# ---- tick: replaces "call 0x448791" in the WM_TIMER handler ----
tick:
  call 0x48b397                     # GetTickCount
  mov edx, eax
  sub edx, dword ptr [LAST]
  cmp edx, 40                       # logic every 40 ms (25 per second, as the original)
  jb render_only
  cmp edx, 160
  jae 1f
  add dword ptr [LAST], 40
  jmp 2f
1:
  mov dword ptr [LAST], eax
2:
  call restore_cur
  jmp 0x448791                      # original tick: logic + drawing

render_only:
  sub esp, 8
  push esp
  call 0x48b649                     # GetCursorPos(&pt)
  pop eax
  pop edx
  mov dword ptr [0x60c094], eax
  mov dword ptr [0x60c098], edx
  call restore_cur
  push dword ptr [0x494910]         # snapshot counters the drawing code advances
  mov al, byte ptr [0x494914]
  mov ah, byte ptr [0x494811]
  push eax
  mov al, byte ptr [0x63727e]
  mov ah, byte ptr [0x63727c]
  push eax
  mov dword ptr [RONLY], 1
  call 0x447bc1                     # draw map + cursor
  call 0x4674ea
  mov dword ptr [RONLY], 0
  pop eax
  mov byte ptr [0x63727e], al
  mov byte ptr [0x63727c], ah
  pop eax
  mov byte ptr [0x494914], al
  mov byte ptr [0x494811], ah
  pop dword ptr [0x494910]
  jmp 0x4450ed                      # copy the picture to the screen

# ---- cursor: replaces "call 0x4537df" at the end of the drawing ----
cursor_hook:
  push esi
  push edi
  push ecx
  mov eax, dword ptr [0x60c098]
  add eax, 16
  imul eax, dword ptr [0x60c068]
  add eax, dword ptr [0x60c094]
  add eax, 16
  mov edx, dword ptr [0x60c064]
  mov dword ptr [SAVEBASE], edx
  add eax, edx
  mov dword ptr [SAVEPTR], eax
  mov esi, eax
  mov edi, SAVEBUF
  mov edx, 32
3:
  mov ecx, 8
  rep movsd
  add esi, dword ptr [0x60c068]
  sub esi, 32
  dec edx
  jnz 3b
  mov dword ptr [SAVED], 1
  pop ecx
  pop edi
  pop esi
  jmp 0x4537df

restore_cur:
  cmp dword ptr [SAVED], 0
  je 5f
  mov dword ptr [SAVED], 0
  mov eax, dword ptr [0x60c064]
  cmp eax, dword ptr [SAVEBASE]
  jne 5f
  push esi
  push edi
  push ecx
  mov edi, dword ptr [SAVEPTR]
  mov esi, SAVEBUF
  mov edx, 32
4:
  mov ecx, 8
  rep movsd
  add edi, dword ptr [0x60c068]
  sub edi, 32
  dec edx
  jnz 4b
  pop ecx
  pop edi
  pop esi
5:
  ret

# ---- effects animation: replaces "call 0x44d725" inside the drawing ----
fx_hook:
  cmp dword ptr [RONLY], 0
  jne 6f
  jmp 0x44d725
6:
  ret
