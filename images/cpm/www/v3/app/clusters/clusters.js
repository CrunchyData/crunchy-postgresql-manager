angular.module('uiRouterSample.clusters', [
    'ui.router',
    'ui.bootstrap'
])

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('clusters', {

                abstract: true,

                url: '/clusters',

                templateUrl: 'app/clusters/clusters.html',

                resolve: {
                    clusters: ['$cookieStore', 'clustersFactory',
                        function($cookieStore, clustersFactory) {
                            if (!$cookieStore.get('cpm_token')) {
                                var nothing = [];
                                console.log('returning nothing');
                                return nothing;
                            }

                            return clustersFactory.all();
                        }
                    ]
                },

                controller: ['$cookieStore', '$scope', '$state', 'clusters', 'utils',
                    function($cookieStore, $scope, $state, clusters, utils) {

                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                userId: 'hi'
                            });
                        }

                        $scope.clusters = clusters;

                        $scope.goToFirst = function() {
                            if ($scope.clusters.length > 0) {
                                var randId = $scope.clusters.data[0].ID;

                                $state.go('clusters.detail.details', {
                                    clusterId: randId
                                });
                            }
                        };
                        $scope.goToFirst();
                    }
                ]
            })

            .state('clusters.list', {

                url: '',

                templateUrl: 'app/clusters/clusters.list.html',
                controller: ['$scope', '$state', 'clusters', 'utils', 'clustersFactory',
                    function($scope, $state, clusters, utils, clustersFactory) {

                        $scope.clusters = clusters;

                        clustersFactory.all()
                            .success(function(data) {
                                console.log('successful clusters all with data=' + data);
                                $scope.clusters = data;
                        	$scope.goToFirst();

                            })
                            .error(function(error) {
                                $scope.alerts = [{
                                    type: 'danger',
                                    msg: error.message
                                }];
                                console.log('here is an error ' + error.message);
                            });

                        $scope.goToFirst = function() {
                            if ($scope.clusters.length > 0) {
                                var randId = $scope.clusters[0].ID;

                                $state.go('clusters.detail.details', {
                                    clusterId: randId
                                });
                            }
                        };
                    }
                ]


            })

            .state('clusters.detail', {

                url: '/{clusterId:[0-9]{1,4}}',

                views: {

                    '': {
                        templateUrl: 'app/clusters/clusters.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, utils) {
                                if (!$cookieStore.get('cpm_token')) {
                                    console.log('cpm_token not defined in projects');
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                if ($scope.clusters.data.length > 0) {
                                    angular.forEach($scope.clusters.data, function(item) {
                                        if (item.ID == $stateParams.clusterId) {
                                            $scope.cluster = item;
                                            console.log(JSON.stringify(item));
                                        }
                                    });
                                }

                            }
                        ]
                    },

                }
            })


            .state('clusters.autocluster', {

                url: '/autocluster',
                views: {

                    '': {
                        templateUrl: 'app/clusters/clusters.autocluster.html',
                        controller: ['$scope', '$stateParams', '$state', 'clustersFactory', 'utils', 'usSpinnerService',
                            function($scope, $stateParams, $state, clustersFactory, utils, usSpinnerService) {

                                $scope.ClusterProfile = 'SM';
                                $scope.cluster = [];
                                $scope.cluster.Name = 'cluster1';
                                $scope.cluster.ClusterType = 'asynchronous';

                                console.log('ClusterProfile=' + $scope.ClusterProfile);
                                console.log('cluster.Name=' + $scope.cluster.Name);

                                $scope.create = function() {
                                    usSpinnerService.spin('spinner-1');
                                    console.log('ClusterProfile=' + $scope.ClusterProfile);
                                    console.log('cluster.Name=' + $scope.cluster.Name);
                                    console.log('cluster.ClusterType=' + $scope.cluster.ClusterType);
                                    clustersFactory.autocluster($scope.cluster, $scope.ClusterProfile)
                                        .success(function(data) {
                                            console.log('successful autocreate with data=' + data);
                                            $scope.results = data;
                                            usSpinnerService.stop('spinner-1');
                                            $state.go('clusters.list', $stateParams, {
                                                reload: true,
                                                inherit: false
                                            });
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
                                            usSpinnerService.stop('spinner-1');
                                        });
                                };

                            }
                        ]
                    },
                }
            })

            .state('clusters.detail.details', {
                url: '/details/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/clusters/clusters.detail.details.html',
                        controller: ['$scope', '$stateParams', '$state', 'clustersFactory', 'utils',
                            function($scope, $stateParams, $state, clustersFactory, utils) {
                                console.log('cluster name here is ' + $scope.cluster.Name);
                                $scope.save = function() {
                                    console.log('saved cluster is ' + $scope.cluster.Name);
                                    serversFactory.add($scope.cluster)
                                        .success(function(data) {
                                            console.log('successful update with data=' + data);
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

            .state('clusters.detail.delete', {

                url: '/delete/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/clusters/clusters.detail.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'clustersFactory', 'utils', 'usSpinnerService',
                            function($scope, $stateParams, $state, clustersFactory, utils, usSpinnerService) {
                                var cluster = $scope.cluster;

                                $scope.delete = function() {
                                    usSpinnerService.spin('spinner-1');
                                    console.log('delete cluster called');
                                    clustersFactory.delete($stateParams.clusterId)
                                        .success(function(data) {
                                            console.log('successful delete with data=' + data);
                                            usSpinnerService.stop('spinner-1');
                                            $state.go('clusters.list', $stateParams, {
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

            .state('clusters.detail.scale', {

                url: '/scale/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/clusters/clusters.detail.scale.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('clusters.detail.define', {

                url: '/define/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/clusters/clusters.detail.define.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('clusters.detail.add', {
                url: '/add/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/clusters/clusters.detail.add.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.cluster.items, $stateParams.itemId);

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },
                }
            })

            .state('clusters.detail.container', {

                    url: '/container/:itemId',
                    views: {

                        '': {
                            templateUrl: 'app/clusters/clusters.detail.item.html',
                            controller: ['$scope', '$stateParams', '$state', 'utils',
                                function($scope, $stateParams, $state, utils) {
                                    $stateParams.containerId = $stateParams.itemId;
                                    $state.go('containers.detail', $stateParams);
                                }
                            ]
                        },

                    }
                })
        }
    ]
);
