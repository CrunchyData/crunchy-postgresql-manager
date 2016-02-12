
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
		github.com/crunchydata/crunchy-postgresql-manager/sec \
		github.com/crunchydata/crunchy-postgresql-manager/template \
		github.com/crunchydata/crunchy-postgresql-manager/admindb 

build:
		cd adminserver && make
		cd cpmserver && make
		cd cpmcontainerserver && make
		cd backupcommand && make
		cd taskserver && make
		cd collectserver && make
		cd restorecommand && make
		cd backrestrestorecommand && make

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
		cd images/cpm-backrest-restore-job && make
		cd images/cpm-prometheus && make
		cd images/cpm-collect && make
		cd images/cpm-efk && make

start:
		./sbin/start-cpm.sh

stop:
		./sbin/stop-cpm.sh

clean:
		rm -rf $(GOBIN)/*server* $(GOBIN)/*command*
		rm -rf $(GOPATH)/pkg/linux_amd64/github.com/crunchydata/crunchy-postgresql-manager/*.a


