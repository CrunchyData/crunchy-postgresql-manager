angular.module('uiRouterSample.roles.service', ['ngCookies'])

.factory('rolesFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var rolesFactory = {};

    rolesFactory.all = function() {
        var url = $cookieStore.get('AdminURL') + '/sec/getroles/' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    rolesFactory.get = function(serverid) {

        var url = $cookieStore.get('AdminURL') + '/server/' + serverid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    rolesFactory.delete = function(rolename) {

        var url = $cookieStore.get('AdminURL') + '/sec/deleterole/' + rolename + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    rolesFactory.save = function(role) {
        role.Token = $cookieStore.get('cpm_token');
        var url = $cookieStore.get('AdminURL') + '/sec/updaterole';
        console.log(url);

        return $http.post(url, role);
    };



    rolesFactory.add = function(role) {

        var url = $cookieStore.get('AdminURL') + '/sec/addrole';
        console.log(url);
        role.Token = $cookieStore.get('cpm_token');
        console.log('adding with role=' + role.Token);
        return $http.post(url, {
            'Name': role.Name,
            'Token': role.Token
        });
    };

    return rolesFactory;
}]);
