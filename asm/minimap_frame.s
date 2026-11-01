.intel_syntax noprefix
.globl _start
_start:
  mov esi, dword ptr [0x60c080]
  mov eax, esi
  neg eax
  and eax, 63
  shr eax, 5
  sub edi, eax
  shr esi, 5
  .byte 0x83, 0xc7, 0x00          # add edi, DX (tune @ offset)
  lea eax, [ebx-1]
  cmp esi, eax
  jle 1f
  mov esi, eax
1:
  mov ecx, dword ptr [0x60c084]
  shr ecx, 4
  cmp ecx, eax
  jle 2f
  mov ecx, eax
2:
  sub edi, esi
  mov eax, ebx
  sub eax, esi
  dec eax
  cmp edi, eax
  jle 3f
  mov edi, eax
3:
  test edi, edi
  jge 4f
  xor edi, edi
4:
  mov eax, dword ptr [ebp-0x10]
  add eax, dword ptr [ebp-0xc]
  .byte 0x83, 0xc0, 0x02          # add eax, DY
  mov edx, ebx
  sub edx, ecx
  dec edx
  cmp eax, edx
  jle 5f
  mov eax, edx
5:
  test eax, eax
  jge 6f
  xor eax, eax
6:
  mov dword ptr [ebp-0x14], eax
  lea edx, [eax+ecx]
  push ebx
  push 0x30
  push esi
  push eax
  push edi
  push dword ptr [0x637020]
  call 0x4420ea
  mov dword ptr [esp+8], edx
  call 0x4420ea
  mov dword ptr [esp+8], eax
  mov dword ptr [esp+12], ecx
  call 0x442067
  inc ecx
  mov dword ptr [esp+12], ecx
  add esi, edi
  mov dword ptr [esp+4], esi
  call 0x442067
  add esp, 24
