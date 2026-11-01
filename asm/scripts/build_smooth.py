import subprocess,sys
o=lambda v:v-0x401000+0x600
src,dst=sys.argv[1],sys.argv[2]
d=bytearray(open(src,'rb').read())
orig=open('patch/Mayday_original.exe','rb').read()
sym={}
for l in subprocess.run(['nm','smooth/smooth.elf'],capture_output=True,text=True).stdout.splitlines():
    a,t,n=l.split(); sym[n]=int(a,16)
code=open('smooth/smooth.bin','rb').read()
CAVE=0x442c70
assert d[o(0x442c65):o(0x443463)]==orig[o(0x442c65):o(0x443463)]
assert CAVE+len(code)<=0x443463
d[o(CAVE):o(CAVE)+len(code)]=code
def put(a,old,new):
    old=bytes.fromhex(old)
    assert d[o(a):o(a)+len(old)]==old,(hex(a),d[o(a):o(a)+len(old)].hex())
    assert len(new)==len(old); d[o(a):o(a)+len(new)]=new
def jmp(a,target,length,old):
    put(a,old,b'\xe9'+((target-(a+5))&0xffffffff).to_bytes(4,'little')+b'\x90'*(length-5))
jmp(0x447bc1,sym['render_entry'],6,'558bec83c4dc')
jmp(0x4482a9,sym['edge_hook'],7,'608b0d98c06000')
jmp(0x444e87,sym['lock_hook'],24,'8b15d0a96300891564c060008b0dbca96300890d68c06000')
jmp(0x445087,sym['rect_hook'],37,'c705cc0a640020000000c705d00a6400200000008d4b20890dd40a64008d4620a3d80a6400')
jmp(0x405385,sym['terrain_guard'],8,'555356578b6c2414')
jmp(0x44edf6,sym['frame_cam'],57,'a12c7063008945f48b15287063008955f08b3d247063008b45f403f84f8b55f02bfa8b4df08b45f403c8894dec8b55ec85d27d0533c9894dec')
jmp(0x44a704,sym['click_hook'],6,'8b3580c06000')
put(0x44500a,'8b157cc060008915b4a963008b0d78c06000890db8a96300',bytes.fromhex('a17cc0600083c020a3b4a96300a178c06000a3b8a9630090'))
put(0x42d892,'83c240',bytes.fromhex('83ea80'))
put(0x42d888,'e4090000',(2614).to_bytes(4,'little'))
put(0x44e338,'e4090000',(2614).to_bytes(4,'little'))
put(0x44e1e2,'a0070000',(2016).to_bytes(4,'little'))
put(0x44e2f7,'a0070000',(2016).to_bytes(4,'little'))
# top/bottom clipping: never branch into the half-tile corner variants
put(0x407a0f,'0f8715f2ffff',b'\x90'*6)
put(0x407a1b,'0f82edf8ffff',b'\x90'*6)
put(0x406c02,'0f8775ffffff',b'\x90'*6)
put(0x406c0e,'72a4',b'\x90'*2)
# sloped tiles (two triangles): at the left/right screen edge draw both triangles instead of one half
put(0x409eb2,'e91a010000',bytes.fromhex('e928000000'))     # jmp 0x409edf (normal path)
put(0x409ed8,'33d28955f0eb54',bytes.fromhex('eb059090909090')) # jmp 0x409edf
open(dst,'wb').write(d); print('built',dst)
