var ContainerDetailController = function($scope, $state, $cookieStore, $stateParams, utils) {
    if (!$cookieStore.get('cpm_token')) {
        console.log('cpm_token not defined in projects');
        $state.go('login', {
            userId: 'hi'
        });
    }

    //$state.go('projects.container.details', $stateParams);

};

var ContainerTaskSchedulesController = function($scope, $stateParams, $state, containersFactory, utils) {

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
};

var ContainerUsersController = function($scope, $stateParams, $state, containersFactory, utils) {
    containersFactory.getallusers($stateParams.containerId)
        .success(function(data) {
            console.log('successful get with data=' + data);
            $scope.users = data;
        })
        .error(function(error) {
            $scope.alerts = [{
                type: 'danger',
                msg: error.Error
            }];
            console.log('here is an error ' + error.Error);
        });

    $scope.edit = function() {
        $state.go('.edit', $stateParams);
    };
};


var ContainerUsersDeleteController = function($scope, $stateParams, $state, containersFactory, utils) {
    console.log('before doing delete user id=' + $stateParams.containerId + ' name=' + $stateParams.itemId);
    $scope.rolname = $stateParams.itemId;

    $scope.delete = function() {
        console.log('doing delete user id=' + $stateParams.containerId + ' name=' + $scope.rolname);
        containersFactory.deleteuser($stateParams.containerId, $scope.rolname)
            .success(function(data) {
                console.log('successful deleteuser with data=' + data);
                $state.go('projects.container.users', $stateParams, {
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
            });

    };
};

var ContainerStartController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var container = $scope.container;
    console.log('here in start top');

    $scope.start = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.start($stateParams.containerId)
            .success(function(data) {
                console.log('successful start with data=' + data);
                $state.go('projects.container.details', $stateParams);
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

var ContainerStopController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var container = $scope.container;

    $scope.stop = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.stop($stateParams.containerId)
            .success(function(data) {
                console.log('successful stop with data=' + data);
                $state.go('projects.container.details', $stateParams);
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

var ContainerMonitorpgstatdatabaseController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.refresh = function() {
        containersFactory.pgstatdatabase($stateParams.containerId)
            .success(function(data) {
                console.log('successful get with data=' + data);
                $scope.statdbresults = data;
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
};


var ContainerMonitorpgstatreplicationController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.refresh = function() {
        containersFactory.pgstatreplication($stateParams.containerId)
            .success(function(data) {
                console.log('successful get with data=' + data);
                $scope.statreplresults = data;
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

};

var ContainerMonitorbgwriterController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.refresh = function() {
        containersFactory.bgwriter($stateParams.containerId)
            .success(function(data) {
                console.log('successful get with data=' + data);
                $scope.bgwriterresults = data;
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
};


var ContainerMonitorbadgerController = function($sce, $scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    $scope.badgerreportlink = $sce.trustAsResourceUrl(
        'http://' + $scope.container.Name + ':10001/static/badger.html');
    $scope.refresh = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.badger($stateParams.containerId)
            .success(function(data) {
                console.log('successful get with data=' + data);
                usSpinnerService.stop('spinner-1');
                $state.go('projects.container.monitor.badger', $stateParams, {
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
        usSpinnerService.stop('spinner-1');
    };

    //$scope.refresh();
};


var ContainerMonitorpgsettingsController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.refresh = function() {
        containersFactory.pgsettings($stateParams.containerId)
            .success(function(data) {
                console.log('successful get with data=' + data);
                $scope.settingsresults = data;
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
};

var ContainerMonitorpgstatstatementsController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.refresh = function() {
        containersFactory.pgstatstatements($stateParams.containerId)
            .success(function(data) {
                $scope.statementresults = data;
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
};



var ContainerMonitorpgcontroldataController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.refresh = function() {
        containersFactory.pgcontroldata($stateParams.containerId)
            .success(function(data) {
                console.log('successful get with data=' + data);
                $scope.controldataresults = data;
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
};


var ContainerMonitorloadtestController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    $scope.slidervaluehigh = "10000";
    $scope.slidervaluelow = "1000";
    $scope.slidervalue = "1000";

    $scope.refresh = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.loadtest($stateParams.containerId, $scope.slidervalue)
            .success(function(data) {
                console.log('successful get with data=' + data);
                $scope.loadtestresults = data;
                usSpinnerService.stop('spinner-1');
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                usSpinnerService.stop('spinner-1');
                console.log('here is an error ' + error.Error);
            });
    };

};

var ContainerMonitordatabasesizeController = function($sce, $scope, $stateParams, $state, containersFactory, utils) {

    console.log('dbsize called with container Name ' + $scope.container.Name);
    console.log('dbsize called with container ID ' + $stateParams.containerId);
    $scope.dbsizegraphlink = $sce.trustAsResourceUrl('http://cpm-promdash:3000/embed/dbsizedashboard#!?var.container=' + $scope.container.Name);


};

var ContainerDeleteController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var container = $scope.container;

    $scope.delete = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.delete($stateParams.containerId)
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

var ContainerFailoverController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var container = $scope.container;

    $scope.failover = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.failover($stateParams.containerId)
            .success(function(data) {
                console.log('successful failover with data=' + data);
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


var ContainerAccessRulesController = function($cookieStore, $scope, $stateParams, $state, utils, containersFactory, usSpinnerService) {
		console.log('in access rules controller');
		console.log('containerid=' + $stateParams.containerId);

        containersFactory.getaccessrules($stateParams.containerId)
            .success(function(data) {
                console.log('successful getaccessrules');
                $scope.cars = data;
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });


    $scope.save = function() {
    	usSpinnerService.spin('spinner-1');
        console.log('save called cars=' + $scope.cars);
            angular.forEach($scope.cars, function(car) {
	    	car.Token = $cookieStore.get('cpm_token');
		car.ContainerID = $stateParams.containerId;
            });
        containersFactory.updateaccessrules($scope.cars)
            .success(function(data) {
                console.log('successful updateaccessrules');
                $scope.alerts = [{
                    type: 'success',
                    msg: 'successfully saved access rules'
                }];
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

var ContainerUsersEditController = function($scope, $stateParams, $state, containersFactory, utils) {
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
                $state.go('projects.container.users', $stateParams, {
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

var ContainerUsersDeleteController = function($scope, $stateParams, $state, containersFactory, utils) {
    console.log('before doing delete user id=' + $stateParams.containerId + ' name=' + $stateParams.itemId);
    $scope.rolname = $stateParams.itemId;

    $scope.delete = function() {
        console.log('doing delete user id=' + $stateParams.containerId + ' name=' + $scope.rolname);
        containersFactory.deleteuser($stateParams.containerId, $scope.rolname)
            .success(function(data) {
                console.log('successful deleteuser with data=' + data);
                $state.go('projects.container.users', $stateParams, {
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
            });

    };
};

var ContainerUsersAddController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.user = {};
    $scope.user.Password = '';
    $scope.user.Password2 = '';
    $scope.user.Rolsuper = false;
    $scope.user.Rolinherit = false;
    $scope.user.Rolcreaterole = false;
    $scope.user.Rolcreatedb = false;
    $scope.user.Rolcatupdate = false;
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
                $state.go('projects.container.users', $stateParams, {
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

var ContainerScheduleExecuteController = function($scope, $stateParams, $state, tasksFactory, utils, usSpinnerService) {

    console.log('in schedule execute' + $stateParams.scheduleID);
    usSpinnerService.spin('spinner-1');

    tasksFactory.execute($scope.schedule)
        .success(function(data) {
            console.log('successful post execute with data=' + data);
            $scope.alerts = [{
                type: 'success',
                msg: 'success'
            }];
            console.log(JSON.stringify(data));
    		usSpinnerService.stop('spinner-1');
        })
        .error(function(error) {
            $scope.alerts = [{
                type: 'danger',
                msg: error.Error
            }];
            usSpinnerService.stop('spinner-1');
            console.log('here is an error ' + error.Error);
        });
};

var ContainerScheduleHistoryController = function($scope, $stateParams, $state, tasksFactory, utils) {

    $scope.edit = function() {
        $state.go('.edit', $stateParams);
    };
    console.log('in schedule history' + $stateParams.scheduleID);
    tasksFactory.getallstatus($stateParams.scheduleID)
        .success(function(data) {
            console.log('successful get history with data=' + data);
            console.log(JSON.stringify(data));
            $scope.stats = data;
        })
        .error(function(error) {
            $scope.alerts = [{
                type: 'danger',
                msg: error.Error
            }];
            console.log('here is an error ' + error.Error);
        });
};

var ContainerScheduleController = function($scope, $stateParams, $state, tasksFactory, serversFactory, utils) {

    $scope.edit = function() {
        $state.go('.edit', $stateParams);
    };
    console.log('in schedule controller [' + $stateParams.scheduleID + ']');
    serversFactory.all()
        .success(function(data) {
            console.log('successful get servers with data=' + data);
            $scope.servers = data;
            $state.go('projects.container.schedule.edit', $stateParams, {
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

var ContainerScheduleAddController = function($scope, $filter, $stateParams, $state, tasksFactory, serversFactory, utils, servers) {
    console.log('XXXXX in schedule add controller XXXX');
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
        tasksFactory.addschedule($scope.schedule, $scope.container.Name)
            .success(function(data) {
                console.log('successful add schedule with data=' + data);
                $state.go('projects.container.taskschedules', $stateParams, {
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

var ContainerScheduleEditController = function($scope, $filter, $stateParams, $state, tasksFactory, serversFactory, utils, servers) {
    $scope.profiles = [{
        name: 'pg_basebackup'
    }, {
        name: 'pg_other'
    }];
    $scope.currentProfileName = $scope.profiles[0];

    $scope.myServer = [];
    $scope.thething = [{
        'name': 'thething',
        'checked': false
    }];

    console.log('top of scheed edit with schedule ID=' + $stateParams.scheduleID);
    console.log('top of scheed edit with servers len =' + $scope.servers);
    console.log('top of scheed edit with servers len23 =' + servers);

    tasksFactory.getschedule($stateParams.scheduleID)
        .success(function(data) {
            console.log('got the schedule');
            $scope.schedule = data;
            if ($scope.schedule.Enabled == 'YES') {
                $scope.thething.checked = true;
            } else {
                $scope.thething.checked = false;
            }
            $scope.setServer();
        })
        .error(function(error) {
            $scope.alerts = [{
                type: 'danger',
                msg: error.Error
            }];
            console.log('here is an error ' + error.Error);
        });

    $scope.updateCurrentSchedule = function() {

        if ($scope.thething.checked == true) {
            $scope.schedule.Enabled = 'YES';
        } else {
            $scope.schedule.Enabled = 'NO';
        }

    };

    $scope.executenow = function() {
        console.log('execute now called');

        tasksFactory.execute($scope.schedule)
            .success(function(data) {
                console.log('successful post execute with data=' + data);
                console.log(JSON.stringify(data));
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });
    };
    $scope.save = function() {
        console.log('save now called');

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
        console.log($scope.myServer.ID + 'myServer id');

        tasksFactory.updateschedule($scope.schedule)
            .success(function(data) {
                console.log('successful post schedule with data=' + data);
                console.log(JSON.stringify(data));
                $scope.alerts = [{
                    type: 'success',
                    msg: 'success'
                }];
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });


    };
    $scope.setServer = function() {
        if ($scope.servers.length > 0) {
            angular.forEach($scope.servers, function(item) {
                if (item.ID == $scope.schedule.ServerID) {
                    $scope.myServer = item;
                }
            });
        }
    };


};



var ContainerScheduleDeleteController = function($scope, $stateParams, $state, tasksFactory, utils) {

    console.log('in delete schedule with scheduleID = ' + $stateParams.scheduleID);

    $scope.delete = function() {
        console.log('in schedule delete' + $stateParams.scheduleID);
        tasksFactory.deleteschedule($stateParams.scheduleID)
            .success(function(data) {
                console.log('successful deleteschedule with data=' + data);
                console.log(JSON.stringify(data));
                $scope.stats = data;
                $state.go('projects.container.taskschedules', $stateParams, {
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

var ContainerAddController = function($scope, $stateParams, $state, serversFactory, containersFactory, utils, usSpinnerService) {

    var newcontainer = {};

    console.log('in ContainerAddController with projectId = ' + $stateParams.projectId);
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

        containersFactory.add($scope.container, $scope.standalone, $scope.dockerprofile)
            .success(function(data) {
                console.log('successful add with data=' + data);
                usSpinnerService.stop('spinner-1');
                $state.go('projects.container', $stateParams, {
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
