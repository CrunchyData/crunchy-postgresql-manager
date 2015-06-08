angular.module('uiRouterSample.projects.service', ['ngCookies'])

.factory('projectsFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var projectsFactory = {};

    projectsFactory.all = function() {
        var url = $cookieStore.get('AdminURL') + '/project/getall/' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };


    projectsFactory.get = function(projectid) {

        var url = $cookieStore.get('AdminURL') + '/project/get/' + projectid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    projectsFactory.delete = function(projectid) {

        var url = $cookieStore.get('AdminURL') + '/project/delete/' + projectid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    projectsFactory.containers = function(projectid) {

        var url = $cookieStore.get('AdminURL') + '/nodes/forserver/' + projectid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    projectsFactory.add = function(project) {
        var url = $cookieStore.get('AdminURL') + '/project/add';
        console.log(url);
        return $http.post(url, project);
    };

    projectsFactory.update = function(project) {
        project.Token = $cookieStore.get('cpm_token');
        var url = $cookieStore.get('AdminURL') + '/project/update';
        console.log(url);
        return $http.post(url, project);
    };

    return projectsFactory;
}]);
