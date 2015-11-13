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

                    controller: ['$scope', '$state', '$stateParams', '$cookieStore', 'utils', 'servers',
                        function($scope, $state, $stateParams, $cookieStore, utils, servers) {

                            if (!$cookieStore.get('cpm_token')) {
                                $state.go('login', {
                                    userId: 'hi'
                                });
                            }

                            $scope.servers = servers;

                            $scope.goToFirst = function() {
                                if ($scope.servers.data.length > 0) {
                                    //console.log($scope.servers.data[0].ID);
                                    var randId = $scope.servers.data[0].ID;

                                    $state.go('servers.detail.details', {
                                        serverId: randId
                                    });
                                }
                            };

				console.log('before goToFirst with serverId=' + $stateParams.serverId);
				if ($stateParams.serverId === undefined) {
					console.log('serverId is undefined here jeff');
				} else {
					console.log('serverId is defined here jeff');
				}
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
                                //console.log('successful get in list =' + JSON.stringify(data));
                                $scope.servers = data;
                            })
                            .error(function(error) {
                                $scope.alerts = [{
                                    type: 'danger',
                                    msg: error.message
                                }];
                                console.log('here is an error ' + error.message);
                            });


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
                url: '/iostat',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.monitor.iostat.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {

                                $scope.refresh = function() {
                                    serversFactory.iostat($stateParams.serverId)
                                        .success(function(data) {
                                            $scope.iostatresults = data.iostat;
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
                        ]
                    },
                }
            })

            .state('servers.detail.monitor.df', {
                url: '/df',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.monitor.df.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {

                                $scope.refresh = function() {
                                    serversFactory.df($stateParams.serverId)
                                        .success(function(data) {
                                            $scope.dfresults = data.df;
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
                        ]
                    },
                }
            })

            .state('servers.detail.monitor.cpu', {
                url: '/cpu',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.monitor.cpu.html',
                        controller: ['$sce', '$scope', '$stateParams', '$state', 'utils',
                            function($sce, $scope, $stateParams, $state, utils) {
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
                url: '/mem',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.monitor.mem.html',
                        controller: ['$sce', '$scope', '$stateParams', '$state', 'utils',
                            function($sce, $scope, $stateParams, $state, utils) {
				$scope.servergraphlink=$sce.trustAsResourceUrl('http://cpm-promdash:3000/embed/servermemdashboard#!?var.host=' + $scope.server.Name);

                            }
                        ]
                    },
                }
            })

            .state('servers.detail', {

                url: '/:serverId',

                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'serversFactory', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, serversFactory, utils) {
			    	console.log('in servers.detail with stateParams = ' + JSON.stringify($stateParams));
			    	console.log('in servers.detail with serverId = ' + $stateParams.serverId);
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in servers');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

				//$scope.server = serversFactory.get($stateParams.serverId);
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
                                newserver.ServerClass = 'low';
                                newserver.CreateDate = '00';
                                $scope.server = newserver;

                                $scope.add = function() {
                                    $scope.server.ID = 0; //0 means to do an insert

                                    serversFactory.add($scope.server)
                                        .success(function(data) {
                                            $state.go('servers.list', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.Error
                                            }];
                                        });
                                };
                            }
                        ]
                    },
                }
            })

            .state('servers.detail.details', {
                url: '/details',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.details.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils',
                            function($scope, $stateParams, $state, serversFactory, utils) {
                                $scope.save = function() {
                                    serversFactory.add($scope.server)
                                        .success(function(data) {
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
                url: '/containers',
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
                                            //console.log('successful get in list =' + JSON.stringify(data));
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
                url: '/monitor',
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
                url: '/start',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.containers.start.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils', 'spinnerService',
                            function($scope, $stateParams, $state, serversFactory, utils, spinnerService) {

                                $scope.start = function() {
					spinnerService.show('startallspinner');
                                    serversFactory.startall($scope.server.ID)
                                        .success(function(data) {
                                            //console.log('successful get in list =' + JSON.stringify(data));
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
                                        })
					.finally(function() {
						spinnerService.hide('startallspinner');
					});

                                };
                            }
                        ]
                    },
                }
            })

            .state('servers.detail.containers.stop', {
                url: '/stop',
                views: {

                    '': {
                        templateUrl: 'app/servers/servers.detail.containers.stop.html',
                        controller: ['$scope', '$stateParams', '$state', 'serversFactory', 'utils', 'spinnerService',
                            function($scope, $stateParams, $state, serversFactory, utils, spinnerService) {

                                $scope.stop = function() {
					spinnerService.show('stopallspinner');
                                    serversFactory.stopall($scope.server.ID)
                                        .success(function(data) {
                                            //console.log('successful get in list =' + JSON.stringify(data));
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
                                        })
					.finally(function() {
						spinnerService.hide('stopallspinner');
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
