angular.module('uiRouterSample.settings', [
    'ui.router',
    'ngCookies',
    'ui.bootstrap'
])

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('settings', {

                abstract: true,

                url: '/settings',

                templateUrl: 'app/settings/settings.html',

                resolve: {
                    settings: ['$cookieStore', 'settingsFactory',
                        function($cookieStore, settingsFactory) {
                            if (!$cookieStore.get('cpm_token')) {
                                var nothing = [];
                                console.log('returning nothing');
                                return nothing;
                            }
                            return settingsFactory.all();
                        }
                    ]
                },

                controller: ['$cookieStore', '$scope', '$state', 'settings', 'utils',
                    function($cookieStore, $scope, $state, settings, utils) {

                        console.log('up here with ' + $cookieStore.get('cpm_token'));
                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                userId: 'hi'
                            });
                        }

                        $scope.settings = settings;

                        $scope.goToFirst = function() {
                            console.log($scope.settings.data[0]);
                            var randId = $scope.settings.data[0].Name;

                            $state.go('settings.detail', {
                                settingId: randId
                            });
                        };
                        $scope.goToFirst();
                    }
                ]
            })

            .state('settings.list', {

                url: '',

                templateUrl: 'app/settings/settings.list.html',
                controller: ['$cookieStore', '$scope', '$state', 'settings', 'utils',
                    function($cookieStore, $scope, $state, settings, utils) {

                        console.log('up here with ' + $cookieStore.get('cpm_token'));
                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                userId: 'hi'
                            });
                        }

                        $scope.settings = settings;

                        $scope.goToFirst = function() {
                            console.log($scope.settings.data[0]);
                            var randId = $scope.settings.data[0].Name;

                            $state.go('settings.detail', {
                                settingId: randId
                            });
                        };
                        $scope.goToFirst();
                    }
                ]
            })

            .state('settings.detail', {

                url: '/{settingId}',

                views: {

                    '': {
                        templateUrl: 'app/settings/settings.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'settingsFactory', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, settingsFactory, utils) {
                                if (!$cookieStore.get('cpm_token')) {
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                if ($scope.settings.data.length > 0) {
                                    angular.forEach($scope.settings.data, function(item) {
                                        if (item.Name == $stateParams.settingId) {
                                            console.log('found matching setting');
                                            $scope.setting = item;
                                        }
                                    });
                                }

                                $scope.save = function() {
                                    console.log('saving setting');
                                    settingsFactory.savesetting($scope.setting)
                                        .success(function(data) {
                                            console.log('successful save with data=' + data);
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


            .state('settings.detail.item', {

                url: '/item/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/settings/settings.detail.item.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.setting.items, $stateParams.itemId);

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },

                }
            })

            .state('settings.detail.item.edit', {
                views: {

                    '@settings.detail': {
                        templateUrl: 'app/settings/settings.detail.item.edit.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.setting.items, $stateParams.itemId);
                                $scope.done = function() {
                                    $state.go('^', $stateParams);
                                };
                            }
                        ]
                    }
                }
            });
        }
    ]
);
