angular.module('uiRouterSample.containers', [
    'ui.router',
    'ui.bootstrap'
])

.config(
    ['$stateProvider', '$urlRouterProvider', 
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('containers', {

                    abstract: true,

                    url: '/containers',

                    templateUrl: 'app/containers/containers.html',

                    resolve: {
                        containers: ['$cookieStore', 'containersFactory',
                            function($cookieStore, containersFactory) {
                                console.log('resolving containers at the top');
                                if (!$cookieStore.get('cpm_token')) {
                                    var nothing = [];
                                    console.log('returning nothing');
                                    return nothing;
                                }

                                return containersFactory.all();
                            }
                        ]
                    },

                    controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'containers', 'utils',
                        function($scope, $state, $cookieStore, $stateParams, containers, utils) {

                            console.log('next 1');
                            if (!$cookieStore.get('cpm_token')) {
                                console.log('next 2');
                                $state.go('login', {
                                    userId: 'hi'
                                });
                            }
                            $scope.containers = containers;

                            $scope.goToFirst = function() {
                                console.log('go to first called in containers');
                                var randId = $scope.containers.data[0].ID;

                                $state.go('containers.detail.details', {
                                    containerId: randId
                                });
                            };
                            console.log('stateParams=' + JSON.stringify($stateParams));
                        }
                    ]
                })

            .state('containers.list', {

                url: '',

                templateUrl: 'app/containers/containers.list.html',
                resolve: {
                    containers: ['$cookieStore', 'containersFactory',
                        function($cookieStore, containersFactory) {
                            if (!$cookieStore.get('cpm_token')) {
                                var nothing = [];
                                console.log('returning nothing');
                                return nothing;
                            }

                            console.log('resolving in containers list the list of containers');
                            return containersFactory.all();
                        }
                    ]
                },
                controller: ['$cookieStore', '$scope', '$state', '$stateParams', 'containers', 'utils',
                    function($cookieStore, $scope, $state, $stateParams, containers, utils) {

                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                userId: 'hi'
                            });
                        }

                        $scope.containers = containers;

                        $scope.goToFirst = function() {
                            console.log('go to first called in containers');
                            if ($scope.containers.data.length > 0) {
                                var randId = $scope.containers.data[0].ID;
                                $state.go('containers.detail.details', {
                                    containerId: randId
                                });
                            }
                        };
                        $scope.goToFirst();
                    }
                ]
            })

            .state('containers.detail', {

                url: '/{containerId}',

                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'containersFactory', 'serversFactory', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, containersFactory, serversFactory, utils) {
                                if (!$cookieStore.get('cpm_token')) {
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                serversFactory.all()
                                    .success(function(data) {
                                        console.log('got me servers =' + data);
                                        $scope.servers = data;
                                    })
                                    .error(function(error) {
                                        $scope.alerts = [{
                                            type: 'danger',
                                            msg: error.message
                                        }];
                                        console.log('here is an error ' + error.message);
                                    });

                                console.log('containers len here is ' + $scope.containers.data.length);
                                if ($scope.containers.data.length > 0) {
                                    angular.forEach($scope.containers.data, function(item) {
                                        if (item.ID == $stateParams.containerId) {
                                            $scope.container = item;
                                        }
                                    });
                                }

                            }
                        ]
                    },

                }
            })

            .state('containers.detail.schedule.history', {

                url: '/history/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.schedule.history.html',
                        controller: ['$scope', '$stateParams', '$state', 'tasksFactory', 'utils',
                            function($scope, $stateParams, $state, tasksFactory, utils) {

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
                                            msg: error.message
                                        }];
                                        console.log('here is an error ' + error.message);
                                    });
                            }
                        ]
                    }

                }
            })

            .state('containers.detail.schedule.delete', {

                url: '/delete/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.schedule.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'tasksFactory', 'utils',
                            function($scope, $stateParams, $state, tasksFactory, utils) {

                                $scope.delete = function() {
                                    console.log('in schedule delete' + $stateParams.scheduleID);
                                    tasksFactory.deleteschedule($stateParams.scheduleID)
                                        .success(function(data) {
                                            console.log('successful deleteschedule with data=' + data);
                                            console.log(JSON.stringify(data));
                                            $scope.stats = data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };
                            }
                        ]
                    }

                }
            })

            .state('containers.detail.taskhistory.delete', {
                views: {
                    '@containers.detail.taskhistory': {
                        templateUrl: 'app/containers/containers.detail.taskhistory.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.container.items, $stateParams.itemId);
                                $scope.done = function() {
                                    // Go back up. '^' means up one. '^.^' would be up twice, to the grandparent.
                                    $state.go('^', $stateParams);
                                };
                            }
                        ]
                    }
                }
            })

            .state('containers.detail.taskschedules', {

                url: '/taskschedules/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.taskschedules.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {


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
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                                $scope.refresh();
                            }
                        ]
                    },
                }
            })

            .state('containers.detail.schedule', {
                url: '/schedule/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.schedule.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'tasksFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, tasksFactory, utils) {

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
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                                $scope.refresh = function() {
                                    tasksFactory.getschedule($stateParams.scheduleID)
                                        .success(function(data) {
                                            console.log('got the schedule');
                                            $scope.schedule = data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };
                                console.log('in schedule...');
                                if ($stateParams.scheduleID != -1) {
                                    $scope.refresh();
                                }
                            }
                        ]
                    },
                }
            })

            .state('containers.detail.schedule.edit', {
                views: {
                    '@containers.detail.schedule': {
                        templateUrl: 'app/containers/containers.detail.schedule.edit.html',
                        resolve: {
                            servers: ['serversFactory',
                                function(serversFactory) {
                                    console.log('in the resolv of servers');
                                    serversFactory.all()
                                        .success(function(data) {
                                            console.log('successful servers all with data=' + data);
                                            servers = data;
                                            return data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                }
                            ]
                        },
                        controller: ['$scope', '$filter', '$stateParams', '$state', 'tasksFactory', 'serversFactory', 'utils', 'servers',
                            function($scope, $filter, $stateParams, $state, tasksFactory, serversFactory, utils, servers) {
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

                                if ($stateParams.scheduleID == -1) {
                                    console.log('jeff');
                                    JSON.stringify($scope.servers);
                                    $scope.myServer = $scope.servers[0];
                                    console.log('setting myServer to ' + $scope.myServer.Name);
                                } else {
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
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                }

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
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
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
                                    if ($stateParams.scheduleID == -1) {
                                        $scope.schedule.ServerID = $scope.myServer.ID;
                                        console.log('jeff value of profile is...' + $scope.currentProfileName.name);
                                        $scope.schedule.ProfileName = $scope.currentProfileName.name;
                                        tasksFactory.addschedule($scope.schedule, $scope.container.Name)
                                            .success(function(data) {
                                                console.log('successful add schedule with data=' + data);
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
                                            });
                                    } else {

                                        tasksFactory.updateschedule($scope.schedule)
                                            .success(function(data) {
                                                console.log('successful post schedule with data=' + data);
                                                console.log(JSON.stringify(data));
                                            })
                                            .error(function(error) {
                                                $scope.alerts = [{
                                                    type: 'danger',
                                                    msg: error.message
                                                }];
                                                console.log('here is an error ' + error.message);
                                            });
                                    }

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


                            }
                        ]
                    }
                }
            })

            .state('containers.detail.users', {

                url: '/users/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.users.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {
                                containersFactory.getallusers($stateParams.containerId)
                                    .success(function(data) {
                                        console.log('successful get with data=' + data);
                                        $scope.users = data;
                                    })
                                    .error(function(error) {
                                        $scope.alerts = [{
                                            type: 'danger',
                                            msg: error.message
                                        }];
                                        console.log('here is an error ' + error.message);
                                    });

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('containers.detail.users.edit', {
                views: {
                    '@containers.detail.users': {
                        templateUrl: 'app/containers/containers.detail.users.edit.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils', 
                            function($scope, $stateParams, $state, containersFactory, utils ) {
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
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                $scope.save = function() {
                                    console.log('save called');
			    	$scope.user.ContainerID = $stateParams.containerId;
				console.log('saving Rolsuper is ' + $scope.user.Rolsuper);
                                    containersFactory.updateuser($scope.user)
                                        .success(function(data) {
                                            console.log('successful updateuser with data=' + data);
                                            $state.go('containers.detail.users', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });

                                        })
                                        .error(function(error) {
						console.log(JSON.stringify(error));
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });

                                };
                            }
                        ]
                    }
                }
            })

            .state('containers.detail.users.add', {
                views: {
                    '@containers.detail.users': {
                        templateUrl: 'app/containers/containers.detail.users.add.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {
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
                                            $state.go('containers.detail.users', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });

                                        })
                                        .error(function(error) {
						console.log(JSON.stringify(error));
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });

                                };
                            }
                        ]
                    }
                }
            })

            .state('containers.detail.users.delete', {
                views: {
                    '@containers.detail.users': {
                        templateUrl: 'app/containers/containers.detail.users.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils', 
                            function($scope, $stateParams, $state, containersFactory, utils) {
                                	console.log('before doing delete user id=' + $stateParams.containerId + ' name=' + $stateParams.itemId);
					$scope.rolname = $stateParams.itemId;

                                $scope.delete = function() {
                                	console.log('doing delete user id=' + $stateParams.containerId + ' name=' + $scope.rolname);
                                    containersFactory.deleteuser($stateParams.containerId, $scope.rolname)
                                        .success(function(data) {
                                            console.log('successful deleteuser with data=' + data);
                                            $state.go('containers.detail.users', $stateParams, {
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
                                        });

                                };
                            }
                        ]
                    }
                }
            })

            .state('containers.detail.accessrules', {

                url: '/accessrules/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.accessrules.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.container.items, $stateParams.itemId);

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('containers.detail.accessrules.edit', {
                views: {
                    '@containers.detail.accessrules': {
                        templateUrl: 'app/containers/containers.detail.accessrules.edit.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.container.items, $stateParams.itemId);
                                $scope.done = function() {
                                    $state.go('^', $stateParams);
                                };
                            }
                        ]
                    }
                }
            })

            .state('containers.detail.monitor', {
                url: '/monitor/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/containers/containers.detail.monitor.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {}
                        ]
                    },
                }
            })

            .state('containers.detail.monitor.pgstatdatabase', {
                url: '/monitor/pgstatdatabase/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/containers/containers.detail.monitor.pgstatdatabase.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {
                                $scope.refresh = function() {
                                    containersFactory.pgstatdatabase($stateParams.containerId)
                                        .success(function(data) {
                                            console.log('successful get with data=' + data);
                                            $scope.statdbresults = data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                                $scope.refresh();

                            }
                        ]
                    },
                }
            })

            .state('containers.detail.monitor.pgstatreplication', {
                url: '/monitor/pgstatreplication/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/containers/containers.detail.monitor.pgstatreplication.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {
                                $scope.refresh = function() {
                                    containersFactory.pgstatreplication($stateParams.containerId)
                                        .success(function(data) {
                                            console.log('successful get with data=' + data);
                                            $scope.statreplresults = data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                                $scope.refresh();

                            }
                        ]
                    },
                }
            })

            .state('containers.detail.monitor.bgwriter', {
                url: '/monitor/bgwriter/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/containers/containers.detail.monitor.bgwriter.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {
                                $scope.refresh = function() {
                                    containersFactory.bgwriter($stateParams.containerId)
                                        .success(function(data) {
                                            console.log('successful get with data=' + data);
                                            $scope.bgwriterresults = data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                                $scope.refresh();
                            }
                        ]
                    },
                }
            })

            .state('containers.detail.monitor.pgsettings', {
                url: '/monitor/pgsettings/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/containers/containers.detail.monitor.pgsettings.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {
                                $scope.refresh = function() {
                                    containersFactory.pgsettings($stateParams.containerId)
                                        .success(function(data) {
                                            console.log('successful get with data=' + data);
                                            $scope.settingsresults = data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                                $scope.refresh();
                            }
                        ]
                    },
                }
            })

            .state('containers.detail.monitor.pgcontroldata', {
                url: '/monitor/pgcontroldata/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/containers/containers.detail.monitor.pgcontroldata.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {
                                $scope.refresh = function() {
                                    containersFactory.pgcontroldata($stateParams.containerId)
                                        .success(function(data) {
                                            console.log('successful get with data=' + data);
                                            $scope.controldataresults = data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                                $scope.refresh();
                            }
                        ]
                    },
                }
            })

            .state('containers.detail.monitor.loadtest', {
                url: '/monitor/loadtest/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/containers/containers.detail.monitor.loadtest.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {
                                $scope.slidervaluehigh = "10000";
                                $scope.slidervaluelow = "1000";
                                $scope.slidervalue = "1000";

                                $scope.refresh = function() {
                                    containersFactory.loadtest($stateParams.containerId, $scope.slidervalue)
                                        .success(function(data) {
                                            console.log('successful get with data=' + data);
                                            $scope.loadtestresults = data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                            }
                        ]
                    },
                }
            })

            .state('containers.detail.monitor.databasesize', {
                url: '/monitor/databasesize/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/containers/containers.detail.monitor.databasesize.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.container.items, $stateParams.itemId);

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })


            .state('containers.detail.details', {

                url: '/details/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.details.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils', 'containersFactory',
                            function($scope, $stateParams, $state, utils, containersFactory) {
console.log('containers details called containerId ' + $stateParams.containerId);
                                    containersFactory.get($stateParams.containerId)
                                        .success(function(data) {
                                            console.log('successful get with data=' + data);
                                            $scope.container=data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
			    }
                        ]
                    },
                }
            })

            .state('containers.detail.add', {

                url: '/add/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.add.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'containersFactory', 'utils', 'usSpinnerService',
                            function($scope, $stateParams, $state, serversFactory, containersFactory, utils, usSpinnerService) {

                                var newcontainer = {};

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
                                        $scope.container = newcontainer;
                                    })
                                    .error(function(error) {
                                        $scope.alerts = [{
                                            type: 'danger',
                                            msg: error.message
                                        }];
                                    });


                                $scope.add = function() {
                                    usSpinnerService.spin('spinner-1');
                                    $scope.container.ServerID = $scope.selectedServer.ID;

                                    $scope.container.ID = 0; //0 means to do an insert

                                    containersFactory.add($scope.container, $scope.standalone, $scope.dockerprofile)
                                        .success(function(data) {
                                            console.log('successful add with data=' + data);
                                            usSpinnerService.stop('spinner-1');
                                            $state.go('containers.list', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            usSpinnerService.stop('spinner-1');
                                        });
                                };
                            }
                        ]
                    },
                }
            })

            .state('containers.detail.delete', {

                url: '/delete/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils', 'usSpinnerService',
                            function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
                                var container = $scope.container;

                                $scope.delete = function() {
                                    usSpinnerService.spin('spinner-1');
                                    containersFactory.delete($stateParams.containerId)
                                        .success(function(data) {
                                            console.log('successful delete with data=' + data);
                                            usSpinnerService.stop('spinner-1');
                                            $state.go('containers.list', $stateParams, {
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

                            }
                        ]
                    },
                }
            })

            .state('containers.detail.start', {

                url: '/start/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.start.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils', 'usSpinnerService',
                            function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
                                var container = $scope.container;
                                console.log('here in start top');

                                $scope.start = function() {
                                    usSpinnerService.spin('spinner-1');
                                    containersFactory.start($stateParams.containerId)
                                        .success(function(data) {
                                            console.log('successful start with data=' + data);
                                            $state.go('containers.list', $stateParams);
                                            usSpinnerService.stop('spinner-1');
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
                            }
                        ]
                    },
                }
            })

            .state('containers.detail.stop', {

                url: '/stop/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/containers/containers.detail.stop.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils', 'usSpinnerService',
                            function($scope, $stateParams, $state, containersFactory, utils, usSpinnerService) {
                                var container = $scope.container;

                                $scope.stop = function() {
                                    usSpinnerService.spin('spinner-1');
                                    containersFactory.stop($stateParams.containerId)
                                        .success(function(data) {
                                            console.log('successful stop with data=' + data);
                                            $state.go('containers.list', $stateParams);
                                            usSpinnerService.stop('spinner-1');
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
                            }
                        ]
                    },
                }
            })
        }
    ]
);
