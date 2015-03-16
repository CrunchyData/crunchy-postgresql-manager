#curl -X POST -H "Content-Type: application/json" \
#--data @cpm-node-template.json  \
#--cacert /home/jeffmc/originmaster/origin/openshift.local.certificates/admin/root.crt \
#https://127.0.0.1:8443/api/v1beta1/pods

#curl -L -X POST -H "Content-Type: application/json" \
#--data @cpm-node-template.json  \
#--cacert /home/jeffmc/originmaster/origin/openshift.local.certificates/admin/root.crt \
#https://127.0.0.1:8443/api/v1beta1/pods

export ROOT=/home/jeffmc/originmaster/origin/openshift.local.certificates
export THING=admin
export KUBECONFIG=$ROOT/$THING/.kubeconfig
curl -X GET https://admin:admin@127.0.0.1:8443/api/v1beta1/pods \
--cacert $ROOT/$THING/cert.crt --insecure

export THING=kube-client
export KUBECONFIG=$ROOT/$THING/.kubeconfig
curl -X GET https://admin:admin@127.0.0.1:8443/api/v1beta1/pods \
--cacert $ROOT/$THING/cert.crt --insecure
