.intel_syntax noprefix
# Frame pacing at the monitor's refresh rate (e.g. 144 Hz), optional frame counter.
# A 1 ms multimedia timer posts the game's WM_TIMER at exactly HZV frames per second
# while the game screen is active; the 25 Hz game logic is unchanged (fps code decides).
.set HZV,      0x646CC0               # frames per second (monitor refresh rate)
.set DEVMODE,  0x646D00               # 156 bytes for EnumDisplaySettingsA
.set INITHZ,   0x646C60
.set TGT,      0x646C64
.set QPC,      0x646C68
.set QPF_LO,   0x646C6C
.set PENDING,  0x646C70
.set LASTPOST, 0x646C74
.set LASTGAME, 0x646C78
.set NEXT,     0x646C7C
.set T0,       0x646C80
.set FRAMES,   0x646C84
.set MAXCUR,   0x646C88
.set SECSTART, 0x646C8C
.set SHOWFPS,  0x646C90
.set SHOWMAX,  0x646C94
.set TMP64,    0x646C98
.set TXT,      0x646CA0
.set GetModuleHandleA, 0x646508
.set GetProcAddress,   0x646554
.set PostMessageA,     0x646650
.set OVX, 200
.set OVY, 40

.section .hzcode,"ax"
.globl _start
_start:
# ---- 0x476258: WM_TIMER in the game screen ----
hz_tick:
  cmp dword ptr [INITHZ], 0
  jne 1f
  call hz_init
1:
  mov eax, dword ptr [TGT]
  test eax, eax
  jz 9f
  call eax                               # timeGetTime (1 ms)
  mov dword ptr [LASTGAME], eax
  cmp dword ptr [PENDING], 0
  je 2f
  mov dword ptr [PENDING], 0             # frame posted by the precise timer
  jmp 3f
2:
  mov edx, eax                           # the game's own 16 ms timer: ignored while
  sub edx, dword ptr [LASTPOST]          # the precise timer is driving the frames
  cmp edx, 60
  jb 8f
3:
  push eax
  call stamp
  pop eax
  jmp 0x42a33b                           # 60 fps code with eax = time in ms
8:
  ret
9:
  call stamp
  jmp 0x42a336

stamp:
  mov eax, dword ptr [QPC]
  test eax, eax
  jz 1f
  push TMP64
  call eax
  mov eax, dword ptr [TMP64]
  or eax, 1
  mov dword ptr [T0], eax
1:
  ret

hz_init:
  mov dword ptr [INITHZ], 1
  pushad
  push offset s_k32
  call dword ptr [GetModuleHandleA]
  mov ebx, eax
  test eax, eax
  jz 9f
  push offset s_ll
  push ebx
  call dword ptr [GetProcAddress]
  test eax, eax
  jz 9f
  push offset s_winmm
  call eax                               # LoadLibraryA("winmm.dll")
  mov esi, eax
  test eax, eax
  jz 9f
  push offset s_qpf
  push ebx
  call dword ptr [GetProcAddress]
  test eax, eax
  jz 5f
  push TMP64
  call eax
  mov eax, dword ptr [TMP64]
  cmp dword ptr [TMP64+4], 0
  jne 5f
  test eax, eax
  jz 5f
  mov dword ptr [QPF_LO], eax
  push offset s_qpc
  push ebx
  call dword ptr [GetProcAddress]
  mov dword ptr [QPC], eax
5:
  mov dword ptr [HZV], 60                # refresh rate of the current display mode
  push offset s_u32
  call dword ptr [GetModuleHandleA]
  test eax, eax
  jz 6f
  push offset s_eds
  push eax
  call dword ptr [GetProcAddress]
  test eax, eax
  jz 6f
  mov word ptr [DEVMODE+36], 156         # dmSize
  push DEVMODE
  push -1                                # ENUM_CURRENT_SETTINGS
  push 0
  call eax
  test eax, eax
  jz 6f
  mov eax, dword ptr [DEVMODE+120]       # dmDisplayFrequency
  cmp eax, 50
  jb 6f
  cmp eax, 360
  ja 6f
  mov dword ptr [HZV], eax
6:
  push offset s_tbp
  push esi
  call dword ptr [GetProcAddress]
  test eax, eax
  jz 9f
  push 1
  call eax                               # timeBeginPeriod(1)
  push offset s_tgt
  push esi
  call dword ptr [GetProcAddress]
  test eax, eax
  jz 9f
  mov edi, eax
  push offset s_tse
  push esi
  call dword ptr [GetProcAddress]
  test eax, eax
  jz 9f
  mov dword ptr [TGT], edi               # before the timer starts
  push 1                                 # TIME_PERIODIC
  push 0
  push offset hz_cb
  push 1
  push 1
  call eax                               # timeSetEvent(1 ms)
  test eax, eax
  jnz 9f
  mov dword ptr [TGT], 0                 # no timer: stay at 60 fps
9:
  popad
  ret

# ---- multimedia timer callback (own thread), every 1 ms ----
hz_cb:
  push ebx
  mov eax, dword ptr [TGT]
  test eax, eax
  jz 9f
  call eax
  mov ebx, eax
  mov edx, eax
  sub edx, dword ptr [LASTGAME]
  cmp edx, 100                           # game screen not active: post nothing
  ja 9f
  mov eax, ebx
  imul eax, dword ptr [HZV]
  mov edx, eax
  sub edx, dword ptr [NEXT]
  js 9f
  cmp edx, 4000
  jl 1f
  mov dword ptr [NEXT], eax
1:
  add dword ptr [NEXT], 1000
  cmp dword ptr [PENDING], 0
  jne 9f
  mov dword ptr [PENDING], 1
  mov dword ptr [LASTPOST], ebx
  push 0                                 # lParam 0 (non-zero would be called as a TimerProc)
  push 1
  push 0x113
  push dword ptr [0x60c070]
  call dword ptr [PostMessageA]
9:
  pop ebx
  ret 20

.section .hzdraw,"ax"
# ---- 0x4450ED: picture goes to the screen: frame counter ----
ov_hook:
  pushad
  mov eax, dword ptr [QPC]
  test eax, eax
  jz 2f
  mov ebx, dword ptr [T0]
  test ebx, ebx
  jz 2f
  mov dword ptr [T0], 0
  push TMP64
  call eax                               # QueryPerformanceCounter (keeps ebx)
  mov eax, dword ptr [TMP64]
  sub eax, ebx
  mov ecx, 10000
  mul ecx
  cmp edx, dword ptr [QPF_LO]
  jae 2f
  div dword ptr [QPF_LO]                 # 0.1 ms units
  cmp eax, dword ptr [MAXCUR]
  jbe 2f
  mov dword ptr [MAXCUR], eax
2:
  inc dword ptr [FRAMES]
  mov eax, dword ptr [TGT]
  test eax, eax
  jz 9f
  call eax
  mov edi, eax
  mov ecx, eax
  sub ecx, dword ptr [SECSTART]
  cmp ecx, 1000
  jb 3f
  mov dword ptr [SECSTART], eax
  mov eax, dword ptr [FRAMES]
  imul eax, eax, 1000
  mov edx, ecx
  shr edx, 1
  add eax, edx
  xor edx, edx
  div ecx
  mov dword ptr [SHOWFPS], eax
  mov eax, dword ptr [MAXCUR]
  mov dword ptr [SHOWMAX], eax
  mov dword ptr [FRAMES], 0
  mov dword ptr [MAXCUR], 0
3:
  cmp byte ptr [show_ov], 0              # frame counter switched on?
  je 9f
  sub edi, dword ptr [LASTGAME]          # only while the game screen is running
  cmp edi, 100
  ja 9f
  call draw_ov
9:
  popad
  push ebx
  push 0
  push 0x640acc
  jmp 0x4450f5

put3:                                    # eax (0..999) -> 3 glyphs at edi, no leading zeros
  mov ecx, 10
  xor edx, edx
  div ecx
  mov byte ptr [edi+2], dl
  xor edx, edx
  div ecx
  mov byte ptr [edi+1], dl
  mov byte ptr [edi], al
  cmp byte ptr [edi], 0
  jne 1f
  mov byte ptr [edi], 17
  cmp byte ptr [edi+1], 0
  jne 1f
  mov byte ptr [edi+1], 17
1:
  ret

draw_ov:
  mov esi, offset tmpl
  mov edi, TXT
  mov ecx, 21
  rep movsb
  mov eax, dword ptr [SHOWFPS]
  cmp eax, 999
  jbe 1f
  mov eax, 999
1:
  mov edi, TXT+4
  call put3
  mov eax, dword ptr [SHOWMAX]
  cmp eax, 9999
  jbe 2f
  mov eax, 9999
2:
  xor edx, edx
  mov ecx, 10
  div ecx
  mov byte ptr [TXT+17], dl
  mov edi, TXT+13
  call put3
  mov esi, dword ptr [0x63a9d0]          # surface memory
  test esi, esi
  jz 9f
  mov ebp, dword ptr [0x63a9bc]          # pitch
  mov eax, dword ptr [0x640ad0]
  add eax, OVY
  imul eax, ebp
  add eax, dword ptr [0x640acc]
  add eax, OVX
  add esi, eax
  mov edi, esi                           # black box
  lea eax, [ebp+ebp*2]
  sub edi, eax
  sub edi, 4
  mov edx, 20
3:
  push edi
  mov ecx, 65
  xor eax, eax
  rep stosd
  pop edi
  add edi, ebp
  dec edx
  jnz 3b
  xor ebx, ebx
4:
  movzx eax, byte ptr [TXT+ebx]
  imul eax, eax, 7
  add eax, offset font
  push esi
  lea edx, [ebx+ebx*2]
  shl edx, 2
  add esi, edx
  mov ecx, 7
5:
  movzx edx, byte ptr [eax]
  inc eax
  push ecx
  mov ecx, 5
  mov edi, esi
6:
  test dl, 0x10
  jz 7f
  mov byte ptr [edi], 0x30
  mov byte ptr [edi+1], 0x30
  mov byte ptr [edi+ebp], 0x30
  mov byte ptr [edi+ebp+1], 0x30
7:
  shl dl, 1
  add edi, 2
  dec ecx
  jnz 6b
  pop ecx
  lea esi, [esi+ebp*2]
  dec ecx
  jnz 5b
  pop esi
  inc ebx
  cmp ebx, 21
  jb 4b
9:
  ret

.section .hzdata,"a"
font:
  .byte 0x0E,0x11,0x13,0x15,0x19,0x11,0x0E   # 0
  .byte 0x04,0x0C,0x04,0x04,0x04,0x04,0x0E   # 1
  .byte 0x0E,0x11,0x01,0x02,0x04,0x08,0x1F   # 2
  .byte 0x1F,0x02,0x04,0x02,0x01,0x11,0x0E   # 3
  .byte 0x02,0x06,0x0A,0x12,0x1F,0x02,0x02   # 4
  .byte 0x1F,0x10,0x1E,0x01,0x01,0x11,0x0E   # 5
  .byte 0x06,0x08,0x10,0x1E,0x11,0x11,0x0E   # 6
  .byte 0x1F,0x01,0x02,0x04,0x08,0x08,0x08   # 7
  .byte 0x0E,0x11,0x11,0x0E,0x11,0x11,0x0E   # 8
  .byte 0x0E,0x11,0x11,0x0F,0x01,0x02,0x0C   # 9
  .byte 0x1F,0x10,0x10,0x1E,0x10,0x10,0x10   # F 10
  .byte 0x1E,0x11,0x11,0x1E,0x10,0x10,0x10   # P 11
  .byte 0x0F,0x10,0x10,0x0E,0x01,0x01,0x1E   # S 12
  .byte 0x11,0x1B,0x15,0x15,0x11,0x11,0x11   # M 13
  .byte 0x0E,0x11,0x11,0x1F,0x11,0x11,0x11   # A 14
  .byte 0x11,0x11,0x0A,0x04,0x0A,0x11,0x11   # X 15
  .byte 0x00,0x00,0x00,0x00,0x00,0x0C,0x0C   # . 16
  .byte 0,0,0,0,0,0,0                        # space 17
tmpl:
  .byte 10,11,12,17, 0,0,0, 17,17, 13,14,15, 17, 0,0,0, 16, 0, 17, 13,12
show_ov: .byte 0
s_u32:   .asciz "user32.dll"
s_eds:   .asciz "EnumDisplaySettingsA"
s_k32:   .asciz "kernel32.dll"
s_ll:    .asciz "LoadLibraryA"
s_winmm: .asciz "winmm.dll"
s_qpf:   .asciz "QueryPerformanceFrequency"
s_qpc:   .asciz "QueryPerformanceCounter"
s_tbp:   .asciz "timeBeginPeriod"
s_tgt:   .asciz "timeGetTime"
s_tse:   .asciz "timeSetEvent"
