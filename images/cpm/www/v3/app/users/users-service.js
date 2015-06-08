angular.module('uiRouterSample.users.service', ['ngCookies'])

.factory('usersFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var usersFactory = {};

    usersFactory.all = function() {
        var url = $cookieStore.get('AdminURL') + '/sec/getusers/' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    usersFactory.get = function(serverid) {

        var url = $cookieStore.get('AdminURL') + '/server/' + serverid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    usersFactory.delete = function(username) {

        var url = $cookieStore.get('AdminURL') + '/sec/deleteuser/' + username + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    usersFactory.add = function(user) {

        var url = $cookieStore.get('AdminURL') + '/sec/adduser';
        console.log(url);

        return $http.post(url, {
            'Name': user.Name,
            'Password': user.Password,
            'Token': $cookieStore.get('cpm_token')

        });
    };

    usersFactory.save = function(thisuser) {
        console.log('user token ' + thisuser.Token);
        var url = $cookieStore.get('AdminURL') + '/sec/updateuser';
        console.log(url);

        return $http.post(url, thisuser);
    };

    usersFactory.changepsw = function(username, password) {
        var url = $cookieStore.get('AdminURL') + '/sec/cp';
        console.log(url);

        return $http.post(url, {
            'Username': username,
            'Password': password,
            'Token': $cookieStore.get('cpm_token')
        });
    };

    return usersFactory;
}]);
