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
        var token = $cookieStore.get('cpm_token');

        $scope.containerid = window.containerid;
        $http.get($cookieStore.get('AdminURL') + '/node/' + window.containerid + '.' + token).success(function(data, status, headers, config) {
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
            var token = $cookieStore.get('cpm_token');

            var thing4 = $cookieStore.get('AdminURL') + '/monitor/container-loadtest/' + $scope.containerid + '.loadtest.' + $scope.slidervalue + "." + token;
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

            var token = $cookieStore.get('cpm_token');

            var thing2 = $cookieStore.get('AdminURL') + '/monitor/container-getinfo/' + $scope.containerid + '.statreplication.' + token;
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
            var token = $cookieStore.get('cpm_token');
            console.log('working on bgwriter');
            var thing3 = $cookieStore.get('AdminURL') + '/monitor/container-getinfo/' + $scope.containerid + '.bgwriter.' + token;
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
            var token = $cookieStore.get('cpm_token');
            var thing = $cookieStore.get('AdminURL') + '/monitor/container-getinfo/' + $scope.containerid + '.statdatabase.' + token;
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

        var loaded = false;
        var pg2seriesData = [];
        var pg2graph;
        var pg2axes;
        var pg2yAxis;
	$scope.refreshTime8h = '8h';
	$scope.refreshTime24h = '24h';
	$scope.refreshTime48h = '48h';
	$scope.refreshTime1w = '1w';

	var pg2graph = new Rickshaw.Graph( {
                        width: 800,
                        height: 300,
                        element: document.getElementById('pg2chart'),
                        renderer: 'line',
                        series: pg2seriesData,
	} );



       	var pg2hoverDetail = new Rickshaw.Graph.HoverDetail( { graph: pg2graph } );

	pg2axes = new Rickshaw.Graph.Axis.Time( { graph: pg2graph } );

        pg2yAxis = new Rickshaw.Graph.Axis.Y( {
                                graph: pg2graph,
                                tickFormat: Rickshaw.Fixtures.Number.formatKMBT,
                                ticksTreatment: 'glow'
	} );
	pg2graph.render();
       	pg2axes.render();
       	pg2yAxis.render();

        $scope.pg2handleRefresh = function(interval) {
            var token = $cookieStore.get('cpm_token');
	    var query = $cookies.AdminURL + '/mon/container/pg2/' + $scope.container.Name + '.' + interval + '.' + token;
		//console.log(query);
		$http.get(query).success(function(data, status, headers, config) {
			//console.log('pg2 query results 1 ' + JSON.stringify(data[0]));
			pg2loadSeries(data);
			if (loaded == false) {
				var pg2legend = new Rickshaw.Graph.Legend( {
					graph: pg2graph,
					element: document.getElementById('pg2legend')
				});
				loaded = true;
			}
            }).error(function(data, status, headers, config) {
                alert('error happended');
            });

        };

        $scope.pg2handleRefresh($scope.refreshTime8h);

        function pg2loadSeries(data) {
		
		//remove all existing data
		len = pg2seriesData.length;
		for (i=0; i<len; i++) {
			pg2seriesData.shift();
		}
	
		var palette = new Rickshaw.Color.Palette( { scheme: 'colorwheel' } );

		//add new data
                angular.forEach(data, function(d) {
			//console.log('data=' + JSON.stringify(d.Data));
                        pg2seriesData.push( { color: palette.color(), data: d.Data, name: d.Name } );

                });
                //console.log('pg2seriesData  ' + JSON.stringify(pg2seriesData));

		//refresh graph
		pg2graph.update();
        }


    }); /* end of monitor1controller */

})();
