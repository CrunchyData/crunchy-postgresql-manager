angular.module('uiRouterSample.containers.service', ['ngCookies'])

.factory('containersFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var containersFactory = {};

    containersFactory.all = function() {
        var url = $cookieStore.get('AdminURL') + '/nodes/' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    containersFactory.get = function(id) {

        var url = $cookieStore.get('AdminURL') + '/node/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.delete = function(id) {

        var url = $cookieStore.get('AdminURL') + '/deletenode/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.start = function(id) {

        var url = $cookieStore.get('AdminURL') + '/admin/start/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };
    containersFactory.stop = function(id) {

        var url = $cookieStore.get('AdminURL') + '/admin/stop/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.add = function(container, standalone, dockerprofile) {

        var url = $cookieStore.get('AdminURL') + '/provision/' +
            dockerprofile + '.' +
            container.Image + '.' +
            container.ServerID + '.' +
            container.Name + '.' +
            standalone + '.' +
            $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.pgstatdatabase = function(id) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/database/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.bgwriter = function(id) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/bgwriter/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    containersFactory.pgsettings = function(id) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/settings/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    containersFactory.pgcontroldata = function(id) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/controldata/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.loadtest = function(id, slidervalue) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/loadtest/' + id + '.' + slidervalue + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.schedules = function(id) {

        var url = $cookieStore.get('AdminURL') + '/backup/getschedules/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.getallusers = function(id) {

        var url = $cookieStore.get('AdminURL') + '/dbuser/getall/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.adduser = function(user) {

        var url = $cookieStore.get('AdminURL') + '/dbuser/add';
        console.log(url);
        console.log('id is ' + user.ContainerID);
	user.Token = $cookieStore.get('cpm_token');

        return $http.post(url,  {
		'ID' : user.ContainerID,
		'Usename' : user.Usename,
		'Passwd' : user.Password,
		'Superuser' : user.Superuser,
		'Createdb' : user.Createdb,
		'Createrole' : user.Createrole,
		'Login' : user.Login,
		'Token' : $cookieStore.get('cpm_token'),
	});
    };


    return containersFactory;
}]);
