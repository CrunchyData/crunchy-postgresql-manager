// Make sure to include the `ui.router` module as a dependency
angular.module('uiRouterSample', [
    'uiRouterSample.projects',
    'uiRouterSample.projects.service',
    'uiRouterSample.containers.service',
    'uiRouterSample.tasks.service',
    'uiRouterSample.clusters.service',
    'uiRouterSample.servers',
    'uiRouterSample.servers.service',
    'uiRouterSample.settings',
    'uiRouterSample.settings.service',
    'uiRouterSample.users',
    'uiRouterSample.users.service',
    'uiRouterSample.rules',
    'uiRouterSample.rules.service',
    'uiRouterSample.home',
    'uiRouterSample.home.service',
    'uiRouterSample.roles',
    'uiRouterSample.roles.service',
    'uiRouterSample.authn',
    'uiRouterSample.authn.service',
    'uiRouterSample.utils.service',
    'angularSpinner',
    'ui.router',
    'ui.bootstrap',
    'ngCookies',
    'ngAnimate',
    'treeControl'
])

.run(
    ['$rootScope', '$cookies', '$cookieStore', '$state', '$stateParams',
        function($rootScope, $cookies, $cookieStore, $state, $stateParams) {

            // It's very handy to add references to $state and $stateParams to the $rootScope
            // so that you can access them from any scope within your applications.For example,
            // <li ng-class="{ active: $state.includes('containers.list') }"> will set the <li>
            // to active whenever 'containers.list' or one of its decendents is active.
            $rootScope.$state = $state;
            $rootScope.$stateParams = $stateParams;
            $rootScope.$cookies = $cookies;
            $rootScope.$cookieStore = $cookieStore;
        }
    ]
)

.config(
    ['$stateProvider', '$urlRouterProvider',
        function($stateProvider, $urlRouterProvider) {

            /////////////////////////////
            // Redirects and Otherwise //
            /////////////////////////////

            // Use $urlRouterProvider to configure any redirects (when) and invalid urls (otherwise).
            $urlRouterProvider

            // The `when` method says if the url is ever the 1st param, then redirect to the 2nd param
            // Here we are just setting up some convenience urls.
                .when('/p?id', '/projects/:id')
                .when('/project/:id', '/projects/:id')
                //.when('/c?id', '/containers/:id/details')
                .when('/container/:id', '/containers/:id')
                .when('/s?id', '/servers/:id')
                .when('/server/:id', '/servers/:id')
                .when('/t?id', '/settings/:id')
                .when('/setting/:id', '/settings/:id')
                .when('/u?id', '/users/:id')
                .when('/user/:id', '/users/:id')
                .when('/x?id', '/profiles/:id')
                .when('/profile/:id', '/profiles/:id')
                .when('/r?id', '/roles/:id')
                .when('/role/:id', '/roles/:id')
                .when('/n?id', '/login/:id')
                .when('/l?id', '/clusters/:id')
                .when('/cluster/:id', '/clusters/:id')

            // If the url is ever invalid, e.g. '/asdf', then redirect to '/' aka the home state
            .otherwise('/');


            //////////////////////////
            // State Configurations //
            //////////////////////////

            // Use $stateProvider to configure your states.
            $stateProvider

            //////////
            // Home //
            //////////

                .state("/", {

                // Use a url of "/" to set a states as the "index".
                url: "/home",

                // Example of an inline template string. By default, templates
                // will populate the ui-view within the parent state's template.
                // For top level states, like this one, the parent template is
                // the index.html file. So this template will be inserted into the
                // ui-view within index.html.
                //template: '<p class="lead">Welcome to Crunchy PostgreSQL Manager</p>' +
                //'<p>Use the menu above to navigate. ' +
                //'<p><a href="#/authn">Login</a> or ' +
                //'<p>Click these links—<a href="#/c?id=1">Alice</a> or ' +
                //'<a href="#/user/42">Bob</a>—to see a url redirect in action.</p>'
                '': {
                    templateUrl: 'app/home/home.html',
                    controller: ['$rootScope', '$scope', '$state', '$stateParams', 'utils',
                        function($rootScope, $scope, $stateParams, utils) {
			    $scope.isCollapsed = true;
			    $scope.projectId = $rootScope.projectId;;
                            console.log('here in app.js controller for home');
                            //$scope.home = utils.findById($scope.home, $stateParams.userId);
                            $state.go('home.detail', { });
                        }
                    ]
                }
            })

            ///////////
            // About //
            ///////////

            .state('about', {
                url: '/about',

                // Showing off how you could return a promise from templateProvider
                templateProvider: ['$timeout',
                    function($timeout) {
                        return $timeout(function() {
                            return '<p class="lead">CPM Resources</p><ul>' +
                                '<li><a href="https://github.com/crunchydata/crunchy-postgresql-manager">Source for CPM</a></li>' +
                                '<li>Version 0.9.5</li>' +
                                '</ul>';
                        }, 100);
                    }
                ]
            })

        }
    ]
);
