
import subprocess,os

#a= subprocess.Popen('/bin/zcat','/tmp/100MB.gz',shell=False,bufsize=256)
#a = subprocess.Popen(('/bin/zcat', '/tmp/100MB.gz'),stdout=subprocess.PIPE,bufsize=256)
#a = subprocess.Popen(('/bin/zcat', '/tmp/100MB.gz'),stdout=subprocess.PIPE,bufsize=256).stdout
a = subprocess.Popen("/bin/zcat /tmp/100MB.gz", shell=True, bufsize=256, stdout=subprocess.PIPE).stdout

#a= subprocess.Popen('ls /tmp/', shell=True) 
i=0;
for foo in a.read(256):
	i=i+1
	print "len=",len(foo)
	
print i
