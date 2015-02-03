// create the module and name it cpmApp
var cpmApp = angular.module('cpmApp', ['ngRoute', 'ngTable', 'ui.bootstrap', 'ngCookies', 'cpmApp.servers', 'cpmApp.clusters', 'cpmApp.containers', 'cpmApp.tools', 'cpmApp.settings']);

// configure our routes
cpmApp.config(function($routeProvider) {
    $routeProvider

    // route for the home page
        .when('/', {
        templateUrl: 'pages/home.html',
        controller: 'mainController'
    })

    // route for the about page
    .when('/servers', {
        templateUrl: 'pages/servers.html',
        controller: 'serversController'
    })

    // route for the servers page
    .when('/containers', {
        templateUrl: 'pages/containers.html',
        controller: 'containersController'
    })

    // route for the contact page
    .when('/clusters', {
        templateUrl: 'pages/clusters.html',
        controller: 'clustersController'
    })

    // route for the contact page
    .when('/tools', {
        templateUrl: 'pages/tools.html',
        controller: 'toolsController'
    })

    // route for the settings page
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

    $scope.hc = [
{
'ts': '2015-Feb-2 10am',
'rows':[
 [  'Database - two',
    'Database is down',
    'databasedown'
 ],
 [  'Database - one',
    'Database is down',
    'databasedown'
 ]
]
}	];
console.log('hc.ts=' + $scope.hc[0].ts);

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
