#!/usr/bin/env python

import sys
import gzip

fd = gzip.open(sys.argv[1], 'r')

ln = fd.readline()
while ln:
    sys.stdout.write(ln)
    ln = fd.readline()

