angular.module('uiRouterSample.proxy.service', ['ngCookies'])

.factory('proxyFactory', ['$rootScope', '$http', '$cookieStore', 'utils', function($rootScope, $http, $cookieStore, utils) {

    var proxyFactory = {};

    proxyFactory.all = function() {
    	console.log('in proxy all with projectId=' + $rootScope.projectId);
        var url = $cookieStore.get('AdminURL') + '/projectnodes/' + $rootScope.projectId + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    proxyFactory.getbycontainerid = function(containerid) {

        var url = $cookieStore.get('AdminURL') + '/proxy/getbycontainerid/' + containerid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    proxyFactory.delete = function(id) {

        var url = $cookieStore.get('AdminURL') + '/deleteproxy/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    proxyFactory.add = function(proxy, container, standalone, dockerprofile) {

	console.log('here with proxy host=' + proxy.DatabaseHost);

        var url = $cookieStore.get('AdminURL') + '/provisionproxy/' +
            'SM' + '.' +
            'cpm-node-proxy' + '.' +
            container.ServerID + '.' +
            container.ProjectID + '.' +
            container.Name + '.' +
            'true' + '.' +
            $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

 	return proxyFactory;

}]);
