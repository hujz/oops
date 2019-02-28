lscpu
lsmem -b -o SIZE -r | awk 'BEGIN{sum=0} {if(NR != 1){sum+=$1;}} END{print sum "/1024/1024/1024"}' | bc
ip a
