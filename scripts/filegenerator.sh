for i in `seq 1 4`; do dd if=/dev/urandom of=test$i count=16384 bs=16384; done
