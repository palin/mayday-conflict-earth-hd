.intel_syntax noprefix
.globl _start
# Terminal scroll, run at a fixed rate (default 30 steps/s) instead of once per frame.
# pos = [0x492410] (int, 60..1400), vel = [0x492414] (int, reused from the old double),
# mouse x = [0x4923ec] (front-end runs at 640x480), last step time = [0x644F00] (free BSS).
_start:
  call 0x48b397                 # GetTickCount
  mov edx, eax
  sub edx, dword ptr [0x644f00]
  .byte 0x83, 0xfa, 0x21        # cmp edx, STEP (ms, tune byte)  -> offset 0x0f
  jb done
  .byte 0x81, 0xfa, 0x84, 0, 0, 0   # cmp edx, 4*STEP (tune dword) -> offset 0x16
  jae reset
  .byte 0x81, 0x05, 0x00, 0x4f, 0x64, 0x00, 0x21, 0, 0, 0   # add dword [0x644f00], STEP -> offset 0x22
  jmp step
reset:
  mov dword ptr [0x644f00], eax
step:
  mov eax, dword ptr [0x4923ec]
  mov ecx, dword ptr [0x492414]
  cmp eax, 64
  jge right
  imul eax, eax, 15
  add eax, 32
  sar eax, 6
  sub eax, 16                   # -16 at the edge .. -1 near x=63
  cmp eax, ecx
  jge apply
  mov ecx, eax
  jmp apply
right:
  cmp eax, 576
  jle apply
  sub eax, 576
  imul eax, eax, 15
  add eax, 32
  sar eax, 6                    # 0 .. 15 at the right edge
  cmp eax, ecx
  jle apply
  mov ecx, eax
apply:
  test ecx, ecx
  je store
  mov eax, dword ptr [0x492410]
  add eax, ecx
  cmp eax, 60
  jge 1f
  mov eax, 60
1:
  cmp eax, 1400
  jle 2f
  mov eax, 1400
2:
  mov dword ptr [0x492410], eax
  test ecx, ecx
  jg 3f
  inc ecx
  jmp store
3:
  dec ecx
store:
  mov dword ptr [0x492414], ecx
done:
  ret
