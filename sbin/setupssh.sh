#!/bin/expect 

set timeout 2

set ip [lindex $argv 0]
set user [lindex $argv 1]
set password [lindex $argv 2]

spawn ssh "$user\@$ip"
expect "password:"
send "$password\r";
expect "Last Login:"
send "hostname\r"
send "exit\r"
interact

