pid=`/opt/software/jdk/jdk1.8.0_202/bin/jps | awk '{if($2=="Main"){print $1}}' `
ps -L -q $pid  -o 'ppid,pid,lwp,cpu,time,rss,%cpu,cmd' | wc -l

echo ok