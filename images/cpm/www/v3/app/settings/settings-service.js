angular.module('uiRouterSample.settings.service', ['ngCookies'])

.factory('settingsFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var settingsFactory = {};

    settingsFactory.all = function() {
        var url = $cookieStore.get('AdminURL') + '/settings/' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    settingsFactory.savesetting = function(setting) {
        var url = $cookieStore.get('AdminURL') + '/savesetting';
        console.log(url);

        return $http.post(url, {
            'Name': setting.Name,
            'Value': setting.Value,
            'Token': $cookieStore.get('cpm_token')
        });
    };

    return settingsFactory;


}]);
