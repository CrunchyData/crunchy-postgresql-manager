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
               	projectId:  $stateParams.projectId
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
        containersFactory.start($stateParams.containerId)
            .success(function(data) {
                $state.go('projects.container.details', $stateParams, {
			reload: true,
			inherit: false
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


var ContainerMonitorbadgerController = function($sce, $scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
    $scope.badgerreportlink = $sce.trustAsResourceUrl(
        'http://' + $stateParams.containerName + ':10001/static/badger.html?rand=' + Math.round(Math.random() * 10000000));
    $scope.refresh = function() {
        usSpinnerService.spin('spinner-1');
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

    tasksFactory.execute($scope.schedule)
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

var ContainerScheduleHistoryController = function($scope, $stateParams, $state, tasksFactory, utils) {

    $scope.edit = function() {
        $state.go('.edit', $stateParams);
    };
    tasksFactory.getallstatus($stateParams.scheduleID)
        .success(function(data) {
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
        name: 'pg_other'
    }];
    $scope.currentProfileName = $scope.profiles[0];

    $scope.myServer = [];
    $scope.thething = [{
        'name': 'thething',
        'checked': false
    }];


    tasksFactory.getschedule($stateParams.scheduleID)
        .success(function(data) {
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

        containersFactory.add($scope.container, $scope.standalone, $scope.dockerprofile)
            .success(function(data) {
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
