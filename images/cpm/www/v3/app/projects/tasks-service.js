angular.module('uiRouterSample.tasks.service', ['ngCookies'])

.factory('tasksFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var tasksFactory = {};

    tasksFactory.getallschedules = function(containerid) {

        var url = $cookieStore.get('AdminURL') + '/backup/getschedules/' + containerid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    tasksFactory.getschedule = function(scheduleid) {

        var url = $cookieStore.get('AdminURL') + '/backup/getschedule/' + scheduleid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    tasksFactory.getallstatus = function(scheduleid) {

        var url = $cookieStore.get('AdminURL') + '/backup/getallstatus/' + scheduleid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    tasksFactory.deleteschedule = function(scheduleid) {

        var url = $cookieStore.get('AdminURL') + '/backup/deleteschedule/' + scheduleid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    tasksFactory.updateschedule = function(schedule) {

        var url = $cookieStore.get('AdminURL') + '/backup/updateschedule';
        console.log(url);

        return $http.post(url, {
            'Token': $cookieStore.get('cpm_token'),
            'ID': schedule.ID,
            'ServerID': schedule.ServerID,
            'Enabled': schedule.Enabled,
            'Minutes': schedule.Minutes,
            'Hours': schedule.Hours,
            'DayOfMonth': schedule.DayOfMonth,
            'Month': schedule.Month,
            'DayOfWeek': schedule.DayOfWeek,
            'Name': schedule.Name

        });
    };

    tasksFactory.execute = function(schedule) {

        var url = $cookieStore.get('AdminURL') + '/backup/now';
        console.log(url);

        return $http.post(url, {
            'Token': $cookieStore.get('cpm_token'),
            'ServerID': schedule.ServerID,
            'ProfileName': schedule.ProfileName,
            'ScheduleID': schedule.ID
        });
    };

    tasksFactory.addschedule = function(schedule, containerName) {

        var url = $cookieStore.get('AdminURL') + '/backup/addschedule';
        console.log(url);

        console.log('serverid ' + schedule.ServerID);
        console.log('containername ' + containerName);
        console.log('profile ' + schedule.ProfileName);
        console.log('schedulename ' + schedule.Name);
        return $http.post(url, {
            'Token': $cookieStore.get('cpm_token'),
            'ServerID': schedule.ServerID,
            'ContainerName': containerName,
            'ProfileName': schedule.ProfileName,
            'Name': schedule.Name
        });
    };

    return tasksFactory;
}]);
