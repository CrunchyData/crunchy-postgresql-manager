
gendeps:
		godep save \
		github.com/crunchydata/crunchy-postgresql-manager/collect \
		github.com/crunchydata/crunchy-postgresql-manager/util \
		github.com/crunchydata/crunchy-postgresql-manager/logit \
		github.com/crunchydata/crunchy-postgresql-manager/adminapi \
		github.com/crunchydata/crunchy-postgresql-manager/backup \
		github.com/crunchydata/crunchy-postgresql-manager/cpmnodeagent \
		github.com/crunchydata/crunchy-postgresql-manager/cpmserveragent \
		github.com/crunchydata/crunchy-postgresql-manager/dummy \
		github.com/crunchydata/crunchy-postgresql-manager/kubeclient \
		github.com/crunchydata/crunchy-postgresql-manager/myinfluxdb/client \
		github.com/crunchydata/crunchy-postgresql-manager/mon \
		github.com/crunchydata/crunchy-postgresql-manager/sec \
		github.com/crunchydata/crunchy-postgresql-manager/template \
		github.com/crunchydata/crunchy-postgresql-manager/admindb 

build:
		godep go install cmd/adminapi.go
		godep go install cmd/cpmnodeagent.go
		godep go install cmd/cpmserveragent.go
		godep go install cmd/backupcommand.go
		godep go install cmd/backupserver.go
		godep go install cmd/monserver.go
		godep go install cmd/dummyserver.go
		godep go install cmd/dockerapi.go
		godep go install cmd/collectserver.go

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


