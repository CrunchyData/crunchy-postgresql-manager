(function() {
    var app = angular.module('TabsApp', ['ngCookies', 'ui.bootstrap', 'ngRoute', 'ui.slider']);


    // configure our routes
    app.config(function($routeProvider) {
        $routeProvider
        // route for the home page
            .when('/', {
            templateUrl: 'monitor-db-stats.html',
            controller: 'statsController'
        })

        // route for the about page
        .when('/bgwriter', {
            templateUrl: 'monitor-db-bgwriter.html',
            controller: 'bgwriterController'
        })

        // route for the servers page
        .when('/repl', {
            templateUrl: 'monitor-db-repl.html',
            controller: 'replController'
        })

        // route for the contact page
        .when('/loadtest', {
            templateUrl: 'monitor-db-loadtest.html',
            controller: 'loadtestController'
        })

        // route for the monitor 1 page
        .when('/monitor1', {
            templateUrl: 'monitor-db-graph.html',
            controller: 'monitor1Controller'
        })
    });

    app.controller('TabsCtrl', function($rootScope, $scope, $route, $http, $cookies, $cookieStore) {
        $scope.status = {
            isFirstOpen: true,
            isFirstDisabled: false
        };

	$scope.currentUser = [];
	$scope.currentUser = $cookieStore.get('cpmuser');
	console.log('currentUser is ' + $scope.currentUser);

        $scope.oneAtATime = true;
        $scope.currenturl = [];
        $scope.container = [];
        $scope.statdbresults = [];
        $scope.statreplresults = [];
        $scope.bgwriterresults = [];
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }

        $scope.containerid = window.containerid;
        $http.get($cookies.AdminURL + '/node/' + window.containerid + '.' + token).success(function(data, status, headers, config) {
            $scope.container = data;
        }).error(function(data, status, headers, config) {
            alert('error in get container');
        });

    });

    app.controller('loadtestController', function($rootScope, $scope, $route, $http, $cookies, $cookieStore) {
        $scope.slidervaluehigh = "10000";
        $scope.slidervaluelow = "1000";
        $scope.slidervalue = "1000";

        $scope.handleRefresh = function() {
            $scope.isLoading = true;
            console.log('refresh pressed slider=' + $scope.slidervalue);
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                console.log('cookie was undefined');
                alert('login required');
                return;
            }

            var thing4 = $cookies.AdminURL + '/monitor/container-loadtest/' + $scope.containerid + '.loadtest.' + $scope.slidervalue + "." + token;
            $http.get(thing4).success(function(data, status, headers, config) {
                console.log('got loadtest results');
                console.log('now=' + data);
                $scope.loadtestresults = data;
                $scope.isLoading = false;
            }).error(function(data, status, headers, config) {
                alert('error in loadtest call ');
            });
        };

        $scope.handleRefresh();
    });

    app.controller('replController', function($rootScope, $scope, $route, $http, $cookies, $cookieStore) {
        $scope.handleRefresh = function() {
            console.log('working on repl');

            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                console.log('cookie was undefined');
                alert('login required');
                return;
            }

            var thing2 = $cookies.AdminURL + '/monitor/container-getinfo/' + $scope.containerid + '.statreplication.' + token;
            console.log('url=' + thing2);
            $http.get(thing2).success(function(data, status, headers, config) {
                console.log('got statrepl results');
                console.log('pid=' + data[0].pid);
                $scope.statreplresults = data;
            }).error(function(data, status, headers, config) {
                alert('error in monitor container statrepl');
            });
        };


        $scope.handleRefresh();
    });

    app.controller('bgwriterController', function($rootScope, $scope, $route, $http, $cookies, $cookieStore) {
        $scope.handleRefresh = function() {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                console.log('cookie was undefined');
                alert('login required');
                return;
            }
            console.log('working on bgwriter');
            var thing3 = $cookies.AdminURL + '/monitor/container-getinfo/' + $scope.containerid + '.bgwriter.' + token;
            console.log('url=' + thing3);
            $http.get(thing3).success(function(data, status, headers, config) {
                console.log('got bgwriter results');
                console.log('now=' + data.now);
                $scope.bgwriterresults = data;
            }).error(function(data, status, headers, config) {
                alert('error in monitor container bgwriter');
            });

        };
        $scope.handleRefresh();
    });

    app.controller('statsController', function($rootScope, $scope, $route, $http, $cookies, $cookieStore) {
        $scope.handleRefresh = function() {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                console.log('cookie was undefined');
                alert('login required');
                return;
            }
            var thing = $cookies.AdminURL + '/monitor/container-getinfo/' + $scope.containerid + '.statdatabase.' + token;
            console.log('url=' + thing);
            $http.get(thing).success(function(data, status, headers, config) {
                console.log('got statdb results');
                $scope.statdbresults = data;
            }).error(function(data, status, headers, config) {
                alert('error in monitor container statdb');
            });
        };

        $scope.handleRefresh();
    });

    app.controller('monitor1Controller', function($rootScope, $scope, $route, $http, $cookies, $cookieStore) {
	    console.log('hi from monitor1Controller');
    });

})();
