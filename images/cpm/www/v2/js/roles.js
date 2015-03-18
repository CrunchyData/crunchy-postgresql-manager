// create the module and name it cpmApp
	var roleApp = angular.module('RoleApp', ['ngRoute']);


	roleApp.controller('roleController', function($scope, $http) {
		$scope.hithere = 'fromrolecontroller';
		//console.log('hi from roleController');
	});

