package main

import (
	"fmt"
	"strings"
)

func main() {
	var output = `Catalog version number:               201409291
	Database system identifier:           6119948401140838434
	Database cluster state:               in production
	pg_control last modified:             Fri Feb 27 18:11:13 2015
	Latest checkpoint location:           0/17521A8
	Prior checkpoint location:            0/1752140
	Latest checkpoint's REDO location:    0/17521A8
	Latest checkpoint's REDO WAL file:    000000010000000000000001
	Latest checkpoint's TimeLineID:       1
	Latest checkpoint's PrevTimeLineID:   1
	Latest checkpoint's full_page_writes: on
	Latest checkpoint's NextXID:          0/776
	Latest checkpoint's NextOID:          24576
	Latest checkpoint's NextMultiXactId:  1
	Latest checkpoint's NextMultiOffset:  0
	Latest checkpoint's oldestXID:        755
	Latest checkpoint's oldestXID's DB:   1
	Latest checkpoint's oldestActiveXID:  0
	Latest checkpoint's oldestMultiXid:   1
	Latest checkpoint's oldestMulti's DB: 1
	Time of latest checkpoint:            Fri Feb 27 18:11:13 2015
	Fake LSN counter for unlogged rels:   0/1
	Minimum recovery ending location:     0/0
	Min recovery ending loc's timeline:   0
	Backup start location:                0/0
	Backup end location:                  0/0
	End-of-backup record required:        no
	Current wal_level setting:            archive
	Current wal_log_hints setting:        off
	Current max_connections setting:      100
	Current max_worker_processes setting: 8
	Current max_prepared_xacts setting:   0
	Current max_locks_per_xact setting:   64
	Maximum data alignment:               8
	Database block size:                  8192
	Blocks per segment of large relation: 131072
	WAL block size:                       8192
	Bytes per WAL segment:                16777216
	Maximum length of identifiers:        64
	Maximum columns in an index:          32
	Maximum size of a TOAST chunk:        1996
	Size of a large-object chunk:         2048
	Date/time type storage:               64-bit integers
	Float4 argument passing:              by value
	Float8 argument passing:              by value
	Data page checksum version:           0
	`

	fmt.Println("hi")
	lines := strings.Split(output, "\n")
	fmt.Println(len(lines))
	for i := range lines {
		fmt.Println(len(lines[i]))
		if len(lines[i]) > 1 {
			columns := strings.Split(lines[i], ":")
			fmt.Println("name=[" + strings.TrimSpace(columns[0]) + "] value=[" + strings.TrimSpace(columns[1]) + "]")
		}
	}
}
