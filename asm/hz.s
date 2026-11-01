.intel_syntax noprefix
# Frame pacing at the monitor's refresh rate (e.g. 144 Hz).
# v16: the game's main message loop (GetMessageA at 0x476cd9) hands out a synthetic WM_TIMER
# when the next 1/HZV s is due and NO message is waiting - input and the game's own messages
# always go first (like the original low-priority WM_TIMER). Frames run only while the game's
# own timer 1 is alive: KillTimer (pause, objectives, options, surrender, defeat) stops them at once.
# The 25 Hz game logic is unchanged (fps code decides).
.set HZV,      0x646CC0               # frames per second (monitor refresh rate)
.set DEVMODE,  0x646D00               # 156 bytes for EnumDisplaySettingsA
.set INITHZ,   0x646C60
.set TGT,      0x646C64
.set QPC,      0x646C68
.set QPF_LO,   0x646C6C
.set PENDING,  0x646C70
.set LASTSYNTH,0x646C74               # time of the last synthetic frame
.set LASTGAME, 0x646C78               # time of the last real game-timer tick
.set TIMER_ON, 0x646CC4               # game timer 1 alive (0 after KillTimer)
.set MWFMO,    0x646CC8               # MsgWaitForMultipleObjects
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
.set PeekMessageA,     0x646654
.set GetMessageA_thunk, 0x48b637
.set KillTimer_iat,    0x64660c
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
  cmp dword ptr [PENDING], 0
  je 2f
  mov dword ptr [PENDING], 0             # synthetic frame from hz_getmsg
  jmp 3f
2:
  mov dword ptr [LASTGAME], eax          # the game's own timer tick: timer is alive
  mov dword ptr [TIMER_ON], 1
  mov edx, eax
  sub edx, dword ptr [LASTSYNTH]         # ignored while synthetic frames are running
  cmp edx, 60
  jb 8f
3:
  jmp 0x42a33b                           # 60 fps code with eax = time in ms
8:
  ret
9:
  jmp 0x42a336

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
5:
  mov dword ptr [HZV], 60                # refresh rate of the current display mode
  push offset s_u32
  call dword ptr [GetModuleHandleA]
  test eax, eax
  jz 6f
  mov edi, eax
  push offset s_mw
  push edi
  call dword ptr [GetProcAddress]
  mov dword ptr [MWFMO], eax
  push offset s_eds
  push edi
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
  cmp dword ptr [MWFMO], 0
  je 9f
  mov dword ptr [TGT], eax               # frame pacing on
9:
  popad
  ret

.section .hzdraw,"ax"
# ---- 0x476cd9: GetMessageA(&msg, 0, 0, 0) of the main loop ----
hz_getmsg:
  cmp dword ptr [TGT], 0
  je 8f
  cmp dword ptr [TIMER_ON], 0
  je 8f
  push ebp
  mov ebp, esp
  push ebx
1:
  push 1                                 # PM_REMOVE
  push 0
  push 0
  push 0
  push dword ptr [ebp+8]
  call dword ptr [PeekMessageA]
  test eax, eax
  jz 2f
  mov eax, dword ptr [ebp+8]
  cmp dword ptr [eax+4], 0x12            # WM_QUIT -> 0
  setne al
  movzx eax, al
  jmp 7f
2:
  cmp dword ptr [TIMER_ON], 0            # game timer killed: wait normally
  je 6f
  call dword ptr [TGT]
  mov ebx, eax
  sub eax, dword ptr [LASTGAME]
  cmp eax, 100                           # not in the game screen any more
  jae 6f
  mov eax, ebx
  imul eax, dword ptr [HZV]
  mov edx, eax
  sub edx, dword ptr [NEXT]
  js 4f
  cmp edx, 4000
  jl 3f
  mov dword ptr [NEXT], eax
3:
  add dword ptr [NEXT], 1000
  mov eax, dword ptr [ebp+8]             # synthetic WM_TIMER id 1
  mov ecx, dword ptr [0x60c070]
  mov dword ptr [eax], ecx
  mov dword ptr [eax+4], 0x113
  mov dword ptr [eax+8], 1
  mov dword ptr [eax+12], 0
  mov dword ptr [eax+16], ebx
  mov dword ptr [PENDING], 1
  mov dword ptr [LASTSYNTH], ebx
  mov eax, 1
  jmp 7f
4:
  neg edx                                # sleep until the frame is due or a message arrives
  mov eax, edx
  xor edx, edx
  div dword ptr [HZV]
  inc eax
  push 0x4ff                             # QS_ALLINPUT
  push eax
  push 0
  push 0
  push 0
  call dword ptr [MWFMO]
  jmp 1b
6:
  pop ebx
  pop ebp
8:
  jmp GetMessageA_thunk
7:
  pop ebx
  pop ebp
  ret 16

# ---- 0x48b58f: KillTimer thunk: timer 1 stopped -> no more frames ----
kill_hook:
  cmp dword ptr [esp+8], 1
  jne 1f
  mov dword ptr [TIMER_ON], 0
1:
  jmp dword ptr [KillTimer_iat]

.section .hzdata,"a"
s_u32:   .asciz "user32.dll"
s_eds:   .asciz "EnumDisplaySettingsA"
s_mw:    .asciz "MsgWaitForMultipleObjects"
s_k32:   .asciz "kernel32.dll"
s_ll:    .asciz "LoadLibraryA"
s_winmm: .asciz "winmm.dll"
s_tbp:   .asciz "timeBeginPeriod"
s_tgt:   .asciz "timeGetTime"
