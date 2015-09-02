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
    var newproxy = {};
	newproxy.DatabaseName = "postgres";
	newproxy.DatabasePort = "5432";
	$scope.proxy = newproxy;

    console.log('in ProxyAddController with projectId = ' + $stateParams.projectId);
    serversFactory.all()
        .success(function(data) {
            console.log('got servers' + data.length);
            $scope.servers = data;
            newcontainer.ID = 0;
            newcontainer.Name = 'newproxy';
            newcontainer.Image = 'cpm-node-proxy';
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


var ProxyStartController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var proxy = $scope.proxy;
    console.log('here in start top');

    $scope.start = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.start($stateParams.containerId)
            .success(function(data) {
                console.log('successful start with data=' + data);
                $state.go('projects.proxy.details', $stateParams);
                usSpinnerService.stop('spinner-1');
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

var ProxyStopController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var proxy = $scope.proxy;

    $scope.stop = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.stop($stateParams.containerId)
            .success(function(data) {
                console.log('successful stop with data=' + data);
                $state.go('projects.proxy.details', $stateParams);
                usSpinnerService.stop('spinner-1');
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


var ProxyDatabasesizeController = function($sce, 
		$scope, $stateParams, $state, containersFactory, utils) {

		$scope.container = {}
		$scope.container.Name = $scope.proxy.ContainerName;

    		console.log('proxy dbsize called with container Name ' + $scope.proxy.ContainerName);
    		$scope.proxysizegraphlink = $sce.trustAsResourceUrl('http://cpm-promdash:3000/embed/dbsizedashboard#!?var.container=' + $scope.proxy.ContainerName);

};


var GotoproxyController = function($scope, $stateParams, $state, containersFactory, utils) {
	console.log('in GotoproxyController');
  	$state.go('projects.proxy.details', {
		containerId: $stateParams.containerId,
               	containerName: $stateParams.containerName,
               	projectId: $stateParams.projectId
	});
};



var ProxyDeleteController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var proxy = $scope.proxy;
	console.log('in delete ctlr with proxy=' + JSON.stringify(proxy));

    $scope.delete = function() {
        usSpinnerService.spin('spinner-1');
	console.log("deleting proxy with containerId=" + proxy.ContainerID);
        containersFactory.delete(proxy.ContainerID)
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
