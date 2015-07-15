angular.module('uiRouterSample.rules.service', ['ngCookies'])

.factory('rulesFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var rulesFactory = {};

    rulesFactory.all = function() {
        var url = $cookieStore.get('AdminURL') + '/rules/getall/' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    rulesFactory.get = function(ruleid) {

        var url = $cookieStore.get('AdminURL') + '/rules/get/' + ruleid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    rulesFactory.delete = function(ruleid) {

        var url = $cookieStore.get('AdminURL') + '/rules/delete/' + ruleid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    rulesFactory.add = function(rule) {

        var url = $cookieStore.get('AdminURL') + '/rules/insert';
        console.log(url);

	console.log('token ' + $cookieStore.get('cpm_token'));

        return $http.post(url, {
            'Name': rule.Name,
            'Type': rule.Type,
            'Database': rule.Database,
            'Description': rule.Description,
            'User': rule.User,
            'Address': rule.Address,
            'Method': rule.Method,
            'Token': $cookieStore.get('cpm_token')

        });
    };

    rulesFactory.update = function(rule) {
    	rule.Token = $cookieStore.get('cpm_token');
        console.log('rule token ' + rule.Token);
        var url = $cookieStore.get('AdminURL') + '/rules/update';
        console.log(url);

        return $http.post(url, rule);
    };

    return rulesFactory;
}]);
