angular.module('uiRouterSample.servers', [
    'ui.router',
    'ui.bootstrap'
])

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('servers', {

                    abstract: true,

                    url: '/servers',

                    templateUrl: 'app/servers/servers.html',

                    resolve: {
                        servers: ['$cookieStore', 'serversFactory',
                            function($cookieStore, serversFactory) {
                                if (!$cookieStore.get('cpm_token')) {
                                    var nothing = [];
                                    console.log('returning nothing');
                                    return nothing;
                                }

                                return serversFactory.all();
                            }
                        ]
                    },

                    controller: ['$scope', '$state', '$cookieStore', 'utils', 'servers',
                        function($scope, $state, $cookieStore, utils, servers) {

                            if (!$cookieStore.get('cpm_token')) {
                                $state.go('login', {
                                    userId: 'hi'
                                });
                            }

                            $scope.servers = servers;

                            $scope.goToFirst = function() {
                                if ($scope.servers.data.length > 0) {
                                    console.log($scope.servers.data[0].ID);
                                    var randId = $scope.servers.data[0].ID;

                                    $state.go('servers.detail.details', {
                                        serverId: randId
                                    });
                                }
                            };
                            $scope.goToFirst();
                        }
                    ]
                })

            .state('servers.list', {

                url: '',

                templateUrl: 'app/servers/servers.list.html',
                controller: ['$scope', '$state', '$stateParams', 'serversFactory', 'servers', 'utils',
                    function($scope, $state, $stateParams, serversFactory, servers, utils) {

                        serversFactory.all()
                            .success(function(data) {
                                console.log('successful get in list =' + JSON.stringify(data));
                                $scope.servers = data;
                            })
                            .error(function(error) {
                                $scope.alerts = [{
                                    type: 'danger',
                                    msg: error.message
                                }];
                                console.log('here is an error ' + error.message);
                            });


                        console.log('servers= ' + JSON.stringify($scope.servers));
                        if ($scope.servers.data.length > 0) {
                        	var randId = $scope.servers.data[0].ID;
                            $state.go('servers.detail.details', {
                                serverId: randId
                            });
                        }
                    }
                ]

            })

            .state('servers.monitor', {
                url: '/monitor/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.monitor.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
			    	console.log('at monitor serverId=' + $stateParams.serverId);
                                $scope.item = utils.findById($scope.server.items, $stateParams.itemId);

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('servers.detail.monitor.iostat', {
                url: '/monitor/iostat/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.monitor.iostat.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {

                                $scope.refresh = function() {
                                    serversFactory.iostat($stateParams.serverId)
                                        .success(function(data) {
                                            console.log('successful get with data=' + data);
                                            $scope.iostatresults = data.iostat;
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

            .state('servers.detail.monitor.df', {
                url: '/monitor/df/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.monitor.df.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {

                                $scope.refresh = function() {
                                    console.log('refresh called');
                                    serversFactory.df($stateParams.serverId)
                                        .success(function(data) {
                                            console.log('successful df with data=' + data);
                                            $scope.dfresults = data.df;
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

            .state('servers.detail.monitor.cpu', {
                url: '/monitor/cpu/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.monitor.cpu.html',
                        controller: ['$sce', '$scope', '$stateParams', '$state', 'utils',
                            function($sce, $scope, $stateParams, $state, utils) {
                                console.log('in cpu mon with serverId=' + $stateParams.serverId);
                                console.log('in cpu mon with servername=' + $scope.server.Name);
				$scope.servergraphlink=$sce.trustAsResourceUrl('http://cpm-promdash:3000/embed/seconddashboard#!?var.host=' + $scope.server.Name);


                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('servers.detail.monitor.mem', {
                url: '/monitor/mem/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.monitor.mem.html',
                        controller: ['$sce', '$scope', '$stateParams', '$state', 'utils',
                            function($sce, $scope, $stateParams, $state, utils) {
                                console.log('in mem mon with serverId=' + $stateParams.serverId);
                                console.log('in mem mon with servername=' + $scope.server.Name);
				$scope.servergraphlink=$sce.trustAsResourceUrl('http://cpm-promdash:3000/embed/servermemdashboard#!?var.host=' + $scope.server.Name);

                            }
                        ]
                    },
                }
            })

            .state('servers.detail', {

                url: '/{serverId:[0-9]{1,4}}',

                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'serversFactory', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, serversFactory, utils) {
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in servers');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                if ($scope.servers.data.length > 0) {
                                    angular.forEach($scope.servers.data, function(item) {
                                        if (item.ID == $stateParams.serverId) {
                                            $scope.server = item;
                                        }
                                    });
                                }

                            }
                        ]
                    },

                }
            })

            .state('servers.detail.item', {
                url: '/item/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.item.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.server.items, $stateParams.itemId);

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('servers.detail.add', {
                url: '/add/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.add.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {

                                var newserver = {};
                                newserver.ID = '0';
                                newserver.Name = 'newserver';
                                newserver.IPAddress = '1.1.1.1';
                                newserver.DockerBridgeIP = '172.17.42.1';
                                newserver.PGDataPath = '/var/cpm/data/pgsql';
                                newserver.ServerClass = 'low';
                                newserver.CreateDate = '00';
                                $scope.server = newserver;

                                $scope.add = function() {
                                    console.log('add server is ' + $scope.server.ServerClass);
                                    $scope.server.ID = 0; //0 means to do an insert

                                    serversFactory.add($scope.server)
                                        .success(function(data) {
                                            console.log('successful add with data=' + data);
                                            $state.go('servers.list', $stateParams, {
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
                                    console.log('add called');
                                };
                            }
                        ]
                    },
                }
            })

            .state('servers.detail.details', {
                url: '/details/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.details.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {
                                console.log('server name here is ' + $scope.server.Name);
                                console.log('server bridge ip is ' + $scope.server.DockerBridgeIP);
                                $scope.save = function() {
                                    console.log('saved server is ' + $scope.server.ServerClass);
                                    serversFactory.add($scope.server)
                                        .success(function(data) {
                                            console.log('successful update with data=' + data);
                                            $scope.alerts = [{
                                                type: 'success',
                                                msg: 'saved'
                                            }];
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                    console.log('save called');
                                };

                            }
                        ]
                    },
                }
            })

            .state('servers.detail.delete', {
                url: '/delete/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {

                                var server = $scope.server;
                                $scope.delete = function() {
                                    serversFactory.delete($stateParams.serverId)
                                        .success(function(data) {
                                            console.log('successful delete with data=' + data);
                                            $state.go('servers.list', {});
                                            $state.go('servers.list', $stateParams, {
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
                    },
                }
            })

            .state('servers.detail.containers', {
                url: '/containers/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.containers.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {

                                var server = $scope.server;

                                $scope.getContainers = function() {
                                    serversFactory.containers($stateParams.serverId)
                                        .success(function(data) {
                                            $scope.containers = data;
                                            console.log('successful get containers');
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                            console.log('here is an error ' + error.message);
                                        });
                                };

                                $scope.getContainers();

                            }
                        ]
                    },
                }
            })

            .state('servers.detail.monitor', {
                url: '/monitor/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.monitor.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {}
                        ]
                    },
                }
            })

            .state('servers.detail.containers.start', {
                url: '/start/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.containers.start.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {

                                $scope.start = function() {
                                    console.log('starting server ' + $scope.server.iD);
                                    serversFactory.startall($scope.server.ID)
                                        .success(function(data) {
                                            console.log('successful get in list =' + JSON.stringify(data));
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

            .state('servers.detail.containers.stop', {
                url: '/stop/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.containers.stop.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {

                                $scope.stop = function() {
                                    console.log('stopping server ' + $scope.server.iD);
                                    serversFactory.stopall($scope.server.ID)
                                        .success(function(data) {
                                            console.log('successful get in list =' + JSON.stringify(data));
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
        }
    ]
);
