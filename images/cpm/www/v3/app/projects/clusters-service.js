angular.module('uiRouterSample.clusters.service', ['ngCookies'])

.factory('clustersFactory', ['$rootScope', '$http', '$cookieStore', 'utils', function($rootScope, $http, $cookieStore, utils) {

    var clustersFactory = {};

    clustersFactory.all = function() {
    	console.log('in clusters all with projectId = ' + $rootScope.projectId);
        var url = $cookieStore.get('AdminURL') + '/projectclusters/' + $rootScope.projectId + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    clustersFactory.get = function(id) {

        var url = $cookieStore.get('AdminURL') + '/cluster/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    //get list of containers in a cluster
    clustersFactory.getcontainers = function(id) {

        var url = $cookieStore.get('AdminURL') + '/clusternodes/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    clustersFactory.delete = function(id) {

        var url = $cookieStore.get('AdminURL') + '/cluster/delete/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    clustersFactory.scale = function(id) {

        var url = $cookieStore.get('AdminURL') + '/cluster/scale/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    clustersFactory.failover = function(containerid) {

        var url = $cookieStore.get('AdminURL') + '/admin/failover/' + containerid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    clustersFactory.configure = function(id) {

        var url = $cookieStore.get('AdminURL') + '/cluster/configure/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    clustersFactory.join = function(names, masterID, clusterID) {

        var url = $cookieStore.get('AdminURL') + '/event/join-cluster/' + names + '.' + masterID + '.' + clusterID + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    //get list of containers not in a cluster
    clustersFactory.nocluster = function() {

        var url = $cookieStore.get('AdminURL') + '/nodes/nocluster/' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    clustersFactory.add = function(cluster) {

        var url = $cookieStore.get('AdminURL') + '/cluster/';
        console.log(url);

        return $http.post(url, {
            'Name': cluster.Name,
            'ProjectID': $rootScope.projectId,
            'Status': 'uninitialized',
            'ClusterType': cluster.ClusterType,
            'Token': $cookieStore.get('cpm_token')
        });
    };

    clustersFactory.autocluster = function(cluster, profile) {

        var url = $cookieStore.get('AdminURL') + '/autocluster';
        console.log(url);

        return $http.post(url, {
            'Name': cluster.Name,
            'ProjectID': cluster.ProjectID, 
            'ClusterType': cluster.ClusterType,
            'ClusterProfile': profile,
            'Token': $cookieStore.get('cpm_token')
        });
    };

    return clustersFactory;
}]);
