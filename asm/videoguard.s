.intel_syntax noprefix
# VIDEO panel (and info panels drawn by the routine at 0x44f5xx): the tracked object number
# ([esi+0xeb], global 0x63a24a) can point past the end of the object list (seen: 220 with 32
# objects, crash at 0x44f84b while flying units were selected). Check it against the list size;
# an invalid number is treated as "nothing tracked" (-1) instead of reading outside the list.
.set OBJS,  0x6371f4                 # object list (65 bytes per entry)
.set COUNT, 0x637200                 # pointer to the number of entries
.globl _start
_start:
# ---- 0x44f838: "mov ecx, [esi+0xeb]" (6 bytes) ----
video_guard:
  mov ecx, dword ptr [esi+0xeb]
  cmp dword ptr [OBJS], 0
  je 1f
  mov eax, dword ptr [COUNT]
  test eax, eax
  jz 1f
  cmp ecx, dword ptr [eax]
  jae 1f                             # unsigned: also catches negative numbers
  jmp 0x44f83e                       # valid: continue as the original
1:
  mov dword ptr [esi+0xeb], -1
  jmp 0x44fe41                       # skip drawing the tracked object
