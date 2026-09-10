from PIL import Image, ImageDraw
from pathlib import Path
import struct, math

OUT=Path('/mnt/data/ConsultantTimer_v8_3')
sizes=[16,24,32,48,64,128,256]
pngs=[]
recording_pngs=[]
for sz in sizes:
    im=Image.new('RGBA',(sz,sz),(0,0,0,0))
    d=ImageDraw.Draw(im)
    pad=max(1,round(sz*0.09))
    stroke=max(2,round(sz*0.085))
    # white clock face with dark outline
    d.ellipse((pad,pad,sz-pad-1,sz-pad-1), fill=(248,248,248,255), outline=(38,38,38,255), width=stroke)
    cx=cy=sz//2
    hand=max(1,round(sz*0.065))
    # minute hand 12 o'clock
    d.line((cx,cy,cx,round(sz*0.28)), fill=(38,38,38,255), width=hand)
    # hour hand ~4 o'clock
    d.line((cx,cy,round(sz*0.70),round(sz*0.61)), fill=(38,38,38,255), width=hand)
    center_dot=max(1,round(sz*0.045))
    d.ellipse((cx-center_dot,cy-center_dot,cx+center_dot,cy+center_dot), fill=(38,38,38,255))
    fp=OUT/f'icon_{sz}.png'
    im.save(fp, optimize=True)
    pngs.append((sz, fp.read_bytes()))

    # Working/recording variant: red status dot at lower-right with a white rim.
    rec=im.copy()
    rd=ImageDraw.Draw(rec)
    r=max(2, round(sz*0.165))
    rim=max(1, round(sz*0.025))
    rcx=sz-r-rim-max(1, round(sz*0.035))-1
    rcy=sz-r-rim-max(1, round(sz*0.035))-1
    rd.ellipse((rcx-r-rim, rcy-r-rim, rcx+r+rim, rcy+r+rim), fill=(255,255,255,255))
    rd.ellipse((rcx-r, rcy-r, rcx+r, rcy+r), fill=(220,35,35,255))
    rfp=OUT/f'icon_recording_{sz}.png'
    rec.save(rfp, optimize=True)
    recording_pngs.append((sz, rfp.read_bytes()))

# multi-size .ico for source package/reference (normal desktop/file icon)
base=Image.open(OUT/'icon_256.png')
base.save(OUT/'ConsultantTimer.ico', format='ICO', sizes=[(s,s) for s in [16,24,32,48,64,128,256]])

class DataRef:
    def __init__(self,data): self.data=data

# Resource tree:
#   RT_GROUP_ICON #1 = normal icon
#   RT_GROUP_ICON #2 = Working icon with red recording dot
# The normal icon remains the executable/Desktop icon.
icon_type={}
for i,(sz,data) in enumerate(pngs,1):
    icon_type[i]={1033:DataRef(data)}
for i,(sz,data) in enumerate(recording_pngs,8):
    icon_type[i]={1033:DataRef(data)}

def make_group(items, first_id):
    entries=[]
    for offset,(sz,data) in enumerate(items):
        icon_id=first_id+offset
        wh=0 if sz>=256 else sz
        entries.append(struct.pack('<BBBBHHIH', wh, wh, 0, 0, 1, 32, len(data), icon_id))
    return struct.pack('<HHH',0,1,len(entries))+b''.join(entries)

grp_normal=make_group(pngs,1)
grp_recording=make_group(recording_pngs,8)
group_type={1:{1033:DataRef(grp_normal)}, 2:{1033:DataRef(grp_recording)}}
root={3:icon_type,14:group_type}

def align(v,a): return (v+a-1)//a*a

def build_rsrc(section_rva):
    buf=bytearray()
    data_records=[] # (data_entry_offset, bytes)
    def emit_dir(node):
        off=len(buf)
        ids=sorted(node.keys())
        buf.extend(struct.pack('<IIHHHH',0,0,0,0,0,len(ids)))
        ent_offsets=[]
        for rid in ids:
            ent_offsets.append(len(buf))
            buf.extend(struct.pack('<II',rid,0))
        for ent_off,rid in zip(ent_offsets,ids):
            child=node[rid]
            if isinstance(child,DataRef):
                de_off=len(buf)
                buf.extend(b'\0'*16)
                struct.pack_into('<I',buf,ent_off+4,de_off)
                data_records.append((de_off,child.data))
            else:
                child_off=emit_dir(child)
                struct.pack_into('<I',buf,ent_off+4,0x80000000|child_off)
        return off
    emit_dir(root)
    # append payloads and patch IMAGE_RESOURCE_DATA_ENTRY
    for de_off,data in data_records:
        while len(buf)%4: buf.append(0)
        blob_off=len(buf)
        buf.extend(data)
        struct.pack_into('<IIII',buf,de_off,section_rva+blob_off,len(data),0,0)
    while len(buf)%4: buf.append(0)
    return bytes(buf)

def patch_pe(exe_path):
    p=Path(exe_path)
    b=bytearray(p.read_bytes())
    peoff=struct.unpack_from('<I',b,0x3c)[0]
    assert b[peoff:peoff+4]==b'PE\0\0'
    coff=peoff+4
    machine,nsec,tstamp,psym,nsym,opt_size,chars=struct.unpack_from('<HHIIIHH',b,coff)
    assert machine==0x8664
    opt=coff+20
    magic=struct.unpack_from('<H',b,opt)[0]
    assert magic==0x20b
    sect_align=struct.unpack_from('<I',b,opt+32)[0]
    file_align=struct.unpack_from('<I',b,opt+36)[0]
    size_image=struct.unpack_from('<I',b,opt+56)[0]
    size_headers=struct.unpack_from('<I',b,opt+60)[0]
    dd=opt+112
    sh=opt+opt_size
    first_raw=10**18
    max_end_va=0
    for i in range(nsec):
        o=sh+40*i
        vs,va,rs,pr=struct.unpack_from('<IIII',b,o+8)
        if pr: first_raw=min(first_raw,pr)
        max_end_va=max(max_end_va, va+max(vs,rs))
    new_hdr=sh+40*nsec
    if new_hdr+40>first_raw:
        raise RuntimeError('No room for new section header')
    new_va=align(max_end_va,sect_align)
    rsrc=build_rsrc(new_va)
    new_raw=align(len(b),file_align)
    if len(b)<new_raw: b.extend(b'\0'*(new_raw-len(b)))
    raw_size=align(len(rsrc),file_align)
    b.extend(rsrc)
    if len(rsrc)<raw_size: b.extend(b'\0'*(raw_size-len(rsrc)))
    # section header
    name=b'.rsrc\0\0\0'
    b[new_hdr:new_hdr+8]=name
    struct.pack_into('<IIIIIIHHI',b,new_hdr+8,len(rsrc),new_va,raw_size,new_raw,0,0,0,0,0x40000040)
    # NumberOfSections
    struct.pack_into('<H',b,coff+2,nsec+1)
    # SizeOfInitializedData
    old_init=struct.unpack_from('<I',b,opt+8)[0]
    struct.pack_into('<I',b,opt+8,old_init+raw_size)
    # SizeOfImage
    struct.pack_into('<I',b,opt+56,align(new_va+len(rsrc),sect_align))
    # resource data directory
    struct.pack_into('<II',b,dd+2*8,new_va,len(rsrc))
    # checksum 0
    struct.pack_into('<I',b,opt+64,0)
    p.write_bytes(b)
    return new_va,len(rsrc),new_raw,raw_size

if __name__=='__main__':
    exe=OUT/'ConsultantTimer.exe'
    va,vs,raw,rs=patch_pe(exe)
    print(f'embedded .rsrc RVA=0x{va:x}, size={vs}, raw=0x{raw:x}/{rs}')
