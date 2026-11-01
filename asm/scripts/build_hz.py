import subprocess,sys
o=lambda v:v-0x401000+0x600
src,dst=sys.argv[1],sys.argv[2]
d=bytearray(open(src,'rb').read())
orig=open('patch/Mayday_original.exe','rb').read()
sym={}
for l in subprocess.run(['nm','hz/hz.elf'],capture_output=True,text=True).stdout.splitlines():
    a,t,n=l.split(); sym[n]=int(a,16)
for name,at,end in [('hzcode',0x443230,0x443463),('hzdraw',0x4421d4,0x44242f)]:
    b=open(f'hz/{name}.bin','rb').read(); assert at+len(b)<=end
    assert d[o(at):o(at)+len(b)]==orig[o(at):o(at)+len(b)], name
    d[o(at):o(at)+len(b)]=b
b=open('hz/hzdata.bin','rb').read(); assert d[o(0x42a4a0):o(0x42a4a0)+len(b)]==b'\xcc'*len(b)
d[o(0x42a4a0):o(0x42a4a0)+len(b)]=b
def put(a,old,new):
    old=bytes.fromhex(old); assert d[o(a):o(a)+len(old)]==old,(hex(a),d[o(a):o(a)+len(old)].hex()); assert len(new)==len(old); d[o(a):o(a)+len(new)]=new
rel=lambda a,t:((t-(a+5))&0xffffffff).to_bytes(4,'little')
put(0x476258,'e8d940fbff',b'\xe8'+rel(0x476258,sym['hz_tick']))
put(0x4450ed,'536a0068cc0a6400',b'\xe9'+rel(0x4450ed,sym['ov_hook'])+b'\x90'*3)
open(dst,'wb').write(d); print('built',dst)
