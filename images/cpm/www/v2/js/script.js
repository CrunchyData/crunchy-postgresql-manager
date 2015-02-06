// create the module and name it cpmApp
var cpmApp = angular.module('cpmApp', ['ngRoute', 'ngTable', 'ui.bootstrap', 'ngCookies', 'cpmApp.servers', 'cpmApp.clusters', 'cpmApp.containers', 'cpmApp.tools', 'cpmApp.settings']);

// configure our routes
cpmApp.config(function($routeProvider) {
    $routeProvider

        .when('/', {
        templateUrl: 'pages/home.html',
        //controller: 'mainController'
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
    });

});

// create the controller and inject Angular's $scope
cpmApp.controller('mainController', function($scope, $http, $cookieStore, $cookies, $modal, $cookieStore, $filter, ngTableParams) {
    // create a message to display in our view
    $scope.message = 'PostgreSQL Container Management!';
    if ($cookies.AdminURL) {} else {
        alert('CPM AdminURL setting is NOT defined, please update on the Settings page before using CPM');
    }


    $scope.hc = [];
    $scope.hcts;

    $scope.data = [];
    $scope.results = [];
    $scope.items = ['item1', 'item2'];
    $scope.loginValue = '';

    $scope.tableParams = new ngTableParams({
        page: 1, // show first page
        count: 10 // count per page
    }, {
        total: $scope.hc.length, // length of data
        getData: function($defer, params) {
            console.log('getData called hc=' + $scope.hc.length);
		$scope.getStatus();
            // use build-in angular filter
            var orderedData = $scope.hc;

            params.total(orderedData.length); // set total for recalc pagination
            $defer.resolve($scope.hc = orderedData.slice((params.page() - 1) * params.count(), params.page() * params.count()));
        }
    });

   //fix around ng-table bug?
    $scope.tableParams.settings().$scope = $scope;



    $scope.doLogin = function() {
        console.log(' login called');
        var modalInstance = $modal.open({
            size: 'sm',
            templateUrl: 'pages/loginmodal.html',
            controller: LoginController,
            resolve: {
                value: function() {
                    return $scope.loginValue;
                }
            }
        });
        modalInstance.result.then(function(token) {
            console.log('results token=' + token);
            $scope.loginValue = token;
            $cookieStore.put('cpmuser', token);
        }, function() {
            $log.info('Modal dismissed at: ' + new Date());
        });

    }
    $scope.doLogout = function() {
        var token = $cookieStore.get('cpmsession');
        console.log(' logout called for token ' + token);
        $http.get($cookies.AdminURL + '/sec/logout/' + token).success(function(data, status, headers, config) {

            console.log('logout ok');
            $cookieStore.remove('cpmsession')
            $scope.alerts = [{
                type: 'success',
                msg: 'Logout ok.'
            }];

        }).error(function(data, status, headers, config) {
            console.log('error:logout');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
        $scope.loginValue = '';
    }

function convertUTCDateToLocalDate(date) {
    var newDate = new Date(date.getTime()+date.getTimezoneOffset()*60*1000);

    var offset = date.getTimezoneOffset() / 60;
    var hours = date.getHours();

    newDate.setHours(hours - offset);

    return newDate;   
}
        $scope.getStatus = function() {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                alert('login required');
                return;
            }

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

	//$scope.getStatus();
});

var LoginController = function($http, $scope, $cookies, $cookieStore, $modalInstance, value) {
    $scope.value = value;
    $scope.ID = '';
    $scope.Password = '';
    $scope.results = [];

    console.log('LoginController called');
    $scope.ok = function() {
        console.log(' login ok called id=' + $scope.ID + ' psw=' + $scope.Password);
        $http.get($cookies.AdminURL + '/sec/login/' + $scope.ID + "." + $scope.Password).success(function(data, status, headers, config) {

            console.log('login ok');
            console.log('token=' + data.Contents);
            $cookieStore.put('cpmsession', data.Contents);
            $scope.value = data.Contents;
            $scope.alerts = [{
                type: 'success',
                msg: 'Login ok.'
            }];

            //$modalInstance.close(data.Contents);
            $modalInstance.close($scope.ID);
        }).error(function(data, status, headers, config) {
            console.log('error:login');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
    }
    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    }

};
