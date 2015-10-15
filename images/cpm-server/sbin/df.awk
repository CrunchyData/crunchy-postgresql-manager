BEGIN {
  FS=","
  print "{\"df\":["
}
 
# this is a cool trick to print comma except last line
NR > 1 {
  print ","
}
 
{
  print "{\"filesystem\":\"" $1 "\",\"total\":\"" $2 "\",\"used\":\"" $3 "\",\"available\":\"" $4 "\",\"pctused\":\"" $5 "\",\"mountpt\":\"" $6 "\"}"
}
 
END {
  print "]}"
}
