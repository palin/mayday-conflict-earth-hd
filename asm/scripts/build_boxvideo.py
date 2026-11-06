import sys
# Box select set the VIDEO panel's unit to the ground-map occupant at the first selected unit's tile;
# for a flying unit that is a wrong or invalid index (blank thumbnail, crash at 0x44f84b). Use the unit itself.
o=lambda v:v-0x401000+0x600
src,dst=sys.argv[1],sys.argv[2]
d=bytearray(open(src,'rb').read())
def put(a,old,new):
    old=bytes.fromhex(old); assert d[o(a):o(a)+len(old)]==old,(hex(a),d[o(a):o(a)+len(old)].hex()); assert len(new)==len(old); d[o(a):o(a)+len(new)]=new
# box select: VIDEO tracked unit = the selected unit itself (ebx), not the ground-tile occupant at its position
put(0x43122f,'8bc3c1e006',bytes.fromhex('8bc3eb1a90'))  # mov eax,ebx ; jmp 0x43124d
open(dst,'wb').write(d); print('built',dst)
