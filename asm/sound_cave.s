.intel_syntax noprefix
.globl _start
_start:
  push 0x498594
  call 0x401f22
  pop ecx
  push dword ptr [0x49871c]
  call 0x40268b
  pop ecx
  jmp 0x402881
