angular.module('uiRouterSample.authn.service', ['ngCookies'])

.factory('authnFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var authnFactory = {};

    authnFactory.doLogin = function(user_id, password, adminurl) {
        var loginUrl = adminurl + '/sec/login/' + user_id + "." + password;
        console.log(loginUrl);

        return $http.get(loginUrl);
    };

    authnFactory.doOtherLogin = function(user_id, password, adminurl) {
        console.log(user_id + ' is user_id in the service factory');
        var obj = {
            alerts: null,
            content: null
        };

        var loginUrl = adminurl + '/sec/login/' + user_id + "." + password;
        console.log(loginUrl);

        $http.get(loginUrl).
        success(function(data, status, headers, config) {
            //$rootScope.cpm_user_id = user_id;
            $cookieStore.put('cpm_user_id', user_id);
            $cookieStore.put('AdminURL', adminurl);
            $cookieStore.put('cpm_token', data.Contents);
            //$location.path('/home');
            console.log(status + ' was status');
            console.log(data.Contents + ' was returned');
            obj.content = data.Contents;
            obj.alerts = [{
                type: 'success',
                msg: 'Logged out.'
            }];
            $scope.loginResults = obj;
        }).
        error(function(data, status, headers, config) {
            obj.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

        return obj;
    };

    authnFactory.doLogout = function() {
        var adminurl = $cookieStore.get('AdminURL');
        var token = $cookieStore.get('cpm_token');

        var logoutUrl = adminurl + '/sec/logout/' + token;
        console.log(logoutUrl);

        return $http.get(logoutUrl);
    };

    authnFactory.doOtherLogout = function() {
        var obj = {
            alerts: null,
            content: null
        };
        var adminurl = $cookieStore.get('AdminURL');
        var token = $cookieStore.get('cpm_token');

        var logoutUrl = adminurl + '/sec/logout/' + token;
        console.log(logoutUrl);

        $http.get(logoutUrl).
        success(function(data, status, headers, config) {
            $cookieStore.remove('cpm_token');
            obj.content = data.Contents;
            obj.alerts = [{
                type: 'success',
                msg: 'Logged out.'
            }];
        }).
        error(function(data, status, headers, config) {
            obj.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

        $scope.loginResults = obj;
        return obj;
    };


    return authnFactory;
}]);
