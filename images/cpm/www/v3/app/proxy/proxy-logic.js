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


var ProxyScheduleController = function($scope, $stateParams, $state, tasksFactory, serversFactory, utils) {

    $scope.edit = function() {
        $state.go('.edit', $stateParams);
    };
    console.log('in schedule controller [' + $stateParams.scheduleID + ']');
    serversFactory.all()
        .success(function(data) {
            console.log('successful get servers with data=' + data);
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
                console.log('successful get schedule with data=' + data);
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
                console.log('successful get schedules with data=' + data);
                console.log(JSON.stringify(data));
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
    console.log('proxy schedule add controller XXXX');
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
    console.log('jeff');
    //console.log(' servers here is ' + JSON.stringify($scope.servers));
    console.log(' other servers here is ' + JSON.stringify(servers));
    $scope.myServer = servers.data[0].ID;
    console.log('setting myServer to ' + $scope.myServer);
    $scope.servers = servers.data;
    console.log('top of scheed add with schedule ID=' + $stateParams.scheduleID);
    $scope.updateCurrentSchedule = function() {

        if ($scope.thething.checked == true) {
            $scope.schedule.Enabled = 'YES';
        } else {
            $scope.schedule.Enabled = 'NO';
        }

    };

    $scope.create = function() {
        console.log('create now called');
        console.log('with myServer ' + $scope.myServer);

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
        console.log($stateParams.scheduleID + 'stateparams');
        console.log($scope.schedule.ID + 'sched id');
        $scope.schedule.ServerID = $scope.myServer;
        console.log('jeff value of profile is...' + $scope.currentProfileName.name);
        $scope.schedule.ProfileName = $scope.currentProfileName.name;
        tasksFactory.addschedule($scope.schedule, $scope.proxy.ContainerName)
            .success(function(data) {
                console.log('successful add schedule with data=' + data);
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
  console.log('in delete schedule with scheduleID = ' + $stateParams.scheduleID);

    $scope.delete = function() {
        console.log('in schedule delete' + $stateParams.scheduleID);
        tasksFactory.deleteschedule($stateParams.scheduleID)
            .success(function(data) {
                console.log('successful deleteschedule with data=' + data);
                console.log(JSON.stringify(data));
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

var ProxyUsersDeleteController = function($scope, $stateParams, $state, containersFactory, utils) {
    console.log('before doing delete user id=' + $stateParams.containerId + ' name=' + $stateParams.itemId);
    $scope.rolname = $stateParams.itemId;

    $scope.delete = function() {
        console.log('doing delete user id=' + $stateParams.containerId + ' name=' + $scope.rolname);
        containersFactory.deleteuser($stateParams.containerId, $scope.rolname)
            .success(function(data) {
                console.log('successful deleteuser with data=' + data);
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
    console.log('doing add user');
    $scope.create = function() {
        $scope.user.ContainerID = $stateParams.containerId;
        console.log($scope.user);
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
                console.log('successful adduser with data=' + data);
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
            console.log('successful get user with Rolname=' + data.Rolname);

            $scope.user.Rolinherit = (data.Rolinherit === "true");
            $scope.user.Rolsuper = (data.Rolsuper === "true");
            console.log('fetched Rolsuper 1=' + data.Rolsuper);
            console.log('fetched Rolsuper=' + $scope.user.Rolsuper);
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
        console.log('save called');
        $scope.user.ContainerID = $stateParams.containerId;
        console.log('saving Rolsuper is ' + $scope.user.Rolsuper);
        containersFactory.updateuser($scope.user)
            .success(function(data) {

                console.log('successful updateuser with data=' + data);
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
