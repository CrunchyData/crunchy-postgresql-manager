BEGIN {
  FS=","
  print "{\"iostat\":["
}
 
# this is a cool trick to print comma except last line
NR > 1 {
  print ","
}
 
{
  print "{\"device\":\"" $1 "\",\"tps\":\"" $2 "\",\"kbreadpersec\":\"" $3 "\",\"kbwritespersec\":\"" $4 "\",\"kbreads\":\"" $5 "\",\"kbwrites\":\"" $6 "\"}"
}
 
END {
  print "]}"
}
