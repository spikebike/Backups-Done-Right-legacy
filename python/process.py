
import hashlib
import sys
import zlib
import os,tempfile
import timeit
from Crypto.Cipher import AES
from Crypto.Cipher import DES
import time


# compression works 10/16/08 - bill
# with no compression or encryption it works 10/16/08 - bill

def process_file(iobj,oobj,blocksize):
	c=zlib.compressobj()
	s=hashlib.sha256()
	e=AES.new("AES key must be either 16, 24, o",AES.MODE_ECB)
	prepend=""
	m=0
	while True:
		buf=iobj.read(65536)
		if not buf:
			break
		cbuf=prepend+c.compress(buf)
		m=len(cbuf)%16;
#		print len(cbuf),len(prepend),m
		if (m>0):
			ebuf=e.encrypt(cbuf[:-m])
	 		s.update(ebuf)
			oobj.write(ebuf)
			prepend=cbuf[-m:]
		else:
			ebuf=e.encrypt(cbuf)
	 		s.update(ebuf)
			oobj.write(ebuf)
			prepend=""

	cbuf=prepend+c.flush()	
	m=len(cbuf)%16;
	if (m>0):  # must clean up
#open("/dev/urandom", "rb").read(16)
		prand="0123456789ABCDEF"
#		print "rand len=",len(prand)
		cbuf=cbuf+prand[:-m]  # pad up to the nearest 16.
		ebuf=e.encrypt(cbuf)
		s.update(ebuf)
		oobj.write(ebuf)
	else:
		ebuf=e.encrypt(cbuf)
	 	s.update(ebuf)
		oobj.write(ebuf)

	return s.hexdigest()	

def filetemp(suffix,prefix,dir):
  (fd, fname) = tempfile.mkstemp(suffix,prefix,dir)
  return (os.fdopen(fd, "w+b"), fname)

def CESfile(fname,blocksize):
	try:
		input = open(fname, 'rb')
#		(output,oname) = filetemp(".bdr","t",".")
#		oname = "100MB.out"
		oname = "/dev/null"
		output = open(oname,'wb')
#		print "input=",input
#		print "output=",output
#		print "fname=",oname
	except:
		return 'Failed to open file'
	ret =  process_file(input,output,blocksize)
	input.close()
	output.close()
	return ret

if __name__ == '__main__':
	blocksize=sys.argv[1]
	for fname in sys.argv[2:]:
		size=os.stat(fname).st_size
		start=time.time()
		sha = CESfile(fname,blocksize)
		stop=time.time()
		diff=stop-start
		sizeMB=size/(1024*1024)
		print "bs=%6s secs=%4.3f size=%d MB bandwith=%5.3f MB/sec" % ( blocksize,diff, sizeMB, sizeMB/diff )

