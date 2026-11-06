import sys
# Main window class without CS_DBLCLKS: the game never handles WM_LBUTTONDBLCLK, so with the
# flag set every second fast click was lost. Style 0x100b -> 0x1003 (RegisterClassA at 0x476b08).
o=lambda v:v-0x401000+0x600
src,dst=sys.argv[1],sys.argv[2]
d=bytearray(open(src,'rb').read())
a=0x476aec
assert d[o(a):o(a)+7]==bytes.fromhex('c745d80b100000'), d[o(a):o(a)+7].hex()
d[o(a)+3]=0x03
open(dst,'wb').write(d); print('built',dst)
