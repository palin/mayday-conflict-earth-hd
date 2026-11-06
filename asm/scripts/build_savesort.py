import subprocess,sys
# Save list sorted by file date (newest at the bottom); see savesort.s
o=lambda v:v-0x401000+0x600
src,dst=sys.argv[1],sys.argv[2]
d=bytearray(open(src,'rb').read())
orig=open('patch/Mayday_original.exe','rb').read()
sym={}
for l in subprocess.run(['nm','asm/savesort.elf'],capture_output=True,text=True).stdout.splitlines():
    a,t,n=l.split(); sym[n]=int(a,16)
CAVE=0x4423b0
b=open('asm/savesort.text.bin','rb').read(); assert CAVE+len(b)<=0x44242f
assert d[o(CAVE):o(CAVE)+len(b)]==orig[o(CAVE):o(CAVE)+len(b)]
d[o(CAVE):o(CAVE)+len(b)]=b
def put(a,old,new):
    old=bytes.fromhex(old); assert d[o(a):o(a)+len(old)]==old,(hex(a),d[o(a):o(a)+len(old)].hex()); assert len(new)==len(old); d[o(a):o(a)+len(new)]=new
rel=lambda a,t:((t-(a+5))&0xffffffff).to_bytes(4,'little')
put(0x41747e,'e8fdac0600',b'\xe8'+rel(0x41747e,sym['copy_name']))   # first file name
put(0x4174be,'e8bdac0600',b'\xe8'+rel(0x4174be,sym['copy_name']))   # next file names
put(0x417520,'683b754100',b'\x68'+sym['date_cmp'].to_bytes(4,'little'))  # file sort comparator
open(dst,'wb').write(d); print('built',dst, 'cave end', hex(CAVE+len(b)))
