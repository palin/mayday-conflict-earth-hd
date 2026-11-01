.intel_syntax noprefix
.globl _start
# Smooth scrolling for the 1920x1080 mode. Code placed after the 60 fps code (0x42A4A0).
.set VX,   0x646C20      # view position x, 1/256 pixel
.set VY,   0x646C24
.set CI,   0x646C28      # camera tile we last wrote
.set CJ,   0x646C2C
.set APPX, 0x646C30      # shift currently applied to the slot positions
.set APPY, 0x646C34
.set LASTT,0x646C38
.set INIT, 0x646C3C
.set T_I,  0x646C40
.set T_J,  0x646C44
.set T_DX, 0x646C48
.set T_DY, 0x646C4C
.set SCRW, 0x60c080
.set SCRH, 0x60c084
.set CAMI, 0x637028
.set CAMJ, 0x63702c
.set CAMIDX,0x637218
.set MAPW, 0x637024
.set SLOTS,0x60c08c
.set NSLOT,0x4947d4

_start:
# ---- 0x447BC1 (map drawing) jumps here first ----
render_entry:
  pushad
  call smooth_update
  popad
  push ebp
  mov ebp, esp
  add esp, -0x24
  jmp 0x447bc7

# ---- 0x4482A9 (old edge scrolling): skipped in 1920 mode ----
edge_hook:
  cmp dword ptr [SCRW], 1920
  je 1f
  pusha
  mov ecx, dword ptr [0x60c098]
  jmp 0x4482b0
1:
  jmp 0x448397

# ---- 0x444E87: buffer pointer gets a hidden 64 px / 16 row margin in 1920 mode ----
lock_hook:
  mov edx, dword ptr [0x63a9d0]
  mov ecx, dword ptr [0x63a9bc]
  mov dword ptr [0x60c068], ecx
  cmp dword ptr [SCRW], 1920
  jne 1f
  shl ecx, 4
  lea edx, [edx+ecx+64]
1:
  mov dword ptr [0x60c064], edx
  mov eax, 1
  ret

# ---- 0x445087: screen copy rectangle (ebx = width, esi = height) ----
rect_hook:
  mov eax, 32
  mov ecx, 32
  cmp ebx, 1920
  jne 1f
  mov al, 96
  mov cl, 48
1:
  mov dword ptr [0x640acc], eax
  add eax, ebx
  mov dword ptr [0x640ad4], eax
  mov dword ptr [0x640ad0], ecx
  add ecx, esi
  mov dword ptr [0x640ad8], ecx
  jmp 0x4450ac

# ---- per frame: move the view, pick the tile camera, shift the slots ----
smooth_update:
  cmp dword ptr [SCRW], 1920
  jne 9f
  mov esi, dword ptr [SLOTS]
  test esi, esi
  jz 9f
  mov eax, dword ptr [APPX]          # slot 0 must be at x=1952+APPX, y=APPY
  add eax, 1952
  cmp dword ptr [esi], eax
  jne 1f
  mov eax, dword ptr [APPY]
  cmp dword ptr [esi+4], eax
  je 2f
1:
  xor eax, eax                       # slots were rebuilt
  mov dword ptr [APPX], eax
  mov dword ptr [APPY], eax
2:
  mov eax, dword ptr [CAMI]
  mov edx, dword ptr [CAMJ]
  cmp dword ptr [INIT], 0
  je 3f
  cmp eax, dword ptr [CI]
  jne 3f
  cmp edx, dword ptr [CJ]
  je 4f
3:                                   # camera set by the game: follow it
  mov ecx, edx
  sub ecx, eax
  shl ecx, 13
  mov dword ptr [VX], ecx
  lea ecx, [eax+edx]
  shl ecx, 12
  mov dword ptr [VY], ecx
  mov dword ptr [INIT], 1
  call 0x48b397
  mov dword ptr [LASTT], eax
4:
  call 0x48b397                      # GetTickCount
  mov ecx, eax
  sub ecx, dword ptr [LASTT]
  mov dword ptr [LASTT], eax
  cmp ecx, 50
  jbe 5f
  mov ecx, 50
5:
  mov ebx, dword ptr [0x4947f0]      # scroll speed setting
  cmp ebx, 1
  jge 6f
  mov ebx, 1
6:
  imul ebx, ecx
  mov eax, ebx                       # horizontal: 1600 px/s per speed step
  shl eax, 11                        # (n*dt*2048/5 = n*dt*1600*256/1000)
  xor edx, edx
  mov ecx, 5
  div ecx
  mov ebx, eax
  mov eax, dword ptr [0x60c094]
  cmp eax, 5
  jge 7f
  mov edi, ebx
  neg edi
  add dword ptr [VX], edi
7:
  mov ecx, dword ptr [SCRW]
  sub ecx, 5
  cmp eax, ecx
  jle 8f
  add dword ptr [VX], ebx
8:
  sar ebx, 1                         # vertical: half the speed
  mov eax, dword ptr [0x60c098]
  cmp eax, 5
  jge 10f
  mov edi, ebx
  neg edi
  add dword ptr [VY], edi
10:
  mov ecx, dword ptr [SCRH]
  sub ecx, 5
  cmp eax, ecx
  jle 11f
  add dword ptr [VY], ebx
11:
  # keep the view inside the game's camera limits, measured in continuous tile
  # coordinates, so the view slides smoothly along the diamond-shaped map border
  mov eax, dword ptr [VX]
  sar eax, 5                             # u = x / 32   (1/256 tiles)
  mov edx, dword ptr [VY]
  sar edx, 4                             # v = y / 16
  mov ecx, edx
  sub ecx, eax
  sar ecx, 1                             # fi
  add eax, edx
  sar eax, 1                             # fj
  # limits chosen so that the screen centre (camera + 32.375, 2.375 tiles) can reach
  # every tip of the map exactly: i in [-32.375, W-32.375], j in [-2.375, W-2.375]
  xor esi, esi                           # esi = 1 if clamped
  mov edi, dword ptr [MAPW]
  shl edi, 8
  mov edx, -8288                         # lowest i * 256
  cmp ecx, edx
  jge 20f
  mov ecx, edx
  inc esi
20:
  add edx, edi                           # highest i * 256
  cmp ecx, edx
  jle 21f
  mov ecx, edx
  inc esi
21:
  mov edx, -608                          # lowest j * 256
  cmp eax, edx
  jge 22f
  mov eax, edx
  inc esi
22:
  add edx, edi                           # highest j * 256
  cmp eax, edx
  jle 23f
  mov eax, edx
  inc esi
23:
  test esi, esi
  jz 24f
  mov edx, eax
  sub edx, ecx
  shl edx, 5
  mov dword ptr [VX], edx
  add eax, ecx
  shl eax, 4
  mov dword ptr [VY], eax
24:
  mov eax, dword ptr [VX]
  mov edx, dword ptr [VY]
  call derive
  test eax, eax
  jnz 13f
  # outside the game's camera limits: clamp each tile coordinate like the original
  # edge scrolling does (the view then slides along the map border)
  mov ecx, -33                           # lowest i
  mov eax, dword ptr [T_I]
  cmp eax, ecx
  jge 14f
  mov eax, ecx
14:
  mov ecx, dword ptr [MAPW]
  sub ecx, 32                            # highest i
  cmp eax, ecx
  jle 15f
  mov eax, ecx
15:
  mov edx, dword ptr [T_J]
  cmp edx, -3
  jge 16f
  mov edx, -3
16:
  mov ecx, dword ptr [MAPW]
  sub ecx, 2
  cmp edx, ecx
  jle 17f
  mov edx, ecx
17:
  mov dword ptr [T_I], eax
  mov dword ptr [T_J], edx
  mov ecx, edx
  sub ecx, eax
  shl ecx, 13
  mov dword ptr [VX], ecx
  lea ecx, [eax+edx]
  shl ecx, 12
  mov dword ptr [VY], ecx
  xor ecx, ecx
  mov dword ptr [T_DX], ecx
  mov dword ptr [T_DY], ecx
13:
  mov eax, dword ptr [T_I]
  mov edx, dword ptr [T_J]
  mov dword ptr [CAMI], eax
  mov dword ptr [CAMJ], edx
  mov dword ptr [CI], eax
  mov dword ptr [CJ], edx
  imul eax, dword ptr [MAPW]
  add eax, edx
  mov dword ptr [CAMIDX], eax
  mov ecx, dword ptr [T_DX]
  sub ecx, dword ptr [APPX]
  mov edx, dword ptr [T_DY]
  sub edx, dword ptr [APPY]
  mov eax, ecx
  or eax, edx
  jz 9f
  mov esi, dword ptr [SLOTS]
  mov edi, dword ptr [NSLOT]
12:
  add dword ptr [esi], ecx
  add dword ptr [esi+4], edx
  add dword ptr [esi+8], ecx
  add dword ptr [esi+12], edx
  add dword ptr [esi+16], ecx
  add dword ptr [esi+20], edx
  add dword ptr [esi+24], ecx
  add dword ptr [esi+28], edx
  add esi, 32
  dec edi
  jnz 12b
  mov eax, dword ptr [T_DX]
  mov dword ptr [APPX], eax
  mov eax, dword ptr [T_DY]
  mov dword ptr [APPY], eax
9:
  ret

# eax = VX, edx = VY (1/256 px) -> tile camera (T_I,T_J) and shift (T_DX in -63..0, T_DY in -15..0)
# returns eax = 1 if the camera is inside the game's normal limits
derive:
  push ebx
  mov ecx, eax
  sar ecx, 8                         # px
  and ecx, -4                        # whole multiples of 4 px: keeps every tile 4-byte aligned like
                                     # the original (the textured-slope routine reads 4 bytes at a time)
  mov ebx, edx
  sar ebx, 8                         # py
  mov edx, ebx
  sar edx, 4                         # b = floor(py/16)
  mov eax, edx
  shl eax, 4
  sub eax, ebx
  mov dword ptr [T_DY], eax          # 16b - py
  mov eax, ecx
  sar eax, 5                         # a = floor(px/32), same parity as b
  mov ebx, eax
  xor ebx, edx
  test bl, 1
  jz 1f
  dec eax
1:
  mov ebx, eax
  shl ebx, 5
  sub ebx, ecx
  mov dword ptr [T_DX], ebx          # 32a - px
  mov ecx, edx
  sub ecx, eax
  sar ecx, 1                         # i = (b-a)/2
  add eax, edx
  sar eax, 1                         # j = (a+b)/2
  mov dword ptr [T_I], ecx
  mov dword ptr [T_J], eax
  cmp ecx, -34                           # camera limits plus one tile of slack
  jl 2f
  mov edx, dword ptr [MAPW]
  sub edx, 31
  cmp ecx, edx
  jg 2f
  cmp eax, -4
  jl 2f
  mov edx, dword ptr [MAPW]
  dec edx
  cmp eax, edx
  jg 2f
  mov eax, 1
  pop ebx
  ret
2:
  xor eax, eax
  pop ebx
  ret

# ---- 0x405385 (terrain tile blit, ecx = left x, [esp+8] = top y, ret 8):
#      skip tiles that are completely outside the clip area ----
terrain_guard:
  cmp ecx, dword ptr [0x48d2dc]          # left x >= right clip edge
  jge 1f
  push eax
  mov eax, dword ptr [0x48d2d4]
  sub eax, 64
  cmp ecx, eax                           # right edge <= left clip edge
  jle 2f
  mov eax, dword ptr [esp+12]            # y (after push eax)
  cmp eax, dword ptr [0x48d2e0]          # top >= bottom clip edge
  jge 2f
  add eax, 32
  cmp eax, dword ptr [0x48d2d8]          # bottom <= top clip edge
  jle 2f
  pop eax
  push ebp
  push ebx
  push esi
  push edi
  mov ebp, dword ptr [esp+0x14]
  mov esi, edx
  mov ebx, dword ptr [esp+0x18]
  cmp ebx, dword ptr [0x48d2ec]          # keep the game's (working) top/bottom clipping
  jae 0x4079fb
  cmp ebx, dword ptr [0x48d2d8]
  jb 0x406beb
  jmp 0x4053c3                           # no left/right clipping: draw the whole tile,
                                         # the part outside the screen lands in the hidden margin
2:
  pop eax
1:
  ret 8

# ---- 0x44EDF6 (minimap): view frame drawn from the smooth view, clipped to the minimap
#      (it may hang over the minimap edge, like the screen hangs over the map tips) ----
frame_cam:
  mov eax, dword ptr [CAMJ]
  mov dword ptr [ebp-0xc], eax
  mov edx, dword ptr [CAMI]
  mov dword ptr [ebp-0x10], edx
  cmp dword ptr [SCRW], 1920
  jne 1f
  cmp dword ptr [INIT], 0
  jne 5f
1:
  mov edi, dword ptr [MAPW]
  add edi, eax
  dec edi
  sub edi, edx
  lea ecx, [eax+edx]
  test ecx, ecx
  jge 3f
  xor ecx, ecx
3:
  mov dword ptr [ebp-0x14], ecx
  jmp 0x44ee2f
5:
  mov esi, dword ptr [VX]                # left column  x0 = round(VX/32 px) + MAPW - 60
  add esi, 2048
  sar esi, 13
  add esi, dword ptr [MAPW]
  sub esi, 60
  mov edi, dword ptr [VY]                # top row  y0 = VY/16 px + 5.7 (measured on screen)
  add edi, 23347
  sar edi, 12
  mov eax, edi                           # top edge
  mov ecx, esi
  lea edx, [esi+60]
  call hclip
  lea eax, [edi+67]                      # bottom edge
  call hclip
  mov eax, esi                           # left edge
  mov ecx, edi
  lea edx, [edi+67]
  call vclip
  lea eax, [esi+60]                      # right edge
  call vclip
  # minimap scroll position for the panel: frame corner kept inside the minimap
  lea eax, [ebx-1]
  test esi, esi
  jge 6f
  xor esi, esi
6:
  cmp esi, eax
  jle 7f
  mov esi, eax
7:
  test edi, edi
  jge 8f
  xor edi, edi
8:
  cmp edi, eax
  jle 9f
  mov edi, eax
9:
  mov dword ptr [ebp-0x14], edi
  mov edi, esi
  mov esi, 60
  mov ecx, 67
  jmp 0x44eee6

# eax = row, ecx..edx = columns (inclusive); ebx = minimap size
hclip:
  push eax
  push ecx
  push edx
  test eax, eax
  jl 9f
  cmp eax, ebx
  jge 9f
  test ecx, ecx
  jge 1f
  xor ecx, ecx
1:
  cmp edx, ebx
  jl 2f
  lea edx, [ebx-1]
2:
  sub edx, ecx
  jl 9f
  inc edx
  push ebx
  push 0x30
  push edx
  push eax
  push ecx
  push dword ptr [0x637020]
  call 0x4420ea
  add esp, 24
9:
  pop edx
  pop ecx
  pop eax
  ret

# eax = column, ecx..edx = rows (inclusive); ebx = minimap size
vclip:
  push eax
  push ecx
  push edx
  test eax, eax
  jl 9f
  cmp eax, ebx
  jge 9f
  test ecx, ecx
  jge 1f
  xor ecx, ecx
1:
  cmp edx, ebx
  jl 2f
  lea edx, [ebx-1]
2:
  sub edx, ecx
  jl 9f
  inc edx
  push ebx
  push 0x30
  push edx
  push ecx
  push eax
  push dword ptr [0x637020]
  call 0x442067
  add esp, 24
9:
  pop edx
  pop ecx
  pop eax
  ret

# ---- 0x44A704 (minimap click, camera = clicked tile): in smooth mode centre the
#      view exactly on the clicked point (same geometry as the frame above) ----
click_hook:
  cmp dword ptr [SCRW], 1920
  jne 1f
  cmp dword ptr [INIT], 0
  je 1f
  mov eax, dword ptr [CAMJ]
  mov edx, dword ptr [CAMI]
  mov ecx, eax
  sub ecx, edx
  shl ecx, 13
  add ecx, 239616                        # (cj-ci+29.25) * 32 px
  mov dword ptr [VX], ecx
  add eax, edx
  shl eax, 12
  sub eax, 158515                        # (ci+cj-38.7) * 16 px
  mov dword ptr [VY], eax
  mov eax, dword ptr [CI]                # keep the tile camera: the smooth view moves it
  mov edx, dword ptr [CJ]
  mov dword ptr [CAMI], eax
  mov dword ptr [CAMJ], edx
  imul eax, dword ptr [MAPW]
  add eax, edx
  mov dword ptr [CAMIDX], eax
  jmp 0x44a73c
1:
  mov esi, dword ptr [0x60c080]
  jmp 0x44a70a
