
gendeps:
		godep save \
		github.com/crunchydata/crunchy-postgresql-manager/types \
		github.com/crunchydata/crunchy-postgresql-manager/collect \
		github.com/crunchydata/crunchy-postgresql-manager/util \
		github.com/crunchydata/crunchy-postgresql-manager/logit \
		github.com/crunchydata/crunchy-postgresql-manager/adminapi \
		github.com/crunchydata/crunchy-postgresql-manager/cpmserverapi \
		github.com/crunchydata/crunchy-postgresql-manager/cpmcontainerapi \
		github.com/crunchydata/crunchy-postgresql-manager/task \
		github.com/crunchydata/crunchy-postgresql-manager/dummy \
		github.com/crunchydata/crunchy-postgresql-manager/sec \
		github.com/crunchydata/crunchy-postgresql-manager/template \
		github.com/crunchydata/crunchy-postgresql-manager/admindb 

build:
		godep go install cmd/adminapi.go
		godep go install cmd/cpmserverapi.go
		godep go install cmd/cpmcontainerapi.go
		godep go install cmd/backupcommand.go
		godep go install cmd/taskserver.go
		godep go install cmd/dummyserver.go
		godep go install cmd/dockerapi.go
		godep go install cmd/collectserver.go
		godep go install cmd/restorecommand.go

buildimages:
		cd images/cpm-server && make  
		cd images/cpm-admin && make  
		cd images/cpm && make 
		cd images/cpm-node && make  
		cd images/cpm-node-proxy && make  
		cd images/cpm-pgpool && make
		cd images/cpm-task && make
		cd images/cpm-backup-job && make
		cd images/cpm-restore-job && make
		cd images/cpm-prometheus && make
		cd images/cpm-collect && make

clean:
		rm -rf $(GOBIN)/*


