angular.module('uiRouterSample.projects', [
    'ui.router',
    'ngCookies',
    'ui.bootstrap',
    'angularSpinners'
])

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('projects', {

                    abstract: true,

                    url: '/projects',

                    templateUrl: 'app/projects/projects.html',

                    resolve: {
                        projects: ['$cookieStore', 'projectsFactory',
                            function($cookieStore, projectsFactory) {

                                if (!$cookieStore.get('cpm_token')) {
                                    var nothing = [];
                                    console.log('returning nothing');
                                    return nothing;
                                }

                                return projectsFactory.all();
                            }
                        ]
                    },

                    controller: ['$scope', '$cookieStore', '$state', 'projects', 'utils',
                        function($scope, $cookieStore, $state, projects, utils) {

                            if (!$cookieStore.get('cpm_token')) {
                                $state.go('login', {
                                    userId: 'hi'
                                });
                            }

                            $scope.treeOptions = {
                                    nodeChildren: "children",
                                    dirSelectable: true,
                                    multiSelection: false,
                                    injectClasses: {
                                        ul: "a1",
                                        li: "a2",
                                        liSelected: "a7",
                                        iExpanded: "a3",
                                        iCollapsed: "a4",
                                        iLeaf: "a5",
                                        label: "a6",
                                        labelSelected: "a8"
                                    }
                                }

                            $scope.projects = projects.data;
                            $scope.dataForTheTree = projects.data;
                            //$scope.loadTree(projects.data);


                            $scope.showSelected = function(node) {
			    	$scope.lastselectednode = node;
                                console.log('tree ' + node.name + ' selected id=' + node.id + ' type=' + node.type + ' projectid=' + node.projectid);
                                console.log('currently selected=' + JSON.stringify($scope.node1));
                                if (node.type == 'database') {
                                    $state.go('projects.container.details', {
                                        containerId: node.id,
                                        projectId: node.projectid
                                    });
                                } else if (node.type == 'cluster') {
                                    $state.go('projects.cluster.details', {
                                        clusterId: node.id
                                    });
                                } else if (node.type == 'proxy') {
                                    $state.go('projects.proxy.details', {
                                        containerId: node.id,
                                        containerName: node.name,
                                        projectId: node.projectid
                                    });
                                } else if (node.type == 'project') {
                                    $state.go('projects.detail.edit', {
                                        projectId: node.id
                                    });
                                } else if (node.type == 'label') {
                                    $state.go('projects.detail.edit', {
                                        projectId: node.projectid
                                    });
                                }
                            };

                            $scope.goToFirst = function() {
                                    //console.log('jeff projects.data00=' + JSON.stringify(projects.data[0]));
                                    //console.log('jeff projects.name=' + projects.data[0].name);
                                    //console.log('jeff projects.id=' + projects.data[0].id);
                                    var randId = projects.data[0].id;

                                    $state.go('projects.detail', {
                                        projectId: randId
                                    });
                            };
                            //$scope.goToFirst();
                        }
                    ]
                })

            .state('projects.list', {

                url: '',

                templateUrl: 'app/projects/projects.list.html',

                controller: ['$scope', '$cookieStore', '$state', 'projects', 'utils',
                    function($scope, $cookieStore, $state, projects, utils) {

                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                userId: 'hi'
                            });
                        }

                        $scope.projects = projects;

                        $scope.goToFirst = function() {
                            var randId = $scope.projects.data[0].id;
                            $state.go('projects.detail', {
                                projectId: randId
                            });
                        };

                        $scope.goToFirst();
                    }
                ]
            })

            .state('projects.container', {

                url: '/{projectId}/container/{containerId}',

                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'containersFactory',
                            function($scope, $state, $cookieStore, $stateParams, utils, containersFactory) {
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in projects');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                if ($stateParams.containerId != "") {
                                    containersFactory.get($stateParams.containerId)
                                        .success(function(data) {
                                            $scope.container = data;
                                        }).error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                            console.log('here is an error ' + error.Error);
                                        });
                                }

                            }
                        ]
                    },

                }
            })

            .state('projects.detail.addcontainer', {
                url: '/addcontainer',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.add.html',
                        controller: ContainerAddController
                    },
                }
            })

            .state('projects.container.details', {
                url: '/details',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.details.html',
                        controller: ContainerDetailController
                    },
                }
            })

            .state('projects.container.schedule', {
                url: '/schedule/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.html',
                        controller: ContainerScheduleController
                    },
                }
            })

            .state('projects.container.schedule.delete', {
                url: '/schedule/:schedulID',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.delete.html',
                        controller: ContainerScheduleDeleteController
                    },
                }
            })

            .state('projects.container.schedule.edit', {
                url: '',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.edit.html',
                        resolve: {
                            servers: ['serversFactory',
                                function(serversFactory) {
                                    serversFactory.all()
                                        .success(function(data) {
                                            servers = data;
                                            return data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                            console.log('here is an error ' + error.Error);
                                        });
                                }
                            ]
                        },

                        controller: ContainerScheduleEditController
                    },
                }
            })

            .state('projects.container.scheduleadd', {
                url: '',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.scheduleadd.html',
                        resolve: {
                            servers: ['serversFactory',
                                function(serversFactory) {
                                    return serversFactory.all();
                                }
                            ]
                        },

                        controller: ContainerScheduleAddController
                    },
                }
            })


            .state('projects.container.schedule.execute', {
                url: '/schedule/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.execute.html',
                        controller: ContainerScheduleExecuteController
                    },
                }
            })

            .state('projects.container.schedule.history', {
                url: '/schedule/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.history.html',
                        controller: ContainerScheduleHistoryController
                    },
                }
            })

            .state('projects.container.start', {
                url: '/start/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.start.html',
                        controller: ContainerStartController
                    },
                }
            })

            .state('projects.container.delete', {
                url: '/delete',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.delete.html',
                        controller: ContainerDeleteController
                    },
                }
            })

            .state('projects.container.failover', {
                url: '/failover/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.failover.html',
                        controller: ContainerFailoverController
                    },
                }
            })

            .state('projects.container.stop', {
                url: '/stop',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.stop.html',
                        controller: ContainerStopController
                    },
                }
            })

            .state('projects.container.taskschedules', {
                url: '/taskschedules/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.taskschedules.html',
                        controller: ContainerTaskSchedulesController
                    },
                }
            })

            .state('projects.container.taskschedules.delete', {
                url: '/delete',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.taskschedules.delete.html',
                        controller: ContainerScheduleDeleteController
                    },
                }
            })

            .state('projects.container.accessrules', {
                url: '/accessrules',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.accessrules.html',
                        controller: ContainerAccessRulesController
                    },
                }
            })

            .state('projects.container.users', {
                url: '/users/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.html',
                        controller: ContainerUsersController
                    },
                }
            })

            .state('projects.container.users.edit', {
                url: '/edit',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.edit.html',
                        controller: ContainerUsersEditController
                    },
                }
            })

            .state('projects.container.usersadd', {
                url: '/users/add/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.add.html',
                        controller: ContainerUsersAddController
                    },
                }
            })

            .state('projects.container.users.delete', {
                url: '/delete',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.delete.html',
                        controller: ContainerUsersDeleteController
                    },
                }
            })

            .state('projects.container.monitor', {
                url: '/monitor',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.html',
                        controller: ['$scope', '$stateParams', '$state', 'containersFactory', 'utils',
                            function($scope, $stateParams, $state, containersFactory, utils) {}
                        ]

                    },
                }
            })

            .state('projects.container.monitor.pgstatdatabase', {
                url: '/pgstatdatabase',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgstatdatabase.html',
                        controller: ContainerMonitorpgstatdatabaseController
                    },
                }
            })

            .state('projects.container.monitor.bgwriter', {
                url: '/bgwriter',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.bgwriter.html',
                        controller: ContainerMonitorbgwriterController
                    },
                }
            })

            .state('projects.container.monitor.badger', {
                url: '/badger/:containerName',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.badger.html',
                        controller: ContainerMonitorbadgerController
                    },
                }
            })

            .state('projects.container.monitor.pgstatreplication', {
                url: '/pgstatreplication',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgstatreplication.html',
                        controller: ContainerMonitorpgstatreplicationController
                    },
                }
            })

            .state('projects.container.monitor.pgsettings', {
                url: '/pgsettings',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgsettings.html',
                        controller: ContainerMonitorpgsettingsController
                    },
                }
            })

            .state('projects.container.monitor.pgstatstatements', {
                url: '/pgstatements',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgstatstatements.html',
                        controller: ContainerMonitorpgstatstatementsController
                    },
                }
            })


            .state('projects.container.monitor.pgcontroldata', {
                url: '/pgcontroldata',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgcontroldata.html',
                        controller: ContainerMonitorpgcontroldataController
                    },
                }
            })


            .state('projects.container.monitor.loadtest', {
                url: '/loadtest',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.loadtest.html',
                        controller: ContainerMonitorloadtestController
                    },
                }
            })

            .state('projects.container.monitor.databasesize', {
                url: '/databasesize',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.databasesize.html',
                        controller: ContainerMonitordatabasesizeController
                    },
                }
            })

            .state('projects.cluster', {

                url: '/cluster/{clusterId}',

                views: {

                    '': {
                        templateUrl: 'app/projects/projects.cluster.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'clustersFactory',
                            function($scope, $state, $cookieStore, $stateParams, utils, clustersFactory) {
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in projects');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                clustersFactory.get($stateParams.clusterId)
                                    .success(function(data) {
                                        $scope.cluster = data;
                                    }).error(function(error) {
                                        $scope.alerts = [{
                                            type: 'danger',
                                            msg: error.Error
                                        }];
                                        console.log('here is an error ' + error.Error);
                                    });

                            }
                        ]
                    },

                }
            })

            .state('projects.cluster.details', {
                url: '/details',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.cluster.details.html',
                        controller: MyController
                    },
                }
            })

            .state('projects.cluster.delete', {

                url: '/delete',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.cluster.delete.html',
                        controller: ClusterDeleteController
                    },
                }
            })

            .state('projects.cluster.start', {

                url: '/start',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.cluster.start.html',
                        controller: ClusterStartController
                    },
                }
            })

            .state('projects.cluster.stop', {

                url: '/stop',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.cluster.stop.html',
                        controller: ClusterStopController
                    },
                }
            })

            .state('projects.cluster.scale', {

                url: '/scale',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.cluster.scale.html',
                        controller: ClusterScaleController
                    },
                }
            })


            .state('projects.proxy', {

                url: '/{projectId}/proxy/{containerId}',

                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'proxyFactory',
                            function($scope, $state, $cookieStore, $stateParams, utils, proxyFactory) {
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in projects');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                if ($stateParams.proxyId != "") {
                                    proxyFactory.getbycontainerid($stateParams.containerId)
                                        .success(function(data) {
                                            $scope.proxy = data;
                                        }).error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                            console.log('here is an error ' + error.Error);
                                        });
                                }

                            }
                        ]
                    },

                }
            })

            .state('projects.proxy.details', {
                url: '/details/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.detail.html',
                        controller: ProxyDetailController
                    },
                }
            })

            .state('projects.proxy.schedule', {
                url: '/schedule/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.schedule.html',
                        controller: ProxyScheduleController
                    },
                }
            })

            .state('projects.proxy.schedule.edit', {
                url: '',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.edit.html',
                        resolve: {
                            servers: ['serversFactory',
                                function(serversFactory) {
                                    serversFactory.all()
                                        .success(function(data) {
                                            servers = data;
                                            return data;
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                            console.log('here is an error ' + error.Error);
                                        });
                                }
                            ]
                        },

                        controller: ContainerScheduleEditController
                    },
                }
            })


            .state('projects.proxy.schedule.history', {
                url: '/schedule/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.history.html',
                        controller: ContainerScheduleHistoryController
                    },
                }
            })

            .state('projects.proxy.schedule.execute', {
                url: '/schedule/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.execute.html',
                        controller: ContainerScheduleExecuteController
                    },
                }
            })
            .state('projects.proxy.schedule.delete', {
                url: '/schedule/:schedulID',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.delete.html',
                        controller: ContainerScheduleDeleteController
                    },
                }
            })

            .state('projects.proxy.scheduleadd', {
                url: '',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.scheduleadd.html',
                        resolve: {
                            servers: ['serversFactory',
                                function(serversFactory) {
                                    return serversFactory.all();
                                }
                            ]
                        },

                        controller: ProxyScheduleAddController
                    },
                }
            })


            .state('projects.proxy.taskschedules', {
                url: '/taskschedules/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.taskschedules.html',
                        controller: ProxyTaskSchedulesController
                    },
                }
            })
            .state('projects.proxy.taskschedules.delete', {
                url: '/taskschedules/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.taskschedules.delete.html',
                        controller: ProxyScheduleDeleteController
                    },
                }
            })

            .state('projects.proxy.users', {
                url: '/users/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.users.html',
                        controller: ContainerUsersController
                    },
                }
            })

            .state('projects.proxy.users.edit', {
                url: '/users/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.edit.html',
                        controller: ProxyUsersEditController
                    },
                }
            })
            .state('projects.proxy.users.add', {
                url: '/users/add/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.add.html',
                        controller: ProxyUsersAddController
                    },
                }
            })

            .state('projects.proxy.users.delete', {
                url: '/users/delete/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.delete.html',
                        controller: ProxyUsersDeleteController
                    },
                }
            })

            .state('projects.proxy.delete', {
                url: '/delete/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.delete.html',
                        controller: ProxyDeleteController
                    },
                }
            })

            .state('projects.proxy.databasesize', {
                url: '/databasesize/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.databasesize.html',
                        controller: ProxyDatabasesizeController
                    },
                }
            })

            .state('projects.proxy.bgwriter', {
                url: '/bgwriter/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.bgwriter.html',
                        controller: ContainerMonitorbgwriterController
                    },
                }
            })
            .state('projects.proxy.loadtest', {
                url: '/loadtest',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.loadtest.html',
                        controller: ContainerMonitorloadtestController
                    },
                }
            })
            .state('projects.proxy.pgsettings', {
                url: '/pgsettings',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgsettings.html',
                        controller: ContainerMonitorpgsettingsController
                    },
                }
            })
            .state('projects.proxy.pgstatdatabase', {
                url: '/pgstatdatabase',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgstatdatabase.html',
                        controller: ContainerMonitorpgstatdatabaseController
                    },
                }
            })
            .state('projects.proxy.pgstatstatements', {
                url: '/pgstatements',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgstatstatements.html',
                        controller: ContainerMonitorpgstatstatementsController
                    },
                }
            })

            .state('projects.proxy.start', {
                url: '/start/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.start.html',
                        controller: ProxyStartController
                    },
                }
            })

            .state('projects.proxy.stop', {
                url: '/stop/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.stop.html',
                        controller: ProxyStopController
                    },
                }
            })


            .state('projects.detail.addproxy', {
                url: '/addproxy',
                views: {

                    '': {
                        templateUrl: 'app/proxy/proxy.add.html',
                        controller: ProxyAddController
                    },
                }
            })

            .state('projects.detail', {

                url: '/{projectId:[0-9]{1,4}}',

                views: {

                    '': {
                        templateUrl: 'app/projects/projects.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'projectsFactory',
                            function($scope, $state, $cookieStore, $stateParams, utils, projectsFactory) {
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in projects');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                projectsFactory.get($stateParams.projectId)
                                    .success(function(data) {
                                        $scope.project = data;
                                    }).error(function(error) {
                                        $scope.alerts = [{
                                            type: 'danger',
                                            msg: error.Error
                                        }];
                                        console.log('here is an error ' + error.Error);
                                    });

                            }

                        ]
                    },

                }
            })

            .state('projects.detail.gotocontainer', {

                url: '/container/{containerId}',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.detail.html',
                        controller: GotocontainerController
                    },
                }
            })

            .state('projects.detail.gotoproxy', {

                url: '/container/{containerId}',
                views: {
                    '': {
                        templateUrl: 'app/proxy/proxy.html',
                        controller: GotoproxyController
                    },
                }
            })

            .state('projects.detail.autocluster', {

                url: '/autocluster',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.detail.autocluster.html',
                        controller: ClusterAutoClusterController
                    },
                }
            })

            .state('projects.detail.details', {
                url: '/details/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.detail.details.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils', 'projectId',
                            function($scope, $stateParams, $state, utils) {
                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('projects.detail.addproject', {
                url: '/addproject',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.detail.add.html',
                        controller: ['$scope', '$stateParams', '$state', '$cookieStore', 'projectsFactory', 'utils', 'usSpinnerService',
                            function($scope, $stateParams, $state, $cookieStore, projectsFactory, utils, usSpinnerService) {
                                var newproject = {}
                                newproject.ID = '0';
                                newproject.Name = 'newproject';
                                newproject.Desc = 'new project description';
                                newproject.UpdateDate = '000';
                                newproject.Token = $cookieStore.get('cpm_token');
                                $scope.project = newproject;

                                $scope.create = function() {
                                    usSpinnerService.spin('spinner-1');
                                    projectsFactory.add($scope.project)
                                        .success(function(data) {
                                            usSpinnerService.stop('spinner-1');
                                            $state.go('projects.list', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });
                                        }).error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                            console.log('here is an error ' + error.Error);
                                            usSpinnerService.stop('spinner-1');
                                        });


                                };
                            }
                        ]
                    },
                }
            })

            .state('projects.detail.delete', {
                url: '/delete',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.detail.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'projectsFactory', 'utils',
                            function($scope, $stateParams, $state, projectsFactory, utils) {

                                $scope.delete = function() {

                                    projectsFactory.delete($scope.project.ID)
                                        .success(function(data) {
                                            $state.go('projects.list', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });

                                        }).error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                            console.log('here is an error ' + error.Error);
                                        });

                                };
                            }
                        ]
                    },
                }
            })

            .state('projects.detail.edit', {
                url: '/edit/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.detail.edit.html',
                        controller: ['$rootScope', '$scope', '$stateParams', 'projectsFactory', '$state', 'utils',
                            function($rootScope, $scope, $stateParams, projectsFactory, $state, utils) {
                                $rootScope.projectId = $stateParams.projectId;


                                projectsFactory.get($stateParams.projectId)
                                    .success(function(data) {
                                        $scope.project = data;
                                    }).error(function(error) {
                                        $scope.alerts = [{
                                            type: 'danger',
                                            msg: error.Error
                                        }];
                                        console.log('here is an error ' + error.Error);
                                    });

                                $scope.save = function() {
                                    projectsFactory.update($scope.project)
                                        .success(function(data) {
                                            $scope.alerts = [{
                                                type: 'success',
                                                msg: 'success'
                                            }];
                                        }).error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                            console.log('here is an error ' + error.Error);
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
