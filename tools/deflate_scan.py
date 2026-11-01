#!/usr/bin/env python3
"""Find raw-deflate streams inside a file (e.g. Wise installers): prints offset, packed size, size, first bytes.
Usage: deflate_scan.py FILE [start_offset]"""
import zlib,sys,hashlib
d=open(sys.argv[1],'rb').read()
res=[];i=0
start=int(sys.argv[2],0) if len(sys.argv)>2 else 0
i=start
while i<len(d)-2:
    try:
        z=zlib.decompressobj(-15)
        out=z.decompress(d[i:i+4000000])
        if z.eof and len(out)>=64:
            used=len(d[i:i+4000000])-len(z.unused_data)
            res.append((i,used,len(out),out[:16]))
            i+=used; continue
    except Exception: pass
    i+=1
for r in res: print(hex(r[0]),r[1],r[2],r[3])
