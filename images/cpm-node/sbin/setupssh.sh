#!/bin/expect 

set timeout 2

set ip [lindex $argv 0]
set user [lindex $argv 1]

spawn ssh -o StrictHostKeyChecking=no "$user\@$ip"
expect "Last Login:"
send "exit\r"
interact

