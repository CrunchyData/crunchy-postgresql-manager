angular.module('uiRouterSample.containers.service', ['ngCookies'])

.factory('containersFactory', ['$rootScope', '$http', '$cookieStore', 'utils', function($rootScope, $http, $cookieStore, utils) {

    var containersFactory = {};

    containersFactory.all = function() {
    	console.log('in containers all with projectId=' + $rootScope.projectId);
        var url = $cookieStore.get('AdminURL') + '/projectnodes/' + $rootScope.projectId + '.' + $cookieStore.get('cpm_token');
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

    containersFactory.failover = function(id) {

        var url = $cookieStore.get('AdminURL') + '/admin/failover/' + id + '.' + $cookieStore.get('cpm_token');
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
            container.ProjectID + '.' +
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

    containersFactory.pgstatreplication = function(id) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/repl/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.bgwriter = function(id) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/bgwriter/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.badger = function(id) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/badger/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    containersFactory.pgsettings = function(id) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/settings/' + id + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    containersFactory.pgstatstatements = function(id) {

        var url = $cookieStore.get('AdminURL') + '/monitor/container/statements/' + id + '.' + $cookieStore.get('cpm_token');
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

    containersFactory.getuser = function(containerid, rolname) {

        var url = $cookieStore.get('AdminURL') + '/dbuser/get/' + containerid + '.' + rolname + '.' +  $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.deleteuser = function(containerid, rolname) {

        var url = $cookieStore.get('AdminURL') + '/dbuser/delete/' + containerid + '.' + rolname + '.' +  $cookieStore.get('cpm_token');
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
		'Rolname' : user.Rolname,
		'Passwd' : user.Password,
		'Superuser' : user.Superuser,
		'Createdb' : user.Createdb,
		'Createrole' : user.Createrole,
		'TestBool' : user.TestBool,
		'Login' : user.Login,
		'Token' : $cookieStore.get('cpm_token'),
	});
    };

    containersFactory.updateuser = function(user) {

        var url = $cookieStore.get('AdminURL') + '/dbuser/update';
        console.log(url);
        console.log('id is ' + user.ContainerID);
	user.Token = $cookieStore.get('cpm_token');

        return $http.post(url,  {
		'ID' : user.ContainerID,
		'Rolname' : user.Rolname,
		'Passwd' : user.Password,
		'Rolsuper' : user.Rolsuper,
		'Rolinherit' : user.Rolinherit,
		'Rolcreaterole' : user.Rolcreaterole,
		'Rolcreatedb' : user.Rolcreatedb,
		'Rollogin' : user.Rollogin,
		'Rolcatupdate' : user.Rolcatupdate,
		'Rolreplication' : user.Rolreplication,
		'Login' : user.Rollogin,
		'Token' : $cookieStore.get('cpm_token'),
	});
    };

    containersFactory.getaccessrules = function(containerid) {
    	console.log('in containers getaccessrules with containerid=' + containerid);
        var url = $cookieStore.get('AdminURL') + '/containerrules/getall/' + containerid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    containersFactory.updateaccessrules = function(cars) {

	angular.forEach(cars, function(car) {
		car.Token = $cookieStore.get('cpm_token');
	});

        var url = $cookieStore.get('AdminURL') + '/containerrules/update';
        console.log(url);

        return $http.post(url,  cars);
    };

    return containersFactory;
}]);
