var ProxyDetailController = function($scope, $state, $cookieStore, $stateParams, utils, proxyFactory) {
    if (!$cookieStore.get('cpm_token')) {
        console.log('cpm_token not defined in projects');
        $state.go('login', {
            userId: 'hi'
        });
    }
        proxyFactory.getbycontainerid($stateParams.containerId)
            .success(function(data) {
		$scope.proxy = data;
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });
    $scope.save = function() {
        console.log('save logic should go here');
        console.log(JSON.stringify($scope.proxy));
        proxyFactory.update($scope.proxy)
            .success(function(data) {
                $scope.alerts = [{
                    type: 'success',
                    msg: 'proxy saved'
                }];
	    /**
                $state.go('projects.proxy', $stateParams, {
                    reload: true,
                    inherit: false
                });
		*/
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
            });
    };

};


var ProxyScheduleController = function($scope, $stateParams, $state, tasksFactory, serversFactory, utils) {

    $scope.edit = function() {
        $state.go('.edit', $stateParams);
    };
    serversFactory.all()
        .success(function(data) {
            $scope.servers = data;
            $state.go('projects.proxy.schedule.edit', $stateParams, {
                reload: false,
                inherit: true
            });
        })
        .error(function(error) {
            $scope.alerts = [{
                type: 'danger',
                msg: error.Error
            }];
            console.log('here is an error ' + error.Error);
        });

    if ($stateParams.scheduleID != '') {
        tasksFactory.getschedule($stateParams.scheduleID)
            .success(function(data) {
                $scope.schedule = data;
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });
    } else {
        console.log('skipping get schedule on new schedule ');
    }

};

var ProxyTaskSchedulesController = function($scope, $stateParams, $state, serversFactory, containersFactory ) {
    $scope.refresh = function() {
        containersFactory.schedules($stateParams.containerId)
            .success(function(data) {
                $scope.schedules = data;
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });
    };

    $scope.refresh();

}

var ProxyScheduleAddController = function($scope, $filter, $stateParams, $state, tasksFactory, serversFactory, utils, servers) {
    $scope.profiles = [{
        name: 'pg_basebackup'
    }, {
        name: 'pg_other'
    }];
    $scope.currentProfileName = $scope.profiles[0];

    $scope.thething = [{
        'name': 'thething',
        'checked': false
    }];
    //console.log(' servers here is ' + JSON.stringify($scope.servers));
    $scope.myServer = servers.data[0].ID;
    $scope.servers = servers.data;
    $scope.updateCurrentSchedule = function() {

        if ($scope.thething.checked == true) {
            $scope.schedule.Enabled = 'YES';
        } else {
            $scope.schedule.Enabled = 'NO';
        }

    };

    $scope.create = function() {

        if ($scope.schedule.Minutes == '') {
            $scope.schedule.Minutes = '*';
        }
        if ($scope.schedule.Hours == '') {
            $scope.schedule.Hours = '*';
        }
        if ($scope.schedule.DayOfMonth == '') {
            $scope.schedule.DayOfMonth = '*';
        }
        if ($scope.schedule.Month == '') {
            $scope.schedule.Month = '*';
        }
        if ($scope.schedule.DayOfWeek == '') {
            $scope.schedule.DayOfWeek = '*';
        }
        $scope.schedule.ServerID = $scope.myServer;
        $scope.schedule.ProfileName = $scope.currentProfileName.name;
        tasksFactory.addschedule($scope.schedule, $scope.proxy.ContainerName)
            .success(function(data) {
                $state.go('projects.proxy.taskschedules', $stateParams, {
                    reload: false,
                    inherit: false
                });
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });

    };
};

var ProxyScheduleDeleteController = function($scope, $stateParams, $state, serversFactory, proxyFactory, utils, usSpinnerService) {

    $scope.delete = function() {
        tasksFactory.deleteschedule($stateParams.scheduleID)
            .success(function(data) {
                $scope.stats = data;
                $state.go('projects.proxy.taskschedules', $stateParams, {
                    reload: false,
                    inherit: false
                });
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });
    };

}

var ProxyAddController = function($scope, $stateParams, $state, serversFactory, proxyFactory, utils, usSpinnerService) {

    var newcontainer = {};
    var newproxy = {};
	newproxy.DatabaseName = "postgres";
	newproxy.DatabasePort = "5432";
	$scope.proxy = newproxy;

    serversFactory.all()
        .success(function(data) {
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

        $scope.container.ID = 0; //0 means to do an insert

        proxyFactory.add($scope.proxy, $scope.container, $scope.standalone, $scope.dockerprofile)
            .success(function(data) {
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

    $scope.start = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.start($stateParams.containerId)
            .success(function(data) {
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

    		$scope.proxysizegraphlink = $sce.trustAsResourceUrl('http://cpm-promdash:3000/embed/dbsizedashboard#!?var.container=' + $scope.proxy.ContainerName);

};


var GotoproxyController = function($scope, $stateParams, $state, containersFactory, utils) {
  	$state.go('projects.proxy.details', {
		containerId: $stateParams.containerId,
               	containerName: $stateParams.containerName,
               	projectId: $stateParams.projectId
	});
};

var ProxyUsersDeleteController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.rolname = $stateParams.itemId;

    $scope.delete = function() {
        containersFactory.deleteuser($stateParams.containerId, $scope.rolname)
            .success(function(data) {
                $state.go('projects.detail.gotoproxy', $stateParams, {
                    reload: true,
                    inherit: true
                });

            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
                console.log('here is this alerts ' + $scope.alerts);
            });

    };
};

var ProxyUsersAddController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.user = {};
    $scope.user.Password = '';
    $scope.user.Password2 = '';
    $scope.user.Rolsuper = false;
    $scope.user.Rolinherit = false;
    $scope.user.Rolcreaterole = false;
    $scope.user.Rolcreatedb = false;
    $scope.user.Rollogin = false;
    $scope.user.Rolreplication = false;

    $scope.user.Rolname = '';
    $scope.create = function() {
        $scope.user.ContainerID = $stateParams.containerId;
        if ($scope.user.Password !=
            $scope.user.Password2) {
            $scope.alerts = [{
                type: 'danger',
                msg: 'passwords do not match'
            }];
            return;
        }

        containersFactory.adduser($scope.user)
            .success(function(data) {
                $state.go('projects.proxy', $stateParams, {
                    reload: false,
                    inherit: false
                });

            })
            .error(function(error) {
                console.log(JSON.stringify(error));
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });

    };
};

var ProxyUsersEditController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.user = {};
    $scope.user.ContainerID = $stateParams.containerId;

    containersFactory.getuser($stateParams.containerId, $stateParams.itemId)
        .success(function(data) {
            $scope.user = data;

            $scope.user.Rolinherit = (data.Rolinherit === "true");
            $scope.user.Rolsuper = (data.Rolsuper === "true");
            $scope.user.Rolcreatedb = (data.Rolcreatedb === "true");
            $scope.user.Rolcreaterole = (data.Rolcreaterole === "true");
            $scope.user.Rollogin = (data.Rollogin === "true");
            $scope.user.Rolreplication = (data.Rolreplication === "true");
        })
        .error(function(error) {
            $scope.alerts = [{
                type: 'danger',
                msg: error.Error
            }];
            console.log('here is an error ' + error.Error);
        });
    $scope.save = function() {
    	console.log('in save');
        $scope.user.ContainerID = $stateParams.containerId;
        containersFactory.updateuser($scope.user)
            .success(function(data) {

                $state.go('projects.proxy.users', $stateParams, {
                    reload: false,
                    inherit: false
                });

            })
            .error(function(error) {
                console.log(JSON.stringify(error));
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });

    };
};

var ProxyDeleteController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var proxy = $scope.proxy;

    $scope.delete = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.delete(proxy.ContainerID)
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
                    msg: error.Error
                }];
                usSpinnerService.stop('spinner-1');
            });
    };


};
