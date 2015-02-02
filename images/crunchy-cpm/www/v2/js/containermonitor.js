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
	    console.log('hi from monitor1Controller parent is ' + $scope.container.Name);

        var pg1seriesData;
        var pg1graph, pg2graph;
        var pg1axes, pg2axes;
        var pg1yAxis, pg2yAxis;
        var pg1graphCreated = false;
        var pg2graphCreated = false;
	$scope.refreshTime8h = '8h';
	$scope.refreshTime24h = '24h';
	$scope.refreshTime48h = '48h';

        $scope.pg1handleRefresh = function(interval) {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                alert('login required');
                return;
            }
            console.log('graphing pg1 stats interval=' + interval);
            var query = 'http://cluster-mon.crunchy.lab:8086/db/cpm/series?u=root&p=root&q=';
            var query2 = 'select * from pg1 where database = \'' + $scope.container.Name + '\'  and time > now() - ' + interval + ' order asc limit 1000';
	    var es = escape(query2);

            $http.get(query + es).
            success(function(data, status, headers, config) {
                    pg1loadSeries(data[0].points);
                    pg1render();
                //console.log('pg1 flux query results ' + data[0].points);
                //console.log('pg1 first point t=' + data[0].points[0][0] + " v=" + data[0].points[0][2]);
            }).error(function(data, status, headers, config) {
                alert('error happended');
            });

        };

        $scope.pg1handleRefresh($scope.refreshTime8h);

        $scope.pg2handleRefresh = function(interval) {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                alert('login required');
                return;
            }
            console.log('graphing pg2 with interval ' + interval);
            var query = 'http://cluster-mon.crunchy.lab:8086/db/cpm/series?u=root&p=root&q=';
            var query2 = 'select * from pg2 where database = \'' + $scope.container.Name + '\' and time > now() - ' + interval + ' order asc limit 1000';
	    var es = escape(query2);

            $http.get(query + es).
            success(function(data, status, headers, config) {
                    pg2loadSeries(data[0].points);
                    pg2render();
                //console.log('pg2 query results ' + data[0].points);
                //console.log('pg2 first point t=' + data[0].points[0][0] + " v=" + data[0].points[0][2]);
            }).error(function(data, status, headers, config) {
                alert('error happended');
            });

        };

        $scope.pg2handleRefresh($scope.refreshTime8h);

        function pg1render() {

                if (!pg1graphCreated) {
                        pg1graph = new Rickshaw.Graph( {
                        element: document.getElementById("pg1chart"),
                        width: 800,
                        height: 300,
                        renderer: 'line',
                        series: [
                        {
                        color: "#c05020",
                        data: pg1seriesData,
                        name: 'PG1' } ]
                        } );

                        var pg1hoverDetail = new Rickshaw.Graph.HoverDetail( {
                        graph: pg1graph
                        } );

                        var ticksTreatment = 'glow';
                        pg1axes = new Rickshaw.Graph.Axis.Time( {
                                graph: pg1graph
                        } );
                        pg1yAxis = new Rickshaw.Graph.Axis.Y( {
                                graph: pg1graph,
                                tickFormat: Rickshaw.Fixtures.Number.formatKMBT,
                                ticksTreatment: ticksTreatment
                        } );

                        pg1graphCreated = true;
                }

                pg1graph.render();
                pg1axes.render();
                pg1yAxis.render();
        }

        function pg1loadSeries(points) {
                pg1seriesData = [];
                angular.forEach(points, function(p) {
                        //console.log('pg1 loading point x=' + p[0] + ' y=' + p[2]);
			xval = Math.round(p[0] / 1000 );

                        pg1seriesData.push( { x: xval, y: p[2] } );
                });
        }

        function pg2loadSeries(points) {
                pg2seriesData = [];
                angular.forEach(points, function(p) {
                        //console.log('pg2 loading point x=' + p[0] + ' y=' + p[2]);
			xval = Math.round(p[0] / 1000 );
                        pg2seriesData.push( { x: xval, y: p[2] } );
                });
        }

        function pg2render() {

                if (!pg2graphCreated) {
                        pg2graph = new Rickshaw.Graph( {
                        element: document.getElemenById("pg2chart"),
                        width: 800,
                        height: 300,
                        renderer: 'line',
                        series: [
                        { color: "#c05020", data: pg2seriesData,
                        name: 'PG 2' } ]
                        } );

                        var pg2hoverDetail = new Rickshaw.Graph.HoverDetail( {
                                graph: pg2graph
                        } );

                        pg2axes = new Rickshaw.Graph.Axis.Time( {
                                graph: pg2graph
                        } );
                        pg2yAxis = new Rickshaw.Graph.Axis.Y( {
                                graph: pg2graph,
                                tickFormat: Rickshaw.Fixtures.Number.formatKMBT,
                                ticksTreatment: 'glow'
                        } );

                        pg2graphCreated = true;
                }

                pg2graph.render();
                pg2axes.render();
                pg2yAxis.render();
        }

    }); /* end of monitor1controller */

})();
