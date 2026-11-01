#!/usr/bin/env python3
"""Minimal Windows minidump reader (modules, exception, thread context, memory). Used to find the
exact crashing instruction from WER LocalDumps (Mayday_crashdumps_ON.reg)."""
import struct,sys
class Dump:
    def __init__(s,path):
        s.d=d=open(path,'rb').read()
        sig,ver,n,rva=struct.unpack_from('<4sIII',d,0)
        assert sig==b'MDMP'
        s.streams={}
        for k in range(n):
            t,size,r=struct.unpack_from('<III',d,rva+12*k); s.streams.setdefault(t,[]).append((r,size))
        s.mem=[]
        if 9 in s.streams:
            r,_=s.streams[9][0]; cnt,base=struct.unpack_from('<QQ',d,r)
            off=base
            for k in range(cnt):
                a,sz=struct.unpack_from('<QQ',d,r+16+16*k); s.mem.append((a,sz,off)); off+=sz
        if 5 in s.streams:
            r,_=s.streams[5][0]; cnt=struct.unpack_from('<I',d,r)[0]
            for k in range(cnt):
                a,sz,rv=struct.unpack_from('<QII',d,r+4+16*k); s.mem.append((a,sz,rv))
    def read(s,addr,n):
        for a,sz,off in s.mem:
            if a<=addr and addr+n<=a+sz: return s.d[off+addr-a:off+addr-a+n]
        return None
    def u32(s,a):
        b=s.read(a,4); return None if b is None else struct.unpack('<I',b)[0]
    def i32(s,a):
        b=s.read(a,4); return None if b is None else struct.unpack('<i',b)[0]
