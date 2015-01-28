(function(){
	  var app = angular.module('cpm', ['ui.bootstrap', 'ngRoute' ]);

	app.config(['$routeProvider',
			  function($routeProvider) {
				      $routeProvider.
		when('/', {
			              templateUrl: 'templates/home.html',
		              controller: 'HomeController'
			            }).
		when('/servers', {
			              templateUrl: 'templates/servers.html',
		              controller: 'ServersController'
			            }).
	      when('/containers', {
		              templateUrl: 'templates/containers.html',
	              controller: 'ContainersController'
		            }).
	            otherwise({
			            redirectTo: 'templates/servers.html'
			          });
		      }]);


	app.controller('PanelController', function($scope) {
		$scope.mmessage = 'this is the panel controller';
		console.log('in panel controller');
	});
	app.controller('HomeController', function($scope) {
		$scope.mmessage = 'this is the home page';
		console.log('in home controller');
	});

	app.controller('ServersController', function($scope) {
		$scope.mmessage = 'this is the servers page';
		console.log('in serverscontroller');
	});
	app.controller('ContainersController', function($scope) {
		$scope.mmessage = 'this is the Containersservers page';
		console.log('in containerscontroller');
	});
	app.controller('ClustersController', function($scope) {
		$scope.mmessage = 'this is the Clustersservers page';
		console.log('in clusterscontroller');
	});
	app.controller('ToolsController', function($scope) {
		$scope.mmessage = 'this is the Tools page';
		console.log('in toolscontroller');
	});
	app.controller('SettingsController', function($scope) {
		$scope.mmessage = 'this is the settings page';
		console.log('in settingscontroller');
	});


})();

