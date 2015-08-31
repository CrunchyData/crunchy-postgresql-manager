var ProxyDetailController = function($scope, $state, $cookieStore, $stateParams, utils, proxyFactory) {
    if (!$cookieStore.get('cpm_token')) {
        console.log('cpm_token not defined in projects');
        $state.go('login', {
            userId: 'hi'
        });
    }
	console.log("in proxy detail controller with containerId=" + $stateParams.containerId);
        proxyFactory.getbycontainerid($stateParams.containerId)
            .success(function(data) {
                console.log('successful getbycontainerid with data=' + data);
		$scope.proxy = data;
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });

};

var ProxyAddController = function($scope, $stateParams, $state, serversFactory, proxyFactory, utils, usSpinnerService) {

    var newcontainer = {};

    console.log('in ProxyAddController with projectId = ' + $stateParams.projectId);
    serversFactory.all()
        .success(function(data) {
            console.log('got servers' + data.length);
            $scope.servers = data;
            newcontainer.ID = 0;
            newcontainer.Name = 'newcontainer';
            newcontainer.Image = 'cpm-node';
            newcontainer.ServerID = $scope.servers[0].ID;
            $scope.selectedServer = $scope.servers[0];
            $scope.dockerprofile = 'SM';
            $scope.standalone = false;
            $scope.container = newcontainer;
        })
        .error(function(error) {
            $scope.alerts = [{
                type: 'danger',
                msg: error.Error
            }];
        });

    $scope.add = function() {
        usSpinnerService.spin('spinner-1');
        $scope.container.ServerID = $scope.selectedServer.ID;
        $scope.container.ProjectID = $stateParams.projectId;

        console.log('in add database with projectID = ' + $stateParams.projectId);
        $scope.container.ID = 0; //0 means to do an insert
        console.log('standalone is ' + $scope.standalone);

        proxyFactory.add($scope.proxy, $scope.container, $scope.standalone, $scope.dockerprofile)
            .success(function(data) {
                console.log('successful add with data=' + data);
                usSpinnerService.stop('spinner-1');
                $state.go('projects.proxy', $stateParams, {
                    reload: true,
                    inherit: false
                });
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                usSpinnerService.stop('spinner-1');
            });
    };
};
var ProxyDeleteController = function($scope, $stateParams, $state, proxyFactory, utils, usSpinnerService) {
    var container = $scope.container;

    $scope.delete = function() {
        usSpinnerService.spin('spinner-1');
        proxyFactory.delete($stateParams.containerId)
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
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
                usSpinnerService.stop('spinner-1');
            });
    };

};
