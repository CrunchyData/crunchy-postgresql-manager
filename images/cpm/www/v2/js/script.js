// create the module and name it cpmApp
var cpmApp = angular.module('cpmApp',
	['ngRoute',
	 'ngTable',
	 'ui.bootstrap',
	 'ngCookies',
	 'cpmApp.servers',
	 'cpmApp.clusters',
	 'cpmApp.containers',
	 'cpmApp.tools',
	 'cpmApp.settings'
	]
);

// configure our routes
cpmApp.config(function($routeProvider) {
    $routeProvider
        .when('/', {
        templateUrl: 'pages/home.html'
    })

    .when('/servers', {
        templateUrl: 'pages/servers.html',
        controller: 'serversController'
    })

    .when('/containers', {
        templateUrl: 'pages/containers.html',
        controller: 'containersController'
    })

    .when('/clusters', {
        templateUrl: 'pages/clusters.html',
        controller: 'clustersController'
    })

    .when('/tools', {
        templateUrl: 'pages/tools.html',
        controller: 'toolsController'
    })

    .when('/settings', {
        templateUrl: 'pages/settings.html',
        controller: 'settingsController'
    })

	.when('/login', {
		templateUrl: 'pages/login.html',
		controller: 'LoginController'
	})

	.otherwise({ redirectTo: '/' });
});

cpmApp.service('Session', function() {
	this.create = function(user_id, token) {
		this.user_id = user_id;
		this.token = token;
	};

	this.destroy = function() {
		this.user_id = null;
		this.token = null;
	}

	return this;
});

cpmApp.run( function($rootScope, $location, $cookieStore) {
	$rootScope.$on( "$locationChangeStart", function(event, next, current) {
		if ($cookieStore.get('cpm_token') == null) {
			$location.path( "/login" );
		}
	});
});

// create the controller and inject Angular's $scope
cpmApp.controller('mainController',
	[
	 "$scope",
	 "$http",
	 "$location",
	 "$cookieStore",
	 "$cookies",
	 "ngTableParams",
	function($scope, $http, $location, $cookieStore, $cookies, ngTableParams) {

	// create a message to display in our view
	$scope.message = 'PostgreSQL Container Management!';

	$scope.hc = [];
	$scope.hcts;

	$scope.data = [];
	$scope.results = [];
	$scope.items = ['item1', 'item2'];

	console.log('User Id: ' + $scope.cpm_user_id);

	$scope.doLogout = function() {
		var token = $cookieStore.get('cpm_token');

		$http.get($cookieStore.get('AdminURL') + '/sec/logout/' + token).
			success(function(data, status, headers, config) {
				console.log('logout ok');
				$cookieStore.remove('cpm_token');
				$scope.cpm_user_id = '';
				$location.path('/login');
				
				$scope.alerts = [{
					type: 'success',
					msg: 'Successfully logged out.'
				}];
			}).
			error(function(data, status, headers, config) {
				console.log('error:logout');
				$scope.alerts = [{
					type: 'danger',
					msg: data.Error
				}];
			});
	}


function convertUTCDateToLocalDate(date) {
    var newDate = new Date(date.getTime()+date.getTimezoneOffset()*60*1000);

    var offset = date.getTimezoneOffset() / 60;
    var hours = date.getHours();

    newDate.setHours(hours - offset);

    return newDate;   
}

	$scope.getStatus = function() {
		var token = $cookieStore.get('cpm_token');

		console.log('getStatus called ');
		var queryroot = 'http://cpm-mon.crunchy.lab:8086/db/cpm/series?u=root&p=root&q=';
		var query1 = 'select seconds, service, servicetype, status from hc1 limit 1';
		console.log(queryroot + query1);

		$http.get(queryroot + query1).
			success(function(data, status, headers, config) {
				console.log('hc1 first row results ' + data[0].points[0]);
				console.log('hc1 first row ts ' + data[0].points[0][2]);
				firstRow = data;

				query2 = 'select seconds, service, servicetype, status from hc1 where seconds = ' + data[0].points[0][2];
				console.log(queryroot + query2);

				$http.get(queryroot + query2).
					success(function(data, status, headers, config) {
						console.log('hc1 full results ' + data[0].points);
						$scope.hc = data[0].points;
						var date = new Date(null);
						date.setSeconds(data[0].points[0][2]);
						//$scope.hcts = date.toISOString();
						$scope.hcts = convertUTCDateToLocalDate(date).toUTCString();
						//overlay for tooltip hover	
						for (i = 0; i < $scope.hc.length; i++) {
							$scope.hc[i][2] = 'database down';
							$scope.hc[i][4] = 'Database -' + $scope.hc[i][3];
						}
					}).error(function(data, status, headers, config) {
						alert('error 1');
					});

			}).error(function(data, status, headers, config) {
				alert('error 2');
			});
	};
}]);

var LoginController = function($http, $rootScope, $scope, $cookieStore, $location, Session) {
	$scope.submit = function() {
		var user_id = $scope.login.user_id;
		var password = $scope.login.password;
		var admin_url = $scope.login.admin_url;
		
		var loginUrl = admin_url + '/sec/login/' + user_id + "." + password;

		$http.get(loginUrl).
			success(function(data, status, headers, config) {
				$rootScope.cpm_user_id = user_id;
				$cookieStore.put('AdminURL', admin_url);
				$cookieStore.put('cpm_token', data.Contents);
				$location.path('/home');
			}).
			error(function(data, status, headers, config) {
				$scope.alerts = [{
					type: 'danger',
					msg: data.Error
				}];
			});
	};
};

