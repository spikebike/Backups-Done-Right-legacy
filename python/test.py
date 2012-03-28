
# to test a 10MB file create with:
# dd if=/dev/urandom of=/tmp/10MB bs=1024 count=10240
# compressed
# padded
# encrypted
# written to file
# decrypted
# unpadded
# decompressed

import hashlib
import sys
import zlib
import os,tempfile
from Crypto.Cipher import AES

input=open("/tmp/100MB","rb")
output=open("100MB-test","wb")
inp=input.read()
e=AES.new("AES key must be either 16, 24, o",AES.MODE_ECB)
c=zlib.compressobj()
ci=c.compress(inp)
ci=ci+c.flush()
prand="0123456789ABCDEF"
m=len(ci)%16;
ci=ci+prand[:-m]
print len(ci),len(ci)%16,m,prand[:-m]
ebuf=e.encrypt(ci)
output.write(ebuf)
input.close()
output.close()

input=open("100MB-test","rb")
output=open("100MB-test2","wb")
inp=input.read()
e=AES.new("AES key must be either 16, 24, o",AES.MODE_ECB)
ebuf=e.decrypt(inp)
out2=zlib.decompress(ebuf[:-(16-m)])
output.write(out2)
output.close()
input.close()


