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

	var seriesData2;
	var graph, memgraph;
	var axes, memaxes;
	var yAxis, memyAxis;
	var graphCreated = false;
	var memgraphCreated = false;

        $scope.handleRefresh = function() {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                alert('login required');
                return;
            }
            console.log('graphing server stats');
	    var query = 'http://localhost:8086/db/cpm/series?u=root&p=root&q=select * from cpu where server = \'myserver\' order asc limit 100';
            $http.get(query).
            success(function(data, status, headers, config) {
		    loadSeries(data[0].points);
		    render();
                console.log('flux query results ' + data[0].points);
                console.log('first point t=' + data[0].points[0][0] + " v=" + data[0].points[0][2]);
            }).error(function(data, status, headers, config) {
                alert('error happended');
            });

        };

	$scope.handleRefresh();

        $scope.memhandleRefresh = function() {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                alert('login required');
                return;
            }
            console.log('graphing server stats');
	    var query = 'http://localhost:8086/db/cpm/series?u=root&p=root&q=select * from mem where server = \'myserver\' order asc limit 100';
            $http.get(query).
            success(function(data, status, headers, config) {
		    memloadSeries(data[0].points);
		    memrender();
                console.log('mem flux query results ' + data[0].points);
                console.log('mem first point t=' + data[0].points[0][0] + " v=" + data[0].points[0][2]);
            }).error(function(data, status, headers, config) {
                alert('error happended');
            });

        };

	$scope.memhandleRefresh();

	function render() {

		if (!graphCreated) {
			graph = new Rickshaw.Graph( {
                        element: document.getElementById("chart"),
                        width: 800,
                        height: 300,
                        renderer: 'line',
                        series: [
                        {
                        color: "#c05020",
                        //data: seriesData[0],
                        data: seriesData2,
                        name: 'CPU Load'
                        }
                        ]
			} );
			
			var hoverDetail = new Rickshaw.Graph.HoverDetail( {
        		graph: graph
			} );

			var ticksTreatment = 'glow';
			axes = new Rickshaw.Graph.Axis.Time( {
				graph: graph
			} );
			yAxis = new Rickshaw.Graph.Axis.Y( {
				graph: graph,
				tickFormat: Rickshaw.Fixtures.Number.formatKMBT,
				ticksTreatment: ticksTreatment
			} );

			graphCreated = true;
		}

		graph.render();
		axes.render();
		yAxis.render();
	}

	function loadSeries(points) {
		seriesData2 = [];
		angular.forEach(points, function(p) {
			console.log('loading point x=' + p[0] + ' y=' + p[2]);
			seriesData2.push( { x: p[0]/1000, y: p[2] } );
		});
	}
	function memloadSeries(points) {
		memseriesData2 = [];
		angular.forEach(points, function(p) {
			console.log('mem loading point x=' + p[0] + ' y=' + p[2]);
			memseriesData2.push( { x: p[0]/1000, y: p[2] } );
		});
	}
	function memrender() {

		if (!memgraphCreated) {
			memgraph = new Rickshaw.Graph( {
                        element: document.getElementById("memchart"),
                        width: 800,
                        height: 300,
                        renderer: 'line',
                        series: [
                        { color: "#c05020", data: memseriesData2,
                        name: 'Mem Usage' }
                        ]
			} );
		
			var memhoverDetail = new Rickshaw.Graph.HoverDetail( {
        			graph: memgraph
			} );
		
			memaxes = new Rickshaw.Graph.Axis.Time( {
				graph: memgraph
			} );
			memyAxis = new Rickshaw.Graph.Axis.Y( {
				graph: memgraph,
				tickFormat: Rickshaw.Fixtures.Number.formatKMBT,
				ticksTreatment: 'glow'
			} );

			memgraphCreated = true;
		}

		memgraph.render();
		memaxes.render();
		memyAxis.render();
	}

    });



})();
