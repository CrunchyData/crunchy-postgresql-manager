var MyController = function($scope, $state, $cookieStore, $stateParams, utils) {
    if (!$cookieStore.get('cpm_token')) {
        console.log('cpm_token not defined in projects');
        $state.go('login', {
            userId: 'hi'
        });
    }
/**
    if ($scope.clusters.data.length > 0) {
        angular.forEach($scope.clusters.data, function(item) {
            if (item.ID == $stateParams.clusterId) {
                $scope.cluster = item;
                console.log(JSON.stringify(item));
            }
        });
    }
    */

};

var ClusterDeleteController = function($scope, $stateParams, $state, clustersFactory, utils, usSpinnerService) {
    var cluster = $scope.cluster;

    $scope.delete = function() {
        usSpinnerService.spin('spinner-1');
        console.log('delete cluster called');
        clustersFactory.delete($stateParams.clusterId)
            .success(function(data) {
                console.log('successful delete with data=' + data);
                usSpinnerService.stop('spinner-1');
                $state.go('projects.list', $stateParams, {
                    reload: true,
                    inherit: false
                });
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.message
                }];
                console.log('here is an error ' + error.message);
                usSpinnerService.stop('spinner-1');
            });
    };

};

var ClusterScaleController = function($scope, $stateParams, $state, clustersFactory, utils, usSpinnerService) {
    var cluster = $scope.cluster;

    $scope.scale = function() {
        usSpinnerService.spin('spinner-1');
        console.log('scale cluster called');
        clustersFactory.scale($stateParams.clusterId)
            .success(function(data) {
                console.log('successful scale with data=' + data);
                usSpinnerService.stop('spinner-1');
                $state.go('projects.list', $stateParams, {
                    reload: true,
                    inherit: false
                });
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.message
                }];
                console.log('here is an error ' + error.message);
                usSpinnerService.stop('spinner-1');
            });
    };

};

var ClusterAutoClusterController = function($scope, $stateParams, $state, clustersFactory, utils, usSpinnerService) {

    $scope.ClusterProfile = 'SM';
    $scope.cluster = [];
    $scope.cluster.Name = 'cluster1';
    $scope.cluster.ClusterType = 'asynchronous';
    console.log('auto cluster with projectId=' + $stateParams.projectId);
    $scope.cluster.ProjectID = $stateParams.projectId;

    console.log('ClusterProfile=' + $scope.ClusterProfile);
    console.log('cluster.Name=' + $scope.cluster.Name);

    $scope.create = function() {
        usSpinnerService.spin('spinner-1');
        console.log('ClusterProfile=' + $scope.ClusterProfile);
        console.log('cluster.Name=' + $scope.cluster.Name);
        console.log('cluster.ClusterType=' + $scope.cluster.ClusterType);
        clustersFactory.autocluster($scope.cluster, $scope.ClusterProfile)
            .success(function(data) {
                console.log('successful autocreate with data=' + data);
                $scope.results = data;
                usSpinnerService.stop('spinner-1');
                $state.go('projects.detail.details', $stateParams, {
                    reload: true,
                    inherit: false
                });
                $scope.alerts = [{
                    type: 'success',
                    msg: 'success'
                }];
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.message
                }];
                console.log('here is an error ' + error.message);
                usSpinnerService.stop('spinner-1');
            });
    };

};
