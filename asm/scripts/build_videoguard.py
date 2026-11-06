# VIDEO panel guard: an invalid tracked unit index (>= object count) closes the panel instead of crashing.
import subprocess,sys
o=lambda v:v-0x401000+0x600
src,dst=sys.argv[1],sys.argv[2]
d=bytearray(open(src,'rb').read())
orig=open('patch/Mayday_original.exe','rb').read()
sym={}
for l in subprocess.run(['nm','asm/videoguard.elf'],capture_output=True,text=True).stdout.splitlines():
    a,t,n=l.split(); sym[n]=int(a,16)
CAVE=0x442380
b=open('asm/videoguard.text.bin','rb').read(); assert CAVE+len(b)<=0x44242f
assert d[o(CAVE):o(CAVE)+len(b)]==orig[o(CAVE):o(CAVE)+len(b)]
d[o(CAVE):o(CAVE)+len(b)]=b
def put(a,old,new):
    old=bytes.fromhex(old); assert d[o(a):o(a)+len(old)]==old,(hex(a),d[o(a):o(a)+len(old)].hex()); assert len(new)==len(old); d[o(a):o(a)+len(new)]=new
rel=lambda a,t:((t-(a+5))&0xffffffff).to_bytes(4,'little')
put(0x44f838,'8b8eeb000000',b'\xe9'+rel(0x44f838,sym['video_guard'])+b'\x90')   # tracked object check
open(dst,'wb').write(d); print('built',dst)
