.intel_syntax noprefix
.globl _start
_start:
  lea edx, [edi+esi+2]
  sar edx, 2
  sub edi, esi
  add edi, 2
  sar edi, 2
  add ecx, edx
  mov eax, dword ptr [0x637028]
  add eax, edi
  mov edx, dword ptr [0x60c07c]
  shr edx, 6
  neg edx
  cmp eax, edx
  jge 1f
  mov eax, edx
1:
  mov ebx, dword ptr [0x637024]
  add edx, ebx
  cmp eax, edx
  jle 2f
  mov eax, edx
2:
  mov dword ptr [0x637028], eax
  test ecx, ecx
  jge 3f
  xor ecx, ecx
3:
  cmp ecx, ebx
  jle 4f
  mov ecx, ebx
4:
  mov dword ptr [0x63702c], ecx
  imul eax, ebx
  add eax, ecx
  mov dword ptr [0x637218], eax
  jmp 0x44a73c
