.intel_syntax noprefix
.globl _start
_start:
  mov esi, dword ptr [0x60c080]
  mov edi, esi
  neg edi
  and edi, 63
  shr edi, 5
  shr esi, 5
  lea esi, [esi+edi*2]
  .byte 0x83, 0xc6, 0x00          # add esi, -2*DX
  mov edi, dword ptr [0x60c084]
  shr edi, 4
  .byte 0x83, 0xc7, 0x04          # add edi, 2*DY
  neg edi
  jmp 0x48b980
