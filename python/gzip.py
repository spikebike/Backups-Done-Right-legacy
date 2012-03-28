import os
import sys
import gzip

input=open('test','rb',16)

#fd = gzip.open(sys.argv[1], 'r')
#
#ln = fd.readline()

for buffer in input.read(16):
	print len(buffer),buffer
#    sys.stdout.write(ln)

#    ln = fd.readline()

