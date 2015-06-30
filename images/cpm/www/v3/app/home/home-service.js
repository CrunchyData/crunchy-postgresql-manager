angular.module('uiRouterSample.home.service', ['ngCookies'])

// A RESTful factory for retrieving home from 'home.json'
.factory('homeFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, utils) {

    var factory = {};

    factory.healthcheck = function() {
     	var url = $cookieStore.get('AdminURL') + '/mon/healthcheck/' + $cookieStore.get('cpm_token');
       	console.log(url);
	return $http.get(url);
    };

    return factory;
}]);
