
gendeps:
		godep save \
		github.com/jmccormick2001/crunchy-postgresql-manager/util \
		github.com/jmccormick2001/crunchy-postgresql-manager/logit \
		github.com/jmccormick2001/crunchy-postgresql-manager/adminapi \
		github.com/jmccormick2001/crunchy-postgresql-manager/backup \
		github.com/jmccormick2001/crunchy-postgresql-manager/cpmagent \
		github.com/jmccormick2001/crunchy-postgresql-manager/dummy \
		github.com/jmccormick2001/crunchy-postgresql-manager/kubeclient \
		github.com/jmccormick2001/crunchy-postgresql-manager/myinfluxdb/client \
		github.com/jmccormick2001/crunchy-postgresql-manager/mon \
		github.com/jmccormick2001/crunchy-postgresql-manager/sec \
		github.com/jmccormick2001/crunchy-postgresql-manager/template \
		github.com/jmccormick2001/crunchy-postgresql-manager/admindb 

build:
		godep go install cmd/adminapi.go
		godep go install cmd/backupcommand.go
		godep go install cmd/backupserver.go
		godep go install cmd/cpmagentserver.go
		godep go install cmd/monserver.go
		godep go install cmd/dummyserver.go

buildimages:
		cd images/cpm-dashboard && make
		cd images/cpm-base && make  
		cd images/cpm-admin && make  
		cd images/cpm && make 
		cd images/cpm-node && make  
		cd images/cpm-pgpool && make
		cd images/cpm-backup && make
		cd images/cpm-backup-job && make
		cd images/cpm-mon && make

clean:
		rm -rf $(GOBIN)/*


