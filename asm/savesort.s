.intel_syntax noprefix
# Save lists (*.mds) sorted by date: oldest first, newest at the bottom.
# The directory scan (0x417390) copies each file name into a 0x204-byte slot (0x4b5bde + i*0x204);
# names are at most 260 bytes, so the slot's tail is free. copy_name also stores the file's
# ftLastWriteTime (WIN32_FIND_DATA at 0x4ad9a0, +0x14) at slot+0x1f0; date_cmp replaces the
# name comparator of the file qsort. Equal times fall back to the original name order.
.set STRCPY, 0x482180
.set NAMECMP, 0x41753b
.set FT_LO, 0x4ad9b4
.set FT_HI, 0x4ad9b8
.globl _start
_start:
copy_name:                      # cdecl (dst, src), like strcpy
    mov eax, dword ptr [esp+4]
    mov ecx, dword ptr [FT_LO]
    mov dword ptr [eax+0x1f0], ecx
    mov ecx, dword ptr [FT_HI]
    mov dword ptr [eax+0x1f4], ecx
    jmp STRCPY
date_cmp:                       # qsort comparator (a, b)
    mov eax, dword ptr [esp+4]
    mov ecx, dword ptr [esp+8]
    mov edx, dword ptr [eax+0x1f4]
    cmp edx, dword ptr [ecx+0x1f4]
    jb 1f
    ja 2f
    mov edx, dword ptr [eax+0x1f0]
    cmp edx, dword ptr [ecx+0x1f0]
    jb 1f
    ja 2f
    jmp NAMECMP
1:  or eax, -1
    ret
2:  mov eax, 1
    ret
