angular.module('uiRouterSample.proxy.service', ['ngCookies'])

.factory('proxyFactory', ['$rootScope', '$http', '$cookieStore', 'utils', function($rootScope, $http, $cookieStore, utils) {

    var proxyFactory = {};

    proxyFactory.all = function() {
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

        var url = $cookieStore.get('AdminURL') + '/provisionproxy';
        console.log(url);

 	return $http.post(url, {
            'Profile': dockerprofile,
            'Image': 'cpm-node-proxy',
            'ServerID': container.ServerID,
            'ProjectID': container.ProjectID,
            'ContainerName': container.Name,
            'Standalone': 'false',
            'Host': proxy.DatabaseHost,
            'Usename': proxy.UserID,
            'Passwd': proxy.UserPassword,
            'Database': proxy.DatabaseName,
            'Port': proxy.DatabasePort,
            'Token': $cookieStore.get('cpm_token')

        });
	};

    proxyFactory.update = function(proxy) {

	console.log("in update service with proxy....");
	console.log(JSON.stringify(proxy));
        var url = $cookieStore.get('AdminURL') + '/proxy/update';
        console.log(url);

 	return $http.post(url, { 
		'ID' : proxy.ID,
		'Port' : proxy.Port,
		'Host' : proxy.Host,
		'Database' : proxy.Database,
            	'Token': $cookieStore.get('cpm_token')
	});
 	};

 	return proxyFactory;

}]);
