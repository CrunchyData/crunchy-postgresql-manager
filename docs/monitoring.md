Monitoring
================
http://cpm-promdash:3000/servermemdashboard

http://cpm-prometheus:9090

### containers
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
cpm-prometheus

cpm-collect
collects metrics for servers and databases and then "sets" them
so that prometheus can pull them.


### collected metrics
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
server cpu
server mem

### embedded graphs

<iframe height="400px" width="100%" 
        ng-src="{{servergraphlink}}">
</iframe>

in javascript:

$scope.servergraphlink=$sce.trustAsResourceUrl('http://cpm-promdash:3000/embed/servermemdashboard#!?var.host=' + $scope.server.Name);

Notice we use a template variable to pass in the Server name which is
a way to uniquely identify the metric as belonging to a particular server.

The 'servermemdashboard' is defined in promdash (cpm-promdash) and saved
as a dashboard.


