import subprocess,sys
o=lambda v:v-0x401000+0x600
src,dst=sys.argv[1],sys.argv[2]
d=bytearray(open(src,'rb').read())
orig=open('patch/Mayday_original.exe','rb').read()
sym={}
for l in subprocess.run(['nm','filescroll/filescroll.elf'],capture_output=True,text=True).stdout.splitlines():
    a,t,n=l.split(); sym[n]=int(a,16)
CAVE=0x4422f0
b=open('filescroll/filescroll.text.bin','rb').read(); assert CAVE+len(b)<=0x44242f
assert d[o(CAVE):o(CAVE)+len(b)]==orig[o(CAVE):o(CAVE)+len(b)]
d[o(CAVE):o(CAVE)+len(b)]=b
def put(a,old,new):
    old=bytes.fromhex(old); assert d[o(a):o(a)+len(old)]==old,(hex(a),d[o(a):o(a)+len(old)].hex()); assert len(new)==len(old); d[o(a):o(a)+len(new)]=new
rel=lambda a,t:((t-(a+5))&0xffffffff).to_bytes(4,'little')
put(0x41181d,'ff0d30a45e00',b'\xe8'+rel(0x41181d,sym['scroll_up'])+b'\x90')     # FILE_SCROLLUP step
put(0x41186a,'ff0530a45e00',b'\xe8'+rel(0x41186a,sym['scroll_down'])+b'\x90')   # FILE_SCROLLDOWN step
put(0x40c79a,'a328a45e00',b'\xe8'+rel(0x40c79a,sym['open_eax']))              # list opened (main menu load)
put(0x40ccc8,'a328a45e00',b'\xe8'+rel(0x40ccc8,sym['open_eax']))              # list opened
put(0x40cd77,'a328a45e00',b'\xe8'+rel(0x40cd77,sym['open_eax']))              # list opened (save)
put(0x41fa5d,'890d28a45e00',b'\xe8'+rel(0x41fa5d,sym['open_ecx'])+b'\x90')    # list opened (terminal monitor)
put(0x4115c2,'891530a45e00',b'\xe8'+rel(0x4115c2,sym['set_bottom'])+b'\x90')  # re-read after FILE_DIRUP
put(0x411b4d,'891530a45e00',b'\xe8'+rel(0x411b4d,sym['set_bottom'])+b'\x90')  # re-read after saving
open(dst,'wb').write(d); print('built',dst)
