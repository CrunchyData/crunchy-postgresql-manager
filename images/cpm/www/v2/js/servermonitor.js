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
        var token = $cookieStore.get('cpm_token');

	$scope.currentUser = [];
	$scope.currentUser = $cookieStore.get('cpmuser');
	console.log('currentUser is ' + $scope.currentUser);

        $http.get($cookieStore.get('AdminURL') + '/server/' + window.serverid + '.' + token).success(function(data, status, headers, config) {
            $scope.server = data;
        }).error(function(data, status, headers, config) {
            alert('error in get server');
        });

        console.log('working on iostat');
        $http.get($cookieStore.get('AdminURL') + '/monitor/server-getinfo/' + serverid + ".cpmiostat." + token).
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
            var token = $cookieStore.get('cpm_token');
            console.log('working on iostat');
            $http.get($cookieStore.get('AdminURL') + '/monitor/server-getinfo/' + serverid + ".cpmiostat." + token).
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
            var token = $cookieStore.get('cpm_token');
            console.log('working on df');
            console.log('calling getinfo service');
            $http.get($cookieStore.get('AdminURL') + '/monitor/server-getinfo/' + serverid + ".cpmdf." + token).
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

	    console.log('graphing server ' + $scope.server.Name); 
	var seriesData2;
	var graph, memgraph;
	var axes, memaxes;
	var yAxis, memyAxis;
	var graphCreated = false;
	var memgraphCreated = false;
	$scope.refreshTime8h = '8h';
	$scope.refreshTime24h = '24h';
	$scope.refreshTime48h = '48h';
	$scope.refreshTime1w = '1w';

        $scope.handleRefresh = function(interval) {
            var token = $cookieStore.get('cpm_token');
            console.log('graphing cpu interval= ' + interval);
	    var query = 'http://cpm-mon.crunchy.lab:8086/db/cpm/series?u=root&p=root&q=';
	    var query2 = 'select * from cpu where server = \'' + $scope.server.Name + '\' and time > now() - ' + interval + ' order asc limit 1000';
	    var es = escape(query2);

            $http.get(query + es).
            success(function(data, status, headers, config) {
		    loadSeries(data[0].points);
		    render();
                //console.log('flux query results ' + data[0].points);
                //console.log('first point t=' + data[0].points[0][0] + ' v=' + data[0].points[0][2]);
            }).error(function(data, status, headers, config) {
                alert('error happended');
            });

        };

	$scope.handleRefresh($scope.refreshTime8h);

        $scope.memhandleRefresh = function(interval) {
            var token = $cookieStore.get('cpm_token');
            console.log('graphing mem ');
	    var query = 'http://cpm-mon.crunchy.lab:8086/db/cpm/series?u=root&p=root&q=';
	    var query2 = 'select * from mem where server = \'' + $scope.server.Name + '\' and time > now() - ' + interval + ' order asc limit 1000';
	    //console.log('flux query ' + query + query2);
	    var es = escape(query2);
	    console.log('encoded query is ' + query + es);
            $http.get(query + es).
            success(function(data, status, headers, config) {
		    memloadSeries(data[0].points);
		    memrender();
                //console.log('mem flux query results ' + data[0].points);
                //console.log('mem first point t=' + data[0].points[0][0] + " v=" + data[0].points[0][2]);
            }).error(function(data, status, headers, config) {
                alert('error happended');
            });

        };

	$scope.memhandleRefresh($scope.refreshTime8h);

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
			//console.log('loading point x=' + p[0] + ' y=' + p[2]);
			//xval = Math.floor(p[0]/1000);
			//xval = new Date(p[0]);
			xval = Math.round(p[0] / 1000 );
			//console.log('loading point x=' + xval + ' y=' + p[2]);
			seriesData2.push( { x: xval, y: p[2] } );
		});
	}
	function memloadSeries(points) {
		memseriesData2 = [];
		angular.forEach(points, function(p) {
			//xval = Math.floor(p[0]/1000);
			//xval = new Date(p[0]);
			//xval = p[0];
			xval = Math.round(p[0]/1000);
			//console.log('mem loading point x=' + xval + ' y=' + p[2]);
			memseriesData2.push( { x: xval,  y: p[2] } );
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
