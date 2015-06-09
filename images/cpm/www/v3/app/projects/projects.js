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


                        $scope.projects = projects;

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
