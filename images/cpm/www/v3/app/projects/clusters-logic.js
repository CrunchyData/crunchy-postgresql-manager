var MyController = function($scope, $state, $cookieStore, $stateParams, utils) {
    if (!$cookieStore.get('cpm_token')) {
        console.log('cpm_token not defined in projects');
        $state.go('login', {
            userId: 'hi'
        });
    }

};


var ClusterStopController = function($scope, $stateParams, $state, clustersFactory, utils, usSpinnerService) {
    var cluster = $scope.cluster;

    $scope.stop = function() {
        usSpinnerService.spin('spinner-1');
        console.log('stop cluster called');
        clustersFactory.stop($stateParams.clusterId)
            .success(function(data) {
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

var ClusterStartController = function($scope, $stateParams, $state, clustersFactory, utils, usSpinnerService) {
    var cluster = $scope.cluster;

    $scope.start = function() {
        usSpinnerService.spin('spinner-1');
        console.log('start cluster called');
        clustersFactory.start($stateParams.clusterId)
            .success(function(data) {
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

var ClusterDeleteController = function($scope, $stateParams, $state, clustersFactory, utils, usSpinnerService) {
    var cluster = $scope.cluster;

    $scope.delete = function() {
        usSpinnerService.spin('spinner-1');
        //console.log('delete cluster called');
        clustersFactory.delete($stateParams.clusterId)
            .success(function(data) {
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
        clustersFactory.scale($stateParams.clusterId)
            .success(function(data) {
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
    $scope.cluster.ProjectID = $stateParams.projectId;

    $scope.create = function() {
        usSpinnerService.spin('spinner-1');
        clustersFactory.autocluster($scope.cluster, $scope.ClusterProfile)
            .success(function(data) {
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
