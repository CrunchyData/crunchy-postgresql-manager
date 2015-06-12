angular.module('uiRouterSample.projects', [
    'ui.router',
    'ngCookies',
    'ui.bootstrap'
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

                            $scope.dataForTheTree = [{
                                "name": "Project 1",
                                "type": "project",
                                "id": "1",
                                "children": [{
                                    "name": "Clusters",
                                    "type": "label",
                                    "projectid": "1",
                                    "children": [{
                                        "name": "Jenifer",
                                        "type": "cluster",
                                        "id": "1",
                                        "projectid": "1",
                                        "children": []
                                    }]
                                }, {
                                    "name": "Databases",
                                    "type": "label",
                                    "id": "3",
                                    "projectid": "1",
                                    "children": [{
                                        "name": "One",
                                        "type": "Database",
                                        "id": "23",
                                        "projectid": "1",
                                        "children": []
                                    }, {
                                        "name": "Two",
                                        "type": "Database",
                                        "id": "28",
                                        "projectid": "1",
                                        "children": []
                                    }]
                                }]
                            }];

                            $scope.loadTree = function(projects) {
                                var data = [];
                                angular.forEach(projects, function(d) {
                                    //load databases
                                    var databases = [];
                                    angular.forEach(d.Containers, function(name, id) {
                                        databases.push({
                                            type: 'database',
                                            name: name,
                                            id: id
                                        });
                                    });

                                    //load clusters
                                    var clusters = [];
                                    angular.forEach(d.Clusters, function(name, id) {
                                        clusters.push({
                                            type: 'cluster',
                                            name: name,
                                            id: id
                                        });
                                    });

                                    //load project children
                                    var childs = [];
                                    childs.push({
                                        type: 'label',
                                        name: 'Clusters',
                                        projectid: d.ID,
                                        children: clusters

                                    });
                                    childs.push({
                                        type: 'label',
                                        name: 'Databases',
                                        projectid: d.ID,
                                        children: databases
                                    });

                                    //load project
                                    data.push({
                                        name: d.Name,
                                        id: d.ID,
                                        type: 'project',
                                        children: childs
                                    });
                                });
                                console.log('data is ' + JSON.stringify(data));
                                $scope.dataForTheTree = data;
                            };

                            $scope.projects = projects;
                            $scope.loadTree(projects.data);

                            $scope.showSelected = function(node) {
                                console.log('tree ' + node.name + ' selected id=' + node.id + ' type=' + node.type + ' projectid=' + node.projectid);
                                if (node.type == 'database') {
                                    $state.go('projects.container.details', {
                                        containerId: node.id
                                    });
                                } else if (node.type == 'cluster') {
                                    $state.go('projects.cluster.details', {
                                        clusterId: node.id
                                    });
                                } else if (node.type == 'project') {
                                    $state.go('projects.detail.details', {
                                        projectId: node.id
                                    });
                                } else if (node.type == 'label') {
                                    $state.go('projects.detail.details', {
                                        projectId: node.projectid
                                    });
                                }
                            };

                            $scope.goToFirst = function() {
                                console.log('projects=' + JSON.stringify($scope.projects));
                                var randId = $scope.projects.data[0].ID;

                                $state.go('projects.detail', {
                                    projectId: randId
                                });
                            };
                            $scope.goToFirst();


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
                            var randId = $scope.projects.data[0].ID;
                            $state.go('projects.detail', {
                                projectId: randId
                            });
                        };

                        $scope.goToFirst();
                    }
                ]
            })

            .state('projects.container', {

                url: '/{containerId}',

                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'containersFactory',
                            function($scope, $state, $cookieStore, $stateParams, utils, containersFactory) {
                                console.log('in projects.container with containerId ' + JSON.stringify($stateParams));
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in projects');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                containersFactory.get($stateParams.containerId)
                                    .success(function(data) {
                                        $scope.container = data;
                                        console.log('success with container get');
					console.log('container ' + JSON.stringify($scope.container));
                                    }).error(function(error) {
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

            .state('projects.container.details', {
                url: '/details/:itemId',
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

            .state('projects.container.schedule.add', {
                url: '',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.schedule.add.html',
                        resolve: {
                            servers: ['serversFactory',
                                function(serversFactory) {
                                    console.log('in the resolv of servers');
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
                url: '/delete/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.delete.html',
                        controller: ContainerDeleteController
                    },
                }
            })

            .state('projects.container.stop', {
                url: '/stop/:itemId',
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
                url: '/taskschedules/:scheduleID',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.taskschedules.delete.html',
                        controller: ContainerScheduleDeleteController
                    },
                }
            })

            .state('projects.container.accessrules', {
                url: '/accessrules/:itemId',
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
                url: '/users/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.edit.html',
                        controller: ContainerUsersEditController
                    },
                }
            })

            .state('projects.container.users.add', {
                url: '/users/add/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.add.html',
                        controller: ContainerUsersAddController
                    },
                }
            })

            .state('projects.container.users.delete', {
                url: '/users/delete/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.users.delete.html',
                        controller: ContainerUsersDeleteController
                    },
                }
            })

            .state('projects.container.monitor', {
                url: '/monitor/:itemId',
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
                url: '/monitor/pgstatdatabase/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgstatdatabase.html',
                        controller: ContainerMonitorpgstatdatabaseController
                    },
                }
            })

            .state('projects.container.monitor.bgwriter', {
                url: '/monitor/bgwriter/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.bgwriter.html',
                        controller: ContainerMonitorbgwriterController
                    },
                }
            })

            .state('projects.container.monitor.pgstatreplication', {
                url: '/monitor/pgstatreplication/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgstatreplication.html',
                        controller: ContainerMonitorpgstatreplicationController
                    },
                }
            })

            .state('projects.container.monitor.pgsettings', {
                url: '/monitor/pgsettings/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgsettings.html',
                        controller: ContainerMonitorpgsettingsController
                    },
                }
            })

            .state('projects.container.monitor.pgcontroldata', {
                url: '/monitor/pgcontroldata/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.pgcontroldata.html',
                        controller: ContainerMonitorpgcontroldataController
                    },
                }
            })


            .state('projects.container.monitor.loadtest', {
                url: '/monitor/loadtest/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.loadtest.html',
                        controller: ContainerMonitorloadtestController
                    },
                }
            })

            .state('projects.container.monitor.databasesize', {
                url: '/monitor/databasesize/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.container.monitor.databasesize.html',
                        controller: ContainerMonitordatabasesizeController
                    },
                }
            })

            .state('projects.cluster', {

                url: '/{clusterId}',

                views: {

                    '': {
                        templateUrl: 'app/projects/projects.cluster.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'clustersFactory',
                            function($scope, $state, $cookieStore, $stateParams, utils, clustersFactory) {
                                console.log('in projects.cluster with clusterId ' + JSON.stringify($stateParams));
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in projects');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                clustersFactory.get($stateParams.clusterId)
                                    .success(function(data) {
                                        $scope.cluster = data;
                                        console.log('success with cluster get');
                                    }).error(function(error) {
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

            .state('projects.cluster.details', {
                url: '/details/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/projects/projects.cluster.details.html',
                        controller: MyController
                    },
                }
            })

            .state('projects.cluster.delete', {

                url: '/delete/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.cluster.delete.html',
                        controller: ClusterDeleteController
                    },
                }
            })



            .state('projects.detail', {

                url: '/{projectId:[0-9]{1,4}}',

                views: {

                    '': {
                        templateUrl: 'app/projects/projects.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, utils) {
                                console.log('in projects.detail');
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in projects');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                if ($scope.projects.data.length > 0) {
                                    angular.forEach($scope.projects.data, function(item) {
                                        if (item.ID == $stateParams.projectId) {
                                            $scope.project = item;
                                            $state.go('projects.detail.edit', {
                                                projectId: item.ID
                                            });
                                        }
                                    });
                                }

                            }
                        ]
                    },

                }
            })

            .state('projects.detail.details', {
                url: '/details/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.detail.details.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                console.log('in detail.details');

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('projects.detail.add', {
                url: '/add/:itemId',
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
                                            console.log('success with add');
                                            $state.go('projects.list', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });
                                        }).error(function(error) {
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

            .state('projects.detail.delete', {
                url: '/delete/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.detail.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'projectsFactory', 'utils',
                            function($scope, $stateParams, $state, projectsFactory, utils) {

                                $scope.delete = function() {

                                    projectsFactory.delete($scope.project.ID)
                                        .success(function(data) {
                                            console.log('success with delete');
                                            $state.go('projects.list', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });

                                        }).error(function(error) {
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

            .state('projects.detail.edit', {
                url: '/edit/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/projects/projects.detail.edit.html',
                        controller: ['$rootScope', '$scope', '$stateParams', 'projectsFactory', '$state', 'utils',
                            function($rootScope, $scope, $stateParams, projectsFactory, $state, utils) {
                                console.log('edit controller called');
                                console.log('projectId=' + $stateParams.projectId);
                                $rootScope.projectId = $stateParams.projectId;

                                console.log('in detail.edit');
                                $scope.save = function() {
                                    projectsFactory.update($scope.project)
                                        .success(function(data) {
                                            console.log('success with update');
                                            $scope.alerts = [{
                                                type: 'success',
                                                msg: 'success'
                                            }];
                                        }).error(function(error) {
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
        }
    ]
);
