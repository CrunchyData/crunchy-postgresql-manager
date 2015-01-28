crunchy-pgpool deployment in Kuber
===========================================

1. Start an OpenShift all-in-one server

        openshift start

2. Use the command line to transform the template, and then send each object to the server:

        openshift kube process -c pgpool-template.json | openshift kube apply -c -

   Note: `-c -` tells the CLI to read a file from STDIN - you can use this in other places as well.

Alternatively, using the Openshift 'deployments' concept, but as of
now, this doesn't work it appears?:

X. You can deploy pgpool-node-config.json with:

	$ openshift kube apply -c crunchy-pgpool-config.json

