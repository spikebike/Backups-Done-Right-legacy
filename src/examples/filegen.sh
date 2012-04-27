 in `seq 1 8`; do dd if=/dev/urandom of=test$i count=8192 bs=16384; done
