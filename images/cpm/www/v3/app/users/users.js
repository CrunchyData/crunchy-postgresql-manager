angular.module('uiRouterSample.users', [
    'ui.router',
    'ui.bootstrap'
])

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
            //////////////
            // Settings //
            //////////////
                .state('users', {

                // With abstract set to true, that means this state can not be explicitly activated.
                // It can only be implicitly activated by activating one of its children.
                abstract: true,

                // This abstract state will prepend '/users' onto the urls of all its children.
                url: '/users',

                // Example of loading a template from a file. This is also a top level state,
                // so this template file will be loaded and then inserted into the ui-view
                // within index.html.
                templateUrl: 'app/users/users.html',

                // Use `resolve` to resolve any asynchronous controller dependencies
                // *before* the controller is instantiated. In this case, since users
                // returns a promise, the controller will wait until users.all() is
                // resolved before instantiation. Non-promise return values are considered
                // to be resolved immediately.
                resolve: {
                    users: ['$cookieStore', 'usersFactory',
                        function($cookieStore, usersFactory) {
                            if (!$cookieStore.get('cpm_token')) {
                                var nothing = [];
                                console.log('returning nothing');
                                return nothing;
                            }

                            return usersFactory.all();
                        }
                    ]
                },

                // You can pair a controller to your template. There *must* be a template to pair with.
                controller: ['$cookieStore', '$scope', '$state', 'users', 'utils',
                    function($cookieStore, $scope, $state, users, utils) {

                        // Add a 'users' field in this abstract parent's scope, so that all
                        // child state views can access it in their scopes. Please note: scope
                        // inheritance is not due to nesting of states, but rather choosing to
                        // nest the templates of those states. It's normal scope inheritance.
                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                userId: 'hi'
                            });
                        }

                        $scope.users = users;

                        console.log('user is ' + users.data[0].Name);
                        //console.log('roles are ' + JSON.stringify(users.data[0].Roles['foy']));
                        angular.forEach(users.data[0].Roles, function(role) {
                            console.log('role is ' + role.Name);
                        });

                        $scope.goToFirst = function() {
                            var randId = $scope.users.data[0].Name;
                            console.log('in first with randId = ' + randId);

                            // $state.go() can be used as a high level convenience method
                            // for activating a state programmatically.
                            $state.go('users.detail.details', {
                                userId: randId
                            });
                        };
                        $scope.goToFirst();
                    }
                ]
            })

            /////////////////////
            // Settings > List //
            /////////////////////

            // Using a '.' within a state name declares a child within a parent.
            // So you have a new state 'list' within the parent 'users' state.
            .state('users.list', {

                // Using an empty url means that this child state will become active
                // when its parent's url is navigated to. Urls of child states are
                // automatically appended to the urls of their parent. So this state's
                // url is '/users' (because '/users' + '').
                url: '',

                // IMPORTANT: Now we have a state that is not a top level state. Its
                // template will be inserted into the ui-view within this state's
                // parent's template; so the ui-view within users.html. This is the
                // most important thing to remember about templates.
                templateUrl: 'app/users/users.list.html',
                controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'usersFactory',
                    function($scope, $state, $cookieStore, $stateParams, utils, usersFactory) {
                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                userId: 'hi'
                            });
                        }

                        usersFactory.all()
                            .success(function(data) {
                                $scope.users = data;
                            })
                            .error(function(error) {
                                $scope.alerts = [{
                                    type: 'danger',
                                    msg: error.message
                                }];
                                console.log('here is an error ' + error.message);
                            });

                        if ($scope.users.data.length > 0) {
                            console.log('here');
                            var randId = $scope.users.data[0].Name;
                            $state.go('users.detail.details', {
                                userId: randId
                            });
                        }
                    }
                ]
            })

            ///////////////////////
            // Settings > Detail //
            ///////////////////////

            // You can have unlimited children within a state. Here is a second child
            // state within the 'users' parent state.
            .state('users.detail', {

                url: '/{userId}',

                views: {

                    // So this one is targeting the unnamed view within the parent state's template.
                    '': {
                        templateUrl: 'app/users/users.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, utils) {
                                if (!$cookieStore.get('cpm_token')) {
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }
                                if ($scope.users.data.length > 0) {
                                    angular.forEach($scope.users.data, function(item) {
                                        if (item.Name == $stateParams.userId) {
                                            $scope.user = item;
                                            $state.go('users.detail.details', {
                                                userId: item.Name
                                            });
                                        }
                                    });
                                }
                            }
                        ]
                    },

                }
            })

            //////////////////////////////
            // Settings > Detail > Item //
            //////////////////////////////

            .state('users.detail.item', {

                // So following what we've learned, this state's full url will end up being
                // '/users/{userId}/item/:itemId'. We are using both types of parameters
                // in the same url, but they behave identically.
                url: '/item/:itemId',
                views: {

                    // This is targeting the unnamed ui-view within the parent state 'user.detail'
                    // We wouldn't have to do it this way if we didn't also want to set the 'hint' view below.
                    // We could instead just set templateUrl and controller outside of the view obj.
                    '': {
                        templateUrl: 'app/users/users.detail.item.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.user.items, $stateParams.itemId);

                                $scope.edit = function() {
                                    // Here we show off go's ability to navigate to a relative state. Using '^' to go upwards
                                    // and '.' to go down, you can navigate to any relative state (ancestor or descendant).
                                    // Here we are going down to the child state 'edit' (full name of 'users.detail.item.edit')
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },

                    // Here we see we are overriding the template that was set by 'users.detail'
                    'hint@': {
                        template: ' This is users.detail.item overriding the "hint" ui-view'
                    }
                }
            })

            .state('users.detail.details', {
                url: '/details/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/users/users.detail.details.html',
                        controller: ['$scope', '$stateParams', '$cookieStore', '$state', 'usersFactory', 'utils',
                            function($scope, $stateParams, $cookieStore, $state, usersFactory, utils) {
                                console.log('in details with user = ' + $scope.user.Name);
                                console.log('in details with user roles = ' + $scope.user.Roles);
                                $scope.save = function() {
                                    console.log('save called');
                                    $scope.user.Token = $cookieStore.get('cpm_token');
                                    usersFactory.save($scope.user)
                                        .success(function(data) {
                                            console.log('success with save');
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

            //////////////////////////////
            // Settings > Detail > add //
            //////////////////////////////
            .state('users.add', {
                url: '/add/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/users/users.detail.add.html',
                        controller: ['$scope', '$stateParams', '$state', 'usersFactory', 'utils',
                            function($scope, $stateParams, $state, usersFactory, utils) {
                                var newuser = [];
                                newuser.Name = 'newuser';
                                newuser.Roles = $scope.users.data[0].Roles;
                                $scope.user = newuser;

                                $scope.add = function() {
                                    console.log('add called');
                                    usersFactory.add($scope.user)
                                        .success(function(data) {
                                            console.log('success add');
                                            $state.go('users.list', $stateParams, {
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

            .state('users.detail.changepassword', {
                url: '/changepassword/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/users/users.detail.changepassword.html',
                        controller: ['$scope', '$stateParams', '$state', 'usersFactory', 'utils',
                            function($scope, $stateParams, $state, usersFactory, utils) {
                                $scope.changepsw = function() {
                                    console.log('change called');
                                    if ($scope.Password != $scope.ConfirmPassword) {
                                        console.log(' passwords did not match ');
                                        $scope.alerts = [{
                                            type: 'danger',
                                            msg: 'passwords did not match'
                                        }];
                                        return;
                                    }

                                    usersFactory.changepsw($scope.user.Name, $scope.Password)
                                        .success(function(data) {
                                            console.log('success chgpsw');
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

            //////////////////////////////
            // Settings > Detail > delete //
            //////////////////////////////
            .state('users.detail.delete', {
                url: '/delete/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/users/users.detail.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'usersFactory', 'utils',
                            function($scope, $stateParams, $state, usersFactory, utils) {

                                $scope.delete = function() {
                                    console.log('delete called');
                                    usersFactory.delete($scope.user.Name)
                                        .success(function(data) {
                                            console.log('success delete');
                                            $state.go('users.list', $stateParams, {
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

            /////////////////////////////////////
            // Settings > Detail > Item > Edit //
            /////////////////////////////////////

            // Notice that this state has no 'url'. States do not require a url. You can use them
            // simply to organize your application into "places" where each "place" can configure
            // only what it needs. The only way to get to this state is via $state.go (or transitionTo)
            .state('users.detail.item.edit', {
                views: {

                    // This is targeting the unnamed view within the 'users.detail' state
                    // essentially swapping out the template that 'users.detail.item' had
                    // inserted with this state's template.
                    '@users.detail': {
                        templateUrl: 'app/users/users.detail.item.edit.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.user.items, $stateParams.itemId);
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
