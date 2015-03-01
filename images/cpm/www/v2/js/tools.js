// create the module and name it cpmApp
var cpmApp = angular.module('cpmApp.tools', ['ngRoute', 'ngCookies']);


cpmApp.controller('toolsController', function($scope, $cookies) {
    console.log('hi from tools controller');
    $scope.message = 'tools page.';
    if ($cookieStore.get('AdminURL')) {} else {
        alert('CPM AdminURL setting is NOT defined, please update on the Settings page before using CPM');
    }

});
