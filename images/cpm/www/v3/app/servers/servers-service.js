angular.module('uiRouterSample.servers.service', ['ngCookies'])

// A RESTful factory for retrieving servers from 'servers.json'
/**
.factory('servers', ['$http', 'utils', function ($http, utils) {
  var path = 'assets/servers.json';
  var servers = $http.get(path).then(function (resp) {
    return resp.data.servers;
  });

  var factory = {};
  factory.all = function () {
    return servers;
  };
  factory.get = function (id) {
    return servers.then(function(){
      return utils.findById(servers, id);
    })
  };
  return factory;
}])
*/

.factory('serversFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var serversFactory = {};

    serversFactory.all = function() {
        var url = $cookieStore.get('AdminURL') + '/servers/' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    serversFactory.get = function(serverid) {

        var cleanip = serverid.replace(/\./g, "_");
        var url = $cookieStore.get('AdminURL') + '/server/' + cleanip + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    serversFactory.startall = function(serverid) {

        var cleanip = serverid.replace(/\./g, "_");
        var url = $cookieStore.get('AdminURL') + '/admin/startall/' + cleanip + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };
    serversFactory.stopall = function(serverid) {

        var cleanip = serverid.replace(/\./g, "_");
        var url = $cookieStore.get('AdminURL') + '/admin/stopall/' + cleanip + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    serversFactory.iostat = function(serverid) {

        var cleanip = serverid.replace(/\./g, "_");
        var url = $cookieStore.get('AdminURL') + '/monitor/server-getinfo/' + cleanip + '.cpmiostat.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };
    serversFactory.df = function(serverid) {

        var cleanip = serverid.replace(/\./g, "_");
        var url = $cookieStore.get('AdminURL') + '/monitor/server-getinfo/' + cleanip + '.cpmdf.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };
    serversFactory.containers = function(serverid) {

        var cleanip = serverid.replace(/\./g, "_");
        var url = $cookieStore.get('AdminURL') + '/nodes/forserver/' + cleanip + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    return serversFactory;
}]);
