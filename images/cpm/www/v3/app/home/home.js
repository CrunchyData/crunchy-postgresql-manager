angular.module('uiRouterSample.home', [
    'ui.router',
    'ngCookies',
    'ui.bootstrap'
])

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {
            $stateProvider
            //////////////
            // Settings //
            //////////////
                .state('home', {

                // With abstract set to true, that means this state can not be explicitly activated.
                // It can only be implicitly activated by activating one of its children.
                abstract: true,

                // This abstract state will prepend '/home' onto the urls of all its children.
                url: '/home',

                // Example of loading a template from a file. This is also a top level state,
                // so this template file will be loaded and then inserted into the ui-view
                // within index.html.
                templateUrl: 'app/home/home.html',

                // You can pair a controller to your template. There *must* be a template to pair with.
                controller: ['$scope', '$state', 'utils',
                    function($scope, $state, utils) {

                        console.log('in home controller here');


                    }
                ]
            })

            ///////////////////////
            // Settings > Detail //
            ///////////////////////

            // You can have unlimited children within a state. Here is a second child
            // state within the 'home' parent state.
            .state('home.detail', {

                url: '',

                views: {

                    // So this one is targeting the unnamed view within the parent state's template.
                    '': {
                        templateUrl: 'app/home/home.detail.html',
                        controller: ['$scope', '$state', '$cookieStore', '$stateParams', 'utils', 'homeFactory',
                            function($scope, $state, $cookieStore, $stateParams, utils, homeFactory) {
                                console.log('here in home 2');
				if (!$cookieStore.get('cpm_token')) {
					$state.go('login', {
						userId: 'hi'
					});
				}

                                homeFactory.healthcheck()
                                    .success(function(data) {
                                        console.log('success with get');
                                        $scope.hcdata = data;
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

        }
    ]
);
