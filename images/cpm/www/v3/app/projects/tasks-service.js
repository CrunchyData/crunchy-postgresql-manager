angular.module('uiRouterSample.tasks.service', ['ngCookies'])

.factory('tasksFactory', ['$http', '$cookieStore', 'utils', function($http, $cookieStore, $scope, utils) {

    var tasksFactory = {};

    tasksFactory.getallschedules = function(containerid) {

        var url = $cookieStore.get('AdminURL') + '/task/getschedules/' + containerid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    tasksFactory.getschedule = function(scheduleid) {

        var url = $cookieStore.get('AdminURL') + '/task/getschedule/' + scheduleid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    tasksFactory.getallstatus = function(scheduleid) {

        var url = $cookieStore.get('AdminURL') + '/task/getallstatus/' + scheduleid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    tasksFactory.deleteschedule = function(scheduleid) {

        var url = $cookieStore.get('AdminURL') + '/task/deleteschedule/' + scheduleid + '.' + $cookieStore.get('cpm_token');
        console.log(url);

        return $http.get(url);
    };

    tasksFactory.updateschedule = function(schedule) {

	console.log('jeff set is ' + schedule.RestoreSet);

        var url = $cookieStore.get('AdminURL') + '/task/updateschedule';
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
            'Name': schedule.Name,
            'RestoreSet': schedule.RestoreSet,
            'RestoreRemotePath': schedule.RestoreRemotePath,
            'RestoreRemoteHost': schedule.RestoreRemoteHost,
            'RestoreRemoteUser': schedule.RestoreRemoteUser,
            'RestoreDbUser': schedule.RestoreDbUser,
            'RestoreDbPass': schedule.RestoreDbPass

        });
    };

    tasksFactory.execute = function(postMessage) {

        var url = $cookieStore.get('AdminURL') + '/task/executenow';
        console.log(url);

        return $http.post(url, {
            'Token': $cookieStore.get('cpm_token'),
            'ServerID': postMessage.ServerID,
            'ProjectID': postMessage.ProjectID,
            'ContainerName': postMessage.ContainerName,
            'ProfileName': postMessage.ProfileName,
            'DockerProfile': postMessage.DockerProfile,
            'ScheduleID': postMessage.ScheduleID,
            'StatusID': postMessage.StatusID
        });
    };

    tasksFactory.addschedule = function(schedule, containerName) {

        var url = $cookieStore.get('AdminURL') + '/task/addschedule';
        console.log(url);

        return $http.post(url, {
            'Token': $cookieStore.get('cpm_token'),
            'ServerID': schedule.ServerID,
            'ContainerName': containerName,
            'ProfileName': schedule.ProfileName,
            'Name': schedule.Name,
            'RestoreSet': schedule.RestoreSet,
            'RestoreRemotePath': schedule.RestoreRemotePath,
            'RestoreRemoteHost': schedule.RestoreRemoteHost,
            'RestoreRemoteUser': schedule.RestoreRemoteUser,
            'RestoreDbUser': schedule.RestoreDbUser,
            'RestoreDbPass': schedule.RestoreDbPass,
            'Serverip': schedule.Serverip
        });
    };

    return tasksFactory;
}]);
