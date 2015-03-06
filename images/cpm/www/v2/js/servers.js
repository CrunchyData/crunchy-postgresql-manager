// create the module and name it cpmApp
var cpmApp = angular.module('cpmApp.servers', ['ngRoute', 'ui.bootstrap', 'ngTable', 'ngCookies']);

cpmApp.run(function($rootScope) {
    $rootScope.$on('LoadingEvent', function(event, args) {
        $rootScope.$broadcast('LoadingEventTarget', args);
    });

    $rootScope.$on('DoneLoadingEvent', function(event, args) {
        $rootScope.$broadcast('DoneLoadingEventTarget', args);
    });

    $rootScope.$on('updateServerPage', function(event, args) {
        $rootScope.$broadcast('updateServerPageTarget', args);
    });
    $rootScope.$on('noServerEvent', function(event, args) {
        $rootScope.$broadcast('noServerTarget', args);
    });
    $rootScope.$on('reloadServers', function(event, args) {
        $rootScope.$broadcast('reloadServersTarget', args);
    });
    $rootScope.$on('changeServerPage', function(event, args) {
        $rootScope.$broadcast('changeServerPageTarget', args);
    });
    $rootScope.$on('createServerEvent', function(event, args) {
        $rootScope.$broadcast('createServerTarget', args);
    });
    $rootScope.$on('deleteServerEvent', function(event, args) {
        $rootScope.$broadcast('deleteServerTarget', args);
    });
});

var csm = function($rootScope, $scope, $modalInstance, $http, $cookies, $cookieStore) {
    $scope.createServerDialogForm = [];
    $scope.results = [];
    $scope.alerts = [];
    $scope.ServerClass = 'low';

    $scope.cancel = function() {
        console.log('in cancel');
        $modalInstance.dismiss('cancel');
    };

    $scope.ok = function() {
        console.log('in CreateServerModalInstanceCtrl');
        ID = "0";
        cleanIP = this.IPAddress.replace(/\./g, "_");
        cleanBridgeIP = this.DockerBridgeIP.replace(/\./g, "_");
        console.log('cleaned IP=' + cleanIP);
        cleanPath = this.PGDataPath.replace(/\//g, "_");

        var token = $cookieStore.get('cpm_token');

        $http.get($cookieStore.get('AdminURL') + '/addserver/' + ID + "." + this.Name + "." + cleanIP + "." + cleanBridgeIP + "." + cleanPath + "." + this.ServerClass + "." + token).success(function(data, status, headers, config) {
            $scope.results = data;
            $rootScope.$emit('createServerEvent', {
                message: $scope.results
            });
            $modalInstance.close();
        }).error(function(data, status, headers, config) {
            console.log('error in Create Server Modal Instance Ctrl');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

    };
};

var UpdateServerModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {
    $scope.value = value;
    $scope.results = [];
    $scope.currentServer = [];

    $scope.ok = function() {
        console.log('updateServer called');
        console.log('update server got ID=' + $scope.value.ID);
        console.log('update server got name=' + $scope.value.Name);
        console.log('update server got ip=' + $scope.value.IPAddress);
        console.log('update server got dockerip=' + $scope.value.DockerBridgeIP);
        console.log('update server got path=' + $scope.value.PGDataPath);
        cleanIP = $scope.value.IPAddress.replace(/\./g, "_");
        cleanDockerIP = $scope.value.DockerBridgeIP.replace(/\./g, "_");
        console.log('cleaned IP=' + cleanIP);
        cleanPath = $scope.value.PGDataPath.replace(/\//g, "_");

        var token = $cookieStore.get('cpm_token');

        $http.get($cookieStore.get('AdminURL') + '/addserver/' +
            $scope.value.ID + "." +
            $scope.value.Name + "." +
            cleanIP + "." +
            cleanDockerIP + "." + cleanPath + "." + $scope.value.ServerClass + "." + token).success(function(data, status, headers, config) {
            $rootScope.$emit('reloadServers', {
                message: $scope.value
            });
        }).
        error(function(data, status, headers, config) {
            console.log('error:GetServerController update server');
        });
        value.status = 'Completed';
        $modalInstance.close();
    };

    $scope.cancel = function() {
        console.log("current server=" + $scope.value.Name);
        $modalInstance.dismiss('cancel');
    };
};

var DeleteServerModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.value = value;
    $scope.results = [];
    $scope.currentServer = [];

    $scope.ok = function() {
        console.log('in DeleteServerModalInstanceCtrl with ID ' + $scope.value.ID);
        var token = $cookieStore.get('cpm_token');

        $http.get($cookieStore.get('AdminURL') + '/deleteserver/' + $scope.value.ID + "." + token).success(function(data, status, headers, config) {
            $scope.results = data;
            $rootScope.$emit('deleteServerEvent', {
                message: ""
            });
        }).error(function(data, status, headers, config) {
            console.log('error in modal delete server');
        });

        value.status = 'Completed';
        $modalInstance.close();
    };

    $scope.cancel = function() {
        console.log("current server=" + $scope.value.Name);
        $modalInstance.dismiss('cancel');
    };
};



cpmApp.controller('serversController', function($rootScope, $scope, $modal) {
    $scope.status = {
        isopen: false
    };
    $scope.currentServer = [];

    $scope.toggled = function(open) {
        console.log('Dropdown is now: ', open);
    };

    console.log('hi from servers controller');
    $scope.message = 'servers page.';


});

cpmApp.controller('getAllServersController', function($rootScope, $scope, $http, $modal, $cookies, $cookieStore) {
    $scope.tab = 1;

    console.log('getAllServersController');
    $scope.results = [];

    if ($cookieStore.get('AdminURL')) {} else {
        alert('AdminURL setting is NOT defined, please update on the Settings page before using CPM');
    }

    $scope.isSelected = function(checkTab) {
        return $scope.tab === checkTab;
    };
    $scope.selectTab = function(setTab) {
        console.log('setting tab to ' + setTab.ID);
        $scope.tab = setTab.ID;
        $scope.currentServer = setTab;
        $rootScope.$emit('changeServerPage', {
            message: setTab
        });
    }

    $rootScope.$on('deleteServerTarget', function(event, args) {
        console.log('server was deleted....here in getAllServers ' + args.message.Name);
        postit();
    });
    $rootScope.$on('reloadServersTarget', function(event, args) {
        console.log('reloading Servers list' + args.message.Name);
        postit();
    });

    $rootScope.$on('changeServerPage2Target', function(event, args) {
        console.log("cookiestore = [" + $cookieStore.get('cpm_token') + "]");
        console.log('server was updated....here in getAllServers ' + args.message.Name);
        $http.get($cookieStore.get('AdminURL') + '/servers/' + $cookieStore.get('cpm_token')).
        success(function(data, status, headers, config) {
            $scope.results = data;
        }).error(function(data, status, headers, config) {
            console.log('error:GetAllServersController.http.get');
        });
    });

    $rootScope.$on('createServerTarget', function(event, args) {
        console.log('server was created....here in getAllServers ' + args.message.Name);
        postit();
        for (i = 0; i < $scope.results.length; i++) {
            if ($scope.results[i].Name == args.message.Name) {
                $scope.results[i].active = true;
            }
        }
        $rootScope.$emit('updateServerPage', {
            message: args.message
        });
    });

    var postit = function() {
        $http.get($cookieStore.get('AdminURL') + '/servers/' + $cookieStore.get('cpm_token')).
        success(function(data, status, headers, config) {
            $scope.results = data;
            if (data.length > 0) {
                console.log('setting tab to ' + data[0].ID);
		$scope.selectTab(data[0]);
            } else {
        	$rootScope.$emit('noServerEvent', {
            	message: 'hi'
        	});
	    }
        }).error(function(data, status, headers, config) {
            console.log('error:GetAllServersController.http.get');
        });
    };

    postit();

});


cpmApp.controller('getServerController', function($scope, $http, $rootScope, $q, $filter, $modal, $cookies, $cookieStore, ngTableParams) {
    console.log('getServerController');

    $scope.currentServer = [];
    $scope.currentServerID = [];
    $scope.results = [];
    $scope.containers = [];
    $scope.selectedContainers = [];
    $scope.users = [];
    $scope.data = [];

    $scope.tableParams = new ngTableParams({
        page: 1, // show first page
        count: 10 // count per page
    }, {
        total: $scope.containers.length, // length of data
        getData: function($defer, params) {
            console.log('getData called containers=' + $scope.containers.length);
            // use build-in angular filter
            var orderedData = $scope.containers;

            params.total(orderedData.length); // set total for recalc pagination
            $defer.resolve($scope.users = orderedData.slice((params.page() - 1) * params.count(), params.page() * params.count()));
        }
    });

    //fix around ng-table bug?
    $scope.tableParams.settings().$scope = $scope;

    $scope.checkboxes = {
        'checked': false,
        items: {}
    };

    // watch for check all checkbox
    $scope.$watch('checkboxes.checked', function(value) {
        angular.forEach($scope.users, function(item) {
            if (angular.isDefined(item.ID)) {
                $scope.checkboxes.items[item.ID] = value;
            }
        });
    });

    // watch for data checkboxes
    $scope.$watch('checkboxes.items', function(values) {
        if (!$scope.users) {
            return;
        }
        var checked = 0,
            unchecked = 0,
            total = $scope.users.length;
        angular.forEach($scope.users, function(item) {
            checked += ($scope.checkboxes.items[item.ID]) || 0;
            unchecked += (!$scope.checkboxes.items[item.ID]) || 0;
        });
        if ((unchecked == 0) || (checked == 0)) {
            $scope.checkboxes.checked = (checked == total);
        }
        // grayed checkbox
        angular.element(document.getElementById("select_all")).prop("indeterminate", (checked != 0 && unchecked != 0));
    }, true);

    function postit(serverid) {
        var token = $cookieStore.get('cpm_token');
        console.log('in GetServerController id=' + serverid);
        $http.get($cookieStore.get('AdminURL') + '/server/' + serverid + "." + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
            $rootScope.$broadcast('CurrentServer', data);
        }).
        error(function(data, status, headers, config) {
            console.log('error GetServerController http.get');
        });

        $http.get($cookieStore.get('AdminURL') + '/nodes/forserver/' + serverid + "." + token).
        success(function(data, status, headers, config) {
            $scope.containers = data;
            $scope.tableParams.reload();
            console.log('containers has ' + $scope.containers.length);
        }).
        error(function(data, status, headers, config) {
            console.log('error: GetServerController.http.get 2');
        });
    }

    $scope.stopContainers = function() {
	$scope.status.isopen = false;
        console.log('stopContainers called');
        var names = '';
        var token = $cookieStore.get('cpm_token');
        angular.forEach($scope.users, function(item) {
            if (angular.isDefined(item.ID)) {
                if ($scope.checkboxes.items[item.ID]) {
                    names += ' ' + item.ID;
                    $rootScope.$emit('LoadingEvent', {
                        message: ""
                    });
                    $http.get($cookieStore.get('AdminURL') + '/admin/stop/' + item.ID + '.' + token).success(function(data, status, headers, config) {
                        console.log('stop container success id=' + item.ID);
                        for (index = 0; index < $scope.containers.length; index++) {
                            if ($scope.containers[index].ID == item.ID) {
                                $scope.containers[index].Status = 'down';
                                $scope.checkboxes.items[item.ID] = false;
                                console.log('setting stopped container to down status');
                            }
                        }
                        $rootScope.$emit('DoneLoadingEvent', {
                            message: ""
                        });
                    }).
                    error(function(data, status, headers, config) {
                        console.log('error:stop container');
                    });
                }
            }
        });
        console.log('stop these:' + names);
    }

    $scope.startContainers = function() {
	$scope.status.isopen = false;
        $rootScope.$emit('LoadingEvent', {
            message: ""
        });
        console.log('startContainers called');
        var token = $cookieStore.get('cpm_token');
        var names = '';
        angular.forEach($scope.users, function(item) {
            if (angular.isDefined(item.ID)) {
                if ($scope.checkboxes.items[item.ID]) {
                    names += ' ' + item.ID;
                    $http.get($cookieStore.get('AdminURL') + '/admin/start/' + item.ID + '.' + token).success(function(data, status, headers, config) {
                        console.log('start container success id=' + item.ID);
                        for (index = 0; index < $scope.containers.length; index++) {
                            if ($scope.containers[index].ID == item.ID) {
                                $scope.containers[index].Status = 'up';
                                $scope.checkboxes.items[item.ID] = false;
                                console.log('setting started container to up status');
                            }
                        }
                    }).
                    error(function(data, status, headers, config) {
                        console.log('error:start container');
                    });
                }
            }
        });
        console.log('start these:' + names);
        $rootScope.$emit('DoneLoadingEvent', {
            message: ""
        });
    }

    $scope.updateServer = function() {
	$scope.status.isopen = false;
        var modalInstance = $modal.open({
            templateUrl: 'pages/updateservermodal.html',
            controller: UpdateServerModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.results;
                }
            }
        });
    }

    $scope.handleAddClick = function(msg) {
        console.log('handleAddClick');
	$scope.status.isopen = false;
        var modalInstance = $modal.open({
            templateUrl: 'pages/createservermodal.html',
            controller: csm
        });
        modalInstance.result.then(function(response) {
            console.log('jeff in handlePlusClick with response =' + response);
            console.log('jeff in handlePlusClick with Name returned=' + this.name);
        });
    };

    $scope.handleDeleteClick = function(msg) {
        console.log('handleDeleteClick');
	$scope.status.isopen = false;
        var modalInstance = $modal.open({
            templateUrl: 'pages/deleteservermodal.html',
            controller: DeleteServerModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.currentServer;
                }
            }
        });
    };

    $scope.handleMonitorClick = function(msg) {
	$scope.status.isopen = false;
        console.log('hi from handleAdminClick id=' + $scope.currentServer.ID);
        var popupWindow = window.open('pages/servermonitor.html');
        console.log('in app.js setting child serverid=' + $scope.currentServer.ID);
        popupWindow.serverid = $scope.currentServer.ID;
    };

    $rootScope.$on('changeServerPageTarget', function(event, args) {
        console.log('here in GetServer ' + args.message.Name);
        $scope.currentServer = args.message;
        $scope.currentServerID = args.message.ID;
        $scope.message = args.message.ID;
        $scope.entryID = $scope.message.ID;
        postit(args.message.ID);
    });

    $rootScope.$on('noServerTarget', function(event, args) {
        console.log('no server event received');
        $scope.currentServer = [];
    	$scope.results = [];
    	$scope.containers = [];
    });


});
