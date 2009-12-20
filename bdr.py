#!/usr/bin/python
import os
from os.path import join, getsize
import time
import datetime
#from pysqlite2 import dbapi2 as sqlite
import sqlite3


startBackup=time.time()

connection = sqlite3.connect('test.db')
connection.row_factory = sqlite3.Row

cursor = connection.cursor()
try:
	cursor.execute('create table dirs (id INTEGER PRIMARY KEY,st_mode INT, st_ino BIGINT, st_uid INT, st_gid INT, name varchar(2048), last_seen ts, deleted INT)')
	cursor.execute('create table files (id INTEGER PRIMARY KEY, st_mode INT, st_ino BIGINT, st_dev BIGINT, st_nlink INT, st_uid INT, st_gid INT, st_size BIGINT, st_atime BIGINT, st_mtime BIGINT, st_ctime BIGINT, name varchar(255), dirID BIGINT, last_seen ts, deleted INT, FOREIGN KEY(dirID) REFERENCES dirs(id))')
except sqlite3.OperationalError, msg:
	print "tables already exist"

def addDir(dirname):
	st=os.stat(dirname)
#	/* do we have a perfect match */
	cursor.execute("select * from dirs where name=? and st_mode=? and st_uid=? and st_gid=? and deleted=0",(dirname,st.st_mode, st.st_uid,st.st_gid))
	try:
#		/* the directory already exists, twiddle the last seen */
		r=cursor.next()
		print "r.id=",r['id']
		cursor.execute('update dirs set last_seen=? where id=?',(datetime.datetime.now(),r['id']))
		print "found dir=",dirname
		return(r)
	except StopIteration:
#		/* either directory is missing or has different perms, in either case add it */
		cursor.execute('INSERT INTO dirs VALUES (null,?,?,?,?,?,?,?)',(st.st_mode,st.st_ino,st.st_uid,st.st_gid,dirname,datetime.datetime.now(),0))
		print "adding dir=",dirname
		return(1)

#/* dont bother to update last-seen if its just a parent lookup */
def getDirID(dirname):
	st=os.stat(dirname)
#	/* do we have a perfect match */
	cursor.execute("select id from dirs where name=? and st_mode=? and st_uid=? and st_gid=? and deleted=0",(dirname,st.st_mode, st.st_uid,st.st_gid))
	try:
#		/* the directory already exists, twiddle the last seen */
		return(cursor.next()[0])
	except StopIteration:
		print "this should never happen"
		return(nil)

def addFile(root,name):
	st=os.stat(os.path.join(root,name))
	dirID=getDirID(root)	
	cursor.execute("select * from files where name=? and st_gid=? and st_mtime=? and dirID=?",(name,st.st_gid,st.st_mtime,dirID))
	try:
#		/* we have an exact match (name, permissions, modtime, and name)
		r=cursor.next()
		cursor.execute('update files set last_seen=? where id=?',(datetime.datetime.now(),r['id']))
		print "found file=",name
		return(cursor.next)
	except StopIteration:
#		/* file is new, add it */
		cursor.execute('INSERT INTO files VALUES (null,?,?,?,?,?,?,?,?,?,?,?,?,?,?)',(st.st_mode,st.st_ino,st.st_dev,st.st_nlink,st.st_uid,st.st_gid,st.st_size,st.st_atime,st.st_mtime,st.st_ctime,name,dirID,datetime.datetime.now(),0))
		return(1)

print "backup started at ",startBackup
 
top="/home/bill/imp/src/bdr"
addDir(top)	
for root, dirs, files in os.walk(top, topdown=True):
#	/* For each dir, if we dont have a match on name, and st_mode, insert it */
	for name in dirs:
		addDir(os.path.join(root,name))
#		print "dir root=",root," name=",name

	connection.commit()
	for name in files:
		addFile(root,name)
#		print "file root=",root," name=",name

print "commiting database changes"
connection.commit()
print ""
print "executing select on dirs"
cursor.execute('SELECT * FROM dirs')
for row in cursor:
	print row[0],row[1],row[2],row[3],row[4],row[5]
	
print "executing select on files"
cursor.execute('SELECT * FROM files')
for row in cursor:
	print row[0],row[1],row[2],row[3],row[4],row[5],row[6],row[7],row[8],row[9],row[10],row[11],row[12],row[13],row[14]
	
