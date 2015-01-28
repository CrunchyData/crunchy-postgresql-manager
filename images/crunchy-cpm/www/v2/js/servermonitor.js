(function() {
    var app = angular.module('TabsApp', ['ui.bootstrap', 'ngRoute', 'ngCookies']);

    // configure our routes
    app.config(function($routeProvider) {
        $routeProvider
        // route for the home page
            .when('/', {
            templateUrl: 'server-monitor-iostat.html',
            controller: 'iostatController'
        })

        .when('/df', {
            templateUrl: 'server-monitor-df.html',
            controller: 'dfController'
        })
        .when('/graph', {
            templateUrl: 'server-graph.html',
            controller: 'graphController'
        })
    });

    app.controller('TabsCtrl', function($scope, $route, $http, $cookies, $cookieStore) {
        $scope.server = [];
        $scope.dfresults = [];
        $scope.iostatresults = [];
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            alert('login required');
            return;
        }

	$scope.currentUser = [];
	$scope.currentUser = $cookieStore.get('cpmuser');
	console.log('currentUser is ' + $scope.currentUser);

        $http.get($cookies.AdminURL + '/server/' + window.serverid + '.' + token).success(function(data, status, headers, config) {
            $scope.server = data;
        }).error(function(data, status, headers, config) {
            alert('error in get server');
        });

        console.log('working on iostat');
        $http.get($cookies.AdminURL + '/monitor/server-getinfo/' + serverid + ".cpmiostat." + token).
        success(function(data, status, headers, config) {
            $scope.iostatresults = data.iostat;
            console.log('getinfo results set ' + data);
        }).
        error(function(data, status, headers, config) {
            alert('error happended');
        });

    });

    app.controller('iostatController', function($rootScope, $scope, $route, $http, $cookies, $cookieStore) {


        $scope.handleRefresh = function() {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                alert('login required');
                return;
            }
            console.log('working on iostat');
            $http.get($cookies.AdminURL + '/monitor/server-getinfo/' + serverid + ".cpmiostat." + token).
            success(function(data, status, headers, config) {
                $scope.iostatresults = data.iostat;
                console.log('getinfo results set ' + data);
            }).
            error(function(data, status, headers, config) {
                alert('error happended');
            });
        };

        $scope.handleRefresh();
    });

    app.controller('dfController', function($rootScope, $scope, $route, $http, $cookies, $cookieStore) {


        $scope.handleRefresh = function() {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                alert('login required');
                return;
            }
            console.log('working on df');
            console.log('calling getinfo service');
            $http.get($cookies.AdminURL + '/monitor/server-getinfo/' + serverid + ".cpmdf." + token).
            success(function(data, status, headers, config) {
                $scope.dfresults = data.df;
                console.log('getinfo results set ' + data.df);
            }).error(function(data, status, headers, config) {
                alert('error happended');
            });

        };

        $scope.handleRefresh();
    });

    app.controller('graphController', function($rootScope, $scope, $route, $http, $cookies, $cookieStore) {


        $scope.handleRefresh = function() {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                alert('login required');
                return;
            }
            console.log('graphing server stats');
            $http.get($cookies.AdminURL + '/monitor/server-getinfo/' + serverid + ".cpmdf." + token).
            success(function(data, status, headers, config) {
                $scope.dfresults = data.df;
                console.log('getinfo results set ' + data.df);
            }).error(function(data, status, headers, config) {
                alert('error happended');
            });

        };

        $scope.handleRefresh();
    });


})();
