angular.module('uiRouterSample.authn', [
    'ui.router',
    'ngCookies'
])

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
            //////////////
            // Settings //
            //////////////
                .state('login', {

                url: '/login',

                views: {
                    '': {
                        templateUrl: 'app/authn/authn.login.html',

                        controller: ['$scope', '$cookieStore', '$state', 'authnFactory', 'utils',
                            function($scope, $cookieStore, $state, authnFactory, utils) {

                                // Add a 'authn' field in this abstract parent's scope, so that all
                                // child state views can access it in their scopes. Please note: scope
                                // inheritance is not due to nesting of states, but rather choosing to
                                // nest the templates of those states. It's normal scope inheritance.
                                //$scope.authn = authn;

                                $scope.alerts = [];

                                $scope.closeAlert = function(index) {
                                    console.log('closeAlert called');
                                    $scope.alerts.splice(index, 1);
                                }
                                $scope.login = function($stateParams) {
                                    console.log('doing login button press');
                                    console.log('user_id=' + $scope.user_id);
                                    console.log('password=' + $scope.password);
                                    console.log('url=' + $scope.admin_url);
                                    authnFactory.doLogin($scope.user_id, $scope.password, $scope.admin_url)
                                        .success(function(results) {
                                            $cookieStore.put('cpm_user_id', $scope.user_id);
                                            $cookieStore.put('AdminURL', $scope.admin_url);
                                            $cookieStore.put('cpm_token', results.Contents);
                                            console.log('token returned was =' + results.Contents);
                                            $state.go('home.detail', {
                                                userId: $scope.user_id
                                            });
                                        })
                                        .error(function(error) {
                                            $scope.alerts = [{
                                                type: 'danger',
                                                msg: error.message
                                            }];
                                        });
                                    //$state.go('authn.detail', { userId: randId });
                                };
                            }
                        ]
                    }
                }
            })

            .state('logout', {

                url: '/logout',

                views: {

                    '': {
                        templateUrl: 'app/authn/authn.logout.html',
                        controller: ['$scope', '$cookieStore', '$state', 'authnFactory', 'utils',
                            function($scope, $cookieStore, $state, authnFactory, utils) {
                                $scope.alerts = [];

                                $scope.closeAlert = function(index) {
                                    console.log('closeAlert called');
                                    $scope.alerts.splice(index, 1);
                                }

                                $scope.logout = function($stateParams) {
                                    console.log('doing logout button press');
                                    console.log('user_id=' + $scope.user_id);
                                    console.log('password=' + $scope.password);
                                    console.log('url=' + $scope.admin_url);

                                    authnFactory.doLogout()
                                        .success(function() {
                                            console.log('here in logout');
                                            $cookieStore.remove('cpm_token');
                                            $scope.alerts = [{
                                                type: 'success',
                                                msg: 'logout successful'
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
                                console.log('here in authn.logout');
                            }
                        ]
                    }
                }
            });

        }
    ]
);
