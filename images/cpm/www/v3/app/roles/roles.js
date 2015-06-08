angular.module('uiRouterSample.roles', [
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
                .state('roles', {

                abstract: true,

                url: '/roles',

                templateUrl: 'app/roles/roles.html',

                resolve: {
                    roles: ['$cookieStore', 'rolesFactory',
                        function($cookieStore, rolesFactory) {
                            if (!$cookieStore.get('cpm_token')) {
                                var nothing = [];
                                console.log('returning nothing');
                                return nothing;
                            }

                            return rolesFactory.all();
                        }
                    ]
                },

                controller: ['$cookieStore', '$scope', '$state', 'roles', 'utils',
                    function($cookieStore, $scope, $state, roles, utils) {

                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                userId: 'hi'
                            });
                        }


                        console.log('in roles');
                        $scope.roles = roles;

                        $scope.goToFirst = function() {
                            var randId = $scope.roles.data[0].Name;

                            $state.go('roles.detail.details', {
                                roleId: randId
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
            // So you have a new state 'list' within the parent 'roles' state.
            .state('roles.list', {

                // Using an empty url means that this child state will become active
                // when its parent's url is navigated to. Urls of child states are
                // automatically appended to the urls of their parent. So this state's
                // url is '/roles' (because '/roles' + '').
                url: '',

                // IMPORTANT: Now we have a state that is not a top level state. Its
                // template will be inserted into the ui-view within this state's
                // parent's template; so the ui-view within roles.html. This is the
                // most important thing to remember about templates.
                templateUrl: 'app/roles/roles.list.html',
                controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'rolesFactory',
                    function($scope, $state, $cookieStore, $stateParams, utils, rolesFactory) {
                        if (!$cookieStore.get('cpm_token')) {
                            $state.go('login', {
                                userId: 'hi'
                            });
                        }

                        console.log('in roles.list');
                        rolesFactory.all()
                            .success(function(data) {
                                $scope.roles = data;
                            })
                            .error(function(error) {
                                $scope.alerts = [{
                                    type: 'danger',
                                    msg: error.message
                                }];
                                console.log('here is an error ' + error.message);
                            });

                        if ($scope.roles.data.length > 0) {
                            console.log('here');
                            var randId = $scope.roles.data[0].Name;
                            $state.go('roles.detail.details', {
                                roleId: randId
                            });
                        }
                    }
                ]

            })

            ///////////////////////
            // Settings > Detail //
            ///////////////////////

            // You can have unlimited children within a state. Here is a second child
            // state within the 'roles' parent state.
            .state('roles.detail', {

                // Urls can have parameters. They can be specified like :param or {param}.
                // If {} is used, then you can also specify a regex pattern that the param
                // must match. The regex is written after a colon (:). Note: Don't use capture
                // groups in your regex patterns, because the whole regex is wrapped again
                // behind the scenes. Our pattern below will only match numbers with a length
                // between 1 and 4.

                // Since this state is also a child of 'roles' its url is appended as well.
                // So its url will end up being '/roles/{roleId:[0-9]{1,4}}'. When the
                // url becomes something like '/roles/42' then this state becomes active
                // and the $stateParams object becomes { roleId: 42 }.
                url: '/{roleId}',

                // If there is more than a single ui-view in the parent template, or you would
                // like to target a ui-view from even higher up the state tree, you can use the
                // views object to configure multiple views. Each view can get its own template,
                // controller, and resolve data.

                // View names can be relative or absolute. Relative view names do not use an '@'
                // symbol. They always refer to views within this state's parent template.
                // Absolute view names use a '@' symbol to distinguish the view and the state.
                // So 'foo@bar' means the ui-view named 'foo' within the 'bar' state's template.
                views: {

                    // So this one is targeting the unnamed view within the parent state's template.
                    '': {
                        templateUrl: 'app/roles/roles.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils',
                            function($scope, $state, $cookieStore, $stateParams, utils) {
                                if (!$cookieStore.get('cpm_token')) {
                                    $state.go('login', {
                                        userId: 'hi'
                                    });
                                }

                                console.log('here in roles.detail');
                                if ($scope.roles.data.length > 0) {
                                    angular.forEach($scope.roles.data, function(item) {
                                        if (item.Name == $stateParams.roleId) {
                                            console.log('matched ' + item.Name);
                                            $scope.role = item;
                                            $state.go('roles.detail.details', {
                                                roleId: item.Name
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

            .state('roles.detail.item', {

                // So following what we've learned, this state's full url will end up being
                // '/roles/{roleId}/item/:itemId'. We are using both types of parameters
                // in the same url, but they behave identically.
                url: '/item/:itemId',
                views: {

                    // This is targeting the unnamed ui-view within the parent state 'role.detail'
                    // We wouldn't have to do it this way if we didn't also want to set the 'hint' view below.
                    // We could instead just set templateUrl and controller outside of the view obj.
                    '': {
                        templateUrl: 'app/roles/roles.detail.item.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                $scope.item = utils.findById($scope.role.items, $stateParams.itemId);

                                console.log('in roles.detail.item');
                                $scope.edit = function() {
                                    // Here we show off go's ability to navigate to a relative state. Using '^' to go upwards
                                    // and '.' to go down, you can navigate to any relative state (ancestor or descendant).
                                    // Here we are going down to the child state 'edit' (full name of 'roles.detail.item.edit')
                                    $state.go('.edit', $stateParams);
                                };
                            }
                        ]
                    },

                    // Here we see we are overriding the template that was set by 'roles.detail'
                    //'hint@': {
                    //template: ' This is roles.detail.item overriding the "hint" ui-view'
                    //}
                }
            })

            .state('roles.detail.details', {
                url: '/details/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/roles/roles.detail.details.html',
                        controller: ['$scope', '$stateParams', '$state', 'rolesFactory', 'utils',
                            function($scope, $stateParams, $state, rolesFactory, utils) {
                                console.log('in detail.details');
                                $scope.save = function() {
                                    console.log('save called');
                                    console.log('perms=' + $scope.role.Permissions);
                                    rolesFactory.save($scope.role)
                                        .success(function(data) {
                                            console.log('successful save with data=' + data);
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

            .state('roles.add', {
                url: '/add/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/roles/roles.detail.add.html',
                        controller: ['$scope', '$stateParams', '$state', 'rolesFactory', 'utils',
                            function($scope, $stateParams, $state, rolesFactory, utils) {
                                $scope.newrole = [];
                                $scope.newrole.Name = 'newrole';
                                //newrole.Permissions = $scope.roles.data[0].Permissions;
                                $scope.newrole.Permissions = $scope.roles.data[0].Permissions;
                                /**
		      'perm-backup': {
			'Name': 'perm-backup',
			'Description': 'perform backups',
			'Selected': false
			},
			'perm-cluster': {
			'Name': 'perm-cluster',
			'Description': 'maintain clusters',
			'Selected': false
			},
			'perm-container': {
			'Name': 'perm-container',
			'Description': 'maintain containers',
			'Selected': false
			},
			'perm-server': {
			'Name': 'perm-server',
			'Description': 'maintain servers',
			'Selected': false
			},
			'perm-setting': {
			'Name': 'perm-setting',
			'Description': 'maintain settings',
			'Selected': false
			},
			'perm-user': {
			'Name': 'perm-user',
			'Description': 'maintain users',
			'Selected': false
			}
		};
		*/

                                $scope.add = function() {
                                    console.log('add called role=' + $scope.newrole.Name);
                                    console.log('add called role perms=' + $scope.newrole.Permissions);

                                    rolesFactory.add($scope.newrole)
                                        .success(function(data) {
                                            console.log('successful add with data=' + data);
                                            $state.go('roles.list', $stateParams, {
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

            .state('roles.detail.add', {
                url: '/add/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/roles/roles.detail.add.html',
                        controller: ['$scope', '$stateParams', '$state', 'rolesFactory', 'utils',
                            function($scope, $stateParams, $state, rolesFactory, utils) {
                                $scope.add = function() {
                                    console.log('add called role=' + $scope.role.Name);

                                    rolesFactory.add($scope.role)
                                        .success(function(data) {
                                            console.log('successful add with data=' + data);
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

            //////////////////////////////
            // Settings > Detail > delete //
            //////////////////////////////
            .state('roles.detail.delete', {
                url: '/delete/:itemId',
                views: {
                    '': {
                        templateUrl: 'app/roles/roles.detail.delete.html',
                        controller: ['$scope', '$stateParams', '$state', 'rolesFactory', 'utils',
                            function($scope, $stateParams, $state, rolesFactory, utils) {
                                $scope.delete = function() {
                                    rolesFactory.delete($scope.role.Name)
                                        .success(function(data) {
                                            console.log('successful delete with data=' + data);
                                            $state.go('roles.list', $stateParams, {
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
                                    console.log('delete called');
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
            .state('roles.detail.item.edit', {
                views: {

                    // This is targeting the unnamed view within the 'roles.detail' state
                    // essentially swapping out the template that 'roles.detail.item' had
                    // inserted with this state's template.
                    '@roles.detail': {
                        templateUrl: 'app/roles/roles.detail.item.edit.html',
                        controller: ['$scope', '$stateParams', '$state', 'utils',
                            function($scope, $stateParams, $state, utils) {
                                console.log('in here too');
                                $scope.item = utils.findById($scope.role.items, $stateParams.itemId);
                                $scope.done = function() {
                                    // Go back up. '^' means up one. '^.^' would be up twice, to the grandparent.
                                    //$state.go('^', $stateParams);
                                };
                            }
                        ]
                    }
                }
            });
        }
    ]
);
