angular.module('uiRouterSample.rules', [
    'ui.router',
    'ui.bootstrap'
])

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('rules', {

                abstract: true,

                url: '/rules',

                templateUrl: 'app/rules/rules.html',

                resolve: {
                    rules: ['$cookieStore', 'rulesFactory',
                        function($cookieStore, rulesFactory) {
                            if (!$cookieStore.get('cpm_token')) {
                                var nothing = [];
                                console.log('returning nothing');
                                return nothing;
                            }

                            return rulesFactory.all();
                        }
                    ]
                },

                controller: ['$cookieStore', '$scope', '$state', 'rules', 'utils',
                    function($cookieStore, $scope, $state, rules, utils) {

                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                ruleId: 'hi'
                            });
                        }

                        $scope.rules = rules;

                        console.log('rule is ' + rules.data[0].Name);
                        //console.log('roles are ' + JSON.stringify(rules.data[0].Roles['foy']));
                        $scope.goToFirst = function() {
                            var randId = $scope.rules.data[0].Name;
                            console.log('in first with randId = ' + randId);

                            // $state.go() can be used as a high level convenience method
                            // for activating a state programmatically.
                            $state.go('rules.detail.details', {
                                ruleId: randId
                            });
                        };
                        $scope.goToFirst();
                    }
                ]
            })

            .state('rules.list', {

                url: '',

                templateUrl: 'app/rules/rules.list.html',
                controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'rulesFactory',
                    function($scope, $state, $cookieStore, $stateParams, utils, rulesFactory) {
                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                ruleId: 'hi'
                            });
                        }

                        rulesFactory.all()
                            .success(function(data) {
                                $scope.rules = data;
                            })
                            .error(function(error) {
                                $scope.alerts = [{
                                    type: 'danger',
                                    msg: error.message
                                }];
                                console.log('here is an error ' + error.message);
                            });

                        if ($scope.rules.data.length > 0) {
                            console.log('here');
                            var randId = $scope.rules.data[0].Name;
                            $state.go('rules.detail.details', {
                                ruleId: randId
                            });
                        }
                    }
                ]
            })

            .state('rules.detail', {

                url: '/{ruleId}',

                views: {

                    '': {
                        templateUrl: 'app/rules/rules.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, utils) {
                                if (!$cookieStore.get('cpm_token')) {
                                    $state.go('login', {
                                        ruleId: 'hi'
                                    });
                                }
                                if ($scope.rules.data.length > 0) {
                                    angular.forEach($scope.rules.data, function(item) {
                                        if (item.Name == $stateParams.ruleId) {
                                            $scope.rule = item;
                                            $state.go('rules.detail.details', {
                                                ruleId: item.Name
                                            });
                                        }
                                    });
                                }
                            }
                        ]
                    },

                }
            })

            .state('rules.detail.item', {

                url: '/item/:itemId',
                views: {

                    '': {
                        templateUrl: 'app/rules/rules.detail.item.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.rule.items, $stateParams.itemId);

                                $scope.edit = function() {
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },

                    'hint@': {
                        template: ' This is rules.detail.item overriding the "hint" ui-view'
                    }
                }
            })

            .state('rules.detail.details', {
                url: '/details/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/rules/rules.detail.details.html',
                        controller: ['$scope', '$stateParams', '$cookieStore', '$state', 'rulesFactory', 'utils', 'usSpinnerService',
                            function($scope, $stateParams, $cookieStore, $state, rulesFactory, utils, usSpinnerService) {
                                console.log('in details with rule itemId = ' + $stateParams.itemId);
                                //console.log('in details with rule = ' + $scope.rule.Name);
                                $scope.save = function() {
                                    console.log('save called');
                                    $scope.rule.Token = $cookieStore.get('cpm_token');
				     usSpinnerService.spin('spinner-1'); 
                                    rulesFactory.update($scope.rule)
                                        .success(function(data) {
                                            console.log('success with save');
                                            $scope.alerts = [{
                                                type: 'success',
                                                msg: 'successfully updated'
                                            }];
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
                    }
                }
            })

            .state('rules.add', {
                url: '/add/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/rules/rules.detail.add.html',
                        controller: ['$scope', '$stateParams', '$state', 'rulesFactory', 'utils',
                            function($scope, $stateParams, $state, rulesFactory, utils) {
                                var newrule = [];
                                newrule.Name = 'newrule';
                                newrule.Method = 'md5';
                                newrule.Type = 'host';
                                $scope.rule = newrule;

                                $scope.add = function() {
                                    console.log('add called');
                                    rulesFactory.add($scope.rule)
                                        .success(function(data) {
                                            console.log('success add');
                                            $state.go('rules.list', $stateParams, {
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

            .state('rules.detail.delete', {
                url: '/delete/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/rules/rules.detail.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'rulesFactory', 'utils',
                            function($scope, $stateParams, $state, rulesFactory, utils) {

                                $scope.delete = function() {
                                    console.log('delete called');
                                    rulesFactory.delete($scope.rule.ID)
                                        .success(function(data) {
                                            console.log('success delete');
                                            $state.go('rules.list', $stateParams, {
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

            .state('rules.detail.item.edit', {
                views: {

                    '@rules.detail': {
                        templateUrl: 'app/rules/rules.detail.item.edit.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.rule.items, $stateParams.itemId);
                                $scope.done = function() {
                                    // Go back up. '^' means up one. '^.^' would be up twice, to the grandparent.
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
