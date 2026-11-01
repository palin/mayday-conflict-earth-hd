#!/usr/bin/env python3
"""Print the import table (IAT address, DLL, function) of a 32-bit PE. Usage: pe_imports.py FILE"""
import struct,sys
d=open(sys.argv[1],'rb').read()
pe=struct.unpack_from('<I',d,0x3c)[0]
nsec=struct.unpack_from('<H',d,pe+6)[0]; opt=pe+24; osz=struct.unpack_from('<H',d,pe+20)[0]
base=struct.unpack_from('<I',d,opt+28)[0]
imp_rva=struct.unpack_from('<I',d,opt+96+8)[0]
secs=[]
for k in range(nsec):
    s=pe+24+osz+40*k
    name=d[s:s+8].rstrip(b'\0').decode(); vs,va,rs,ra=struct.unpack_from('<IIII',d,s+8); secs.append((name,va,vs,ra,rs))
def r2o(r):
    for n,va,vs,ra,rs in secs:
        if va<=r<va+max(vs,rs): return r-va+ra
p=r2o(imp_rva); res={}
while True:
    olt,ts,fc,nm,ft=struct.unpack_from('<IIIII',d,p)
    if nm==0: break
    dll=d[r2o(nm):].split(b'\0')[0].decode(); q=r2o(olt or ft); k=0
    while True:
        t=struct.unpack_from('<I',d,q+4*k)[0]
        if t==0: break
        fn=d[r2o(t)+2:].split(b'\0')[0].decode() if not t&0x80000000 else '#%d'%(t&0xffff)
        res[base+ft+4*k]=(dll,fn); k+=1
    p+=20
print(secs)
for a,(dl,fn) in sorted(res.items()): print(hex(a),dl,fn)
