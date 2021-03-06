var ContainerDetailController = function($scope, $state, $cookieStore, $stateParams, utils) {
    if (!$cookieStore.get('cpm_token')) {
        console.log('cpm_token not defined in projects');
        $state.go('login', {
            userId: 'hi'
        });
    }

};

var GotocontainerController = function($scope, $state, $cookieStore, $stateParams, utils) {
    $state.go('projects.container.details', {
        containerId: $stateParams.containerId,
        projectId: $stateParams.projectId
    });


};

var ContainerTaskSchedulesController = function($scope, $stateParams, $state, containersFactory, utils) {

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


};

var ContainerUsersController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.refresh = function() {
        containersFactory.getallusers($stateParams.containerId)
            .success(function(data) {
                $scope.users = data;
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            });
    }

    $scope.refresh();

    $scope.edit = function() {
        $state.go('.edit', $stateParams);
    };
};


var ContainerUsersDeleteController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.rolname = $stateParams.itemId;

    $scope.delete = function() {
        containersFactory.deleteuser($stateParams.containerId, $scope.rolname)
            .success(function(data) {
                $state.go('projects.container', $stateParams, {
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

var ContainerStartController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var container = $scope.container;

    $scope.start = function() {
        usSpinnerService.spin('spinner-1');
        //console.log('jeff state is ' + JSON.stringify($stateParams))
        containersFactory.start($stateParams.containerId)
            .success(function(data) {
                //$state.go('projects.container.details', $stateParams, {
                //reload: true,
                //inherit: false
                //});
                $state.go('projects.container.details', {
                    containerId: $stateParams.containerId,
                    projectId: $stateParams.projectId
                });

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


var ContainerMonitorbadgerController = function($sce, $scope, $stateParams, $state, containersFactory, utils, usSpinnerService, spinnerService) {
    $scope.badgerreportlink = $sce.trustAsResourceUrl(
        'http://' + $stateParams.containerName + ':10001/static/badger.html?rand=' + Math.round(Math.random() * 10000000));
    $scope.refresh = function() {
        usSpinnerService.spin('spinner-1');
        spinnerService.show('badgerSpinner');
        containersFactory.badger($stateParams.containerId)
            .success(function(data) {
                usSpinnerService.stop('spinner-1');
                //$state.go('projects.container.monitor.badger', $stateParams, {
                $state.go('^', $stateParams, {
                    reload: false,
                    inherit: true
                });
                $scope.alerts = [{
                    type: 'success',
                    msg: 'successfully generated report'
                }];

            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                console.log('here is an error ' + error.Error);
            })
            .finally(function() {
                spinnerService.hide('badgerSpinner');
            });
        usSpinnerService.stop('spinner-1');
    };

    //$scope.refresh();
};


var ContainerMonitorpgsettingsController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.refresh = function() {
        containersFactory.pgsettings($stateParams.containerId)
            .success(function(data) {
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

    $scope.dbsizegraphlink = $sce.trustAsResourceUrl('http://cpm-promdash:3000/embed/dbsizedashboard#!?var.container=' + $scope.container.Name);


};

var ContainerDeleteController = function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    var container = $scope.container;

    $scope.delete = function() {
        usSpinnerService.spin('spinner-1');
        containersFactory.delete($stateParams.containerId)
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

    containersFactory.getaccessrules($stateParams.containerId)
        .success(function(data) {
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
        angular.forEach($scope.cars, function(car) {
            car.Token = $cookieStore.get('cpm_token');
            car.ContainerID = $stateParams.containerId;
        });
        containersFactory.updateaccessrules($scope.cars)
            .success(function(data) {
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
        $scope.user.ContainerID = $stateParams.containerId;
        containersFactory.updateuser($scope.user)
            .success(function(data) {
                $state.go('projects.container', $stateParams, {
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

var ContainerUsersDeleteController = function($scope, $stateParams, $state, containersFactory, utils) {
    $scope.rolname = $stateParams.itemId;

    $scope.delete = function() {
        containersFactory.deleteuser($stateParams.containerId, $scope.rolname)
            .success(function(data) {
                $state.go('projects.container', $stateParams, {
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

var ContainerUsersAddController = function($scope, $stateParams, $state, containersFactory, utils) {
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
                $state.go('projects.container.users', $stateParams, {
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

var ContainerScheduleExecuteController = function($scope, $stateParams, $state, tasksFactory, utils, usSpinnerService) {

    usSpinnerService.spin('spinner-1');

	var postMessage = {};
	postMessage.ServerID = $scope.schedule.ServerID;
	postMessage.ContainerName = $scope.schedule.ContainerName;
	postMessage.ProfileName = $scope.schedule.ProfileName;
	postMessage.ScheduleID = $scope.schedule.ID;
	postMessage.ProjectID = $stateParams.projectId;

    tasksFactory.execute(postMessage)
        .success(function(data) {
            $scope.alerts = [{
                type: 'success',
                msg: 'success'
            }];
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

var ContainerScheduleHistoryRestoreController = function($scope, $stateParams, $state, tasksFactory, utils, usSpinnerService) {

    console.log('in schedule history restore controller');
    console.log('params = ' + JSON.stringify($stateParams));
    console.log('scope container = ' + JSON.stringify($scope.container));
    console.log('scope container = ' + $scope.container["Name"]);

    var newcontainer = {};
    newcontainer.ID = 0;
    newcontainer.Name = 'newcontainer';
    newcontainer.Image = 'cpm-node';
    $scope.dockerprofile = 'SM';
    $scope.container["Name"] = $scope.container["Name"] + '-restored';

    $scope.restore = function() {
        console.log('restore called to restore backup with schedule=' + $scope.schedule.ID);
        console.log('dockerprofile=' + $scope.dockerprofile);
        console.log('name=' + $scope.container.Name);

        console.log(JSON.stringify($scope.schedule));
        usSpinnerService.spin('restore-spinner');

	//need ServerID, ContainerName, ProfileName, ScheduleID
	//
	var postMessage = {};
	postMessage.ServerID = $scope.container.ServerID;
	postMessage.ProjectID = $scope.container.ProjectID;
	postMessage.ContainerName = $scope.container.Name;
	postMessage.ProfileName = 'restore';
	postMessage.DockerProfile = $scope.dockerprofile;
	postMessage.ScheduleID = $scope.schedule.ID;
	postMessage.StatusID = $stateParams.statusID;
        tasksFactory.execute(postMessage)
            .success(function(data) {
                $scope.alerts = [{
                    type: 'success',
                    msg: 'success'
                }];
                usSpinnerService.stop('restore-spinner');
            })
            .error(function(error) {
                $scope.alerts = [{
                    type: 'danger',
                    msg: error.Error
                }];
                usSpinnerService.stop('restore-spinner');
                console.log('here is an error ' + error.Error);
            });
    };
};

var ContainerScheduleHistoryDeleteController = function($scope, $stateParams, $state, tasksFactory, utils, usSpinnerService) {

    console.log('in schedule history Delete controller');
    console.log('params = ' + JSON.stringify($stateParams));
    console.log('scope container = ' + JSON.stringify($scope.container));
    console.log('scope container = ' + $scope.container["Name"]);
    $scope.deletestatusID = $stateParams.statusID;

    $scope.deletehistory = function() {
        console.log('deletehistory called to restore backup with schedule=' + $scope.schedule.ID);
        usSpinnerService.spin('restore-spinner');
	var postMessage = {};
	postMessage.StatusID = $scope.deletestatusID;
	console.log('calling statusdelete with StatusID=' + postMessage.StatusID);

        tasksFactory.deletetaskstatus(postMessage)
            .success(function(data) {
                console.log('successful post deletetaskstatus with data=' + data);
		   $state.go('projects.container.schedule', $stateParams, {
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

        usSpinnerService.stop('restore-spinner');
    };
};

var ContainerScheduleHistoryController = function($scope, $stateParams, $state, tasksFactory, utils) {

    $scope.edit = function() {
        $state.go('.edit', $stateParams);
    };

    $scope.deletehistory = function() {
        console.log('deletehistory not implemented yet');
    };

    $scope.restore = function() {
        //console.log('restore not implemented yet params=' + JSON.stringify($stateParams));
    };

    $scope.refresh = function() {
        tasksFactory.getallstatus($stateParams.scheduleID)
            .success(function(data) {
                $scope.stats = data;
		//console.log(JSON.stringify(data));
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

var ContainerScheduleController = function($scope, $stateParams, $state, tasksFactory, serversFactory, utils) {

    $scope.edit = function() {
        $state.go('.edit', $stateParams);
    };
    serversFactory.all()
        .success(function(data) {
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
    $scope.schedule = {};
    $scope.schedule.RestoreSet = 'latest';
    $scope.schedule.RestoreRemotePath = '';
    $scope.schedule.RestoreRemoteHost = '';
    $scope.schedule.RestoreRemoteUser = '';
    $scope.schedule.RestoreDbUser = '';
    $scope.schedule.RestoreDbPass = '';

    $scope.profiles = [{
        name: 'pg_basebackup'
    }, {
        name: 'pg_backrest_restore'
    }];
    $scope.currentProfileName = $scope.profiles[0];

    $scope.thething = [{
        'name': 'thething',
        'checked': false
    }];
    //console.log(JSON.stringify(servers));
    $scope.myServer = servers.data[0].IPAddress;
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
        $scope.schedule.Serverip = $scope.myServer;
        $scope.schedule.ProfileName = $scope.currentProfileName.name;
        //console.log('adding schedule with schedule=' + JSON.stringify($scope.schedule));
        tasksFactory.addschedule($scope.schedule, $scope.container.Name)
            .success(function(data) {
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
        name: 'pg_backrest_restore'
    }];
    $scope.currentProfileName = $scope.profiles[0];

    $scope.myServer = [];
    $scope.thething = [{
        'name': 'thething',
        'checked': false
    }];


    tasksFactory.getschedule($stateParams.scheduleID)
        .success(function(data) {
            //console.log(JSON.stringify(data));
            $scope.schedule = data;
            if ($scope.schedule.Enabled == 'YES') {
                $scope.thething.checked = true;
            } else {
                $scope.thething.checked = false;
            }

            if ($scope.schedule.ProfileName == 'pg_basebackup') {
                $scope.currentProfileName = $scope.profiles[0];
            } else {
                $scope.currentProfileName = $scope.profiles[1];
            }
            //$scope.setServer();
            console.log('setting server to ' + $scope.schedule.Serverip);
            $scope.myServer = $scope.schedule.Serverip;
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

	var postMessage = {};
	postMessage.ServerID = $scope.ServerID;
	postMessage.ContainerName = '';
	postMessage.ProfileName = $scope.ProfileName;
	postMessage.ProjectID = $scope.ProjectID;
	postMessage.ScheduleID = $scope.ScheduleID;
	console.log('calling executenow with ProjectID=' + $scope.ProjectID);

        tasksFactory.execute(postMessage)
            .success(function(data) {
                console.log('successful post execute with data=' + data);
                //console.log(JSON.stringify(data));
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

        tasksFactory.updateschedule($scope.schedule)
            .success(function(data) {
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


    $scope.scheduleID = $stateParams.scheduleID;

    $scope.delete = function() {
        tasksFactory.deleteschedule($stateParams.scheduleID)
            .success(function(data) {
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

    serversFactory.all()
        .success(function(data) {
            $scope.servers = data;
            newcontainer.ID = 0;
            newcontainer.Name = 'newcontainer';
            newcontainer.Image = 'cpm-node';
            $scope.dockerprofile = 'SM';
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
        $scope.container.ProjectID = $stateParams.projectId;

        $scope.container.ID = 0; //0 means to do an insert

        containersFactory.add($scope.container, $scope.dockerprofile)
            .success(function(data) {
                //console.log('data from add is ' + JSON.stringify(data));
                console.log('current projectid is ' + $stateParams.projectId);
                $state.transitionTo('projects.container.details', {
                    containerId: data.ID,
                    projectId: $stateParams.projectId
                }, {
                    reload: true,
                    inherit: false
                });
                usSpinnerService.stop('spinner-1');
                /**
		$state.go('projects.container.details', {
			containerId: data.ID,
			projectId: $stateParams.projectId
		});

                $state.go('projects.container', $stateParams, {
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
                usSpinnerService.stop('spinner-1');
            });
    };
};
