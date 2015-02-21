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
    })

    .when('/containers', {
        templateUrl: 'pages/containers.html',
    })

    .when('/clusters', {
        templateUrl: 'pages/clusters.html',
    })

    .when('/tools', {
        templateUrl: 'pages/tools.html',
    })

    .when('/settings', {
        templateUrl: 'pages/settings.html',
    })

	.when('/login', {
		templateUrl: 'pages/login.html',
	})

	.otherwise({ redirectTo: '/' });
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
	 "$rootScope",
	 "$http",
	 "$location",
	 "$cookieStore",
	 "$cookies",
	 "ngTableParams",
	function($scope, $rootScope, $http, $location, $cookieStore, $cookies, ngTableParams) {

	// create a message to display in our view
	$scope.message = 'PostgreSQL Container Management!';

	$scope.hc = [];
	$scope.hcts;

	$scope.data = [];
	$scope.results = [];
	$scope.items = ['item1', 'item2'];

	$rootScope.cpm_user_id = $cookieStore.get('cpm_user_id');

	$scope.doLogout = function() {
		var token = $cookieStore.get('cpm_token');

		$http.get($cookieStore.get('AdminURL') + '/sec/logout/' + token).
			success(function(data, status, headers, config) {
				console.log('logout ok');
				$cookieStore.remove('cpm_token');
				$rootScope.cpm_user_id = '';
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
	    var query = $cookies.AdminURL + '/mon/hc1/' + token;

            $http.get(query).
            success(function(data, status, headers, config) {
                console.log('hc1: first row results ' + data[0].points[0]);
                console.log('hc1: first row ts ' + data[0].points[0][2]);
		console.log('hc1: full results ' + data[0].points);
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
		    console.log(data);
                alert('error 2');
            });

        };

}]);

var LoginController = function($http, $rootScope, $scope, $cookieStore, $location) {

	$scope.admin_url = $cookieStore.get('AdminURL');

	$scope.submit = function() {
		var user_id = $scope.user_id;
		var password = $scope.password;
		var admin_url = $scope.admin_url;

		var loginUrl = admin_url + '/sec/login/' + user_id + "." + password;

		$http.get(loginUrl).
			success(function(data, status, headers, config) {
				$rootScope.cpm_user_id = user_id;
				$cookieStore.put('cpm_user_id', user_id);
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

