// create the module and name it cpmApp
var cpmApp = angular.module('cpmApp.containers', ['ngRoute', 'ngCookies']);


cpmApp.controller('containersController', function($scope, $rootScope, $cookies, $routeParams ) {
    console.log('hi from containers controller');
    //var param1 = $routeParams.param1;

    $scope.message = 'containers page.';
});


cpmApp.run(function($rootScope) {
    $rootScope.$on('LoadingEvent', function(event, args) {
        $rootScope.$broadcast('LoadingEventTarget', args);
        $rootScope.isLoading = true;
    });
    $rootScope.$on('DoneLoadingEvent', function(event, args) {
        $rootScope.$broadcast('DoneLoadingEventTarget', args);
        $rootScope.isLoading = false;
    });
    $rootScope.$on('deleteContainerEvent', function(event, args) {
        $rootScope.$broadcast('deleteContainerTarget', args);
    });
    $rootScope.$on('updateContainerPage', function(event, args) {
        $rootScope.$broadcast('updateContainerPageTarget', args);
    });
    $rootScope.$on('setContainer', function(event, args) {
        $rootScope.$broadcast('setContainerTarget', args);
    });
    $rootScope.$on('createContainerEvent', function(event, args) {
        $rootScope.$broadcast('createContainerTarget', args);
        console.log('broadcast of createContainerTarget');
    });
});

cpmApp.directive('aDisabled', function() {
    return {
        compile: function(tElement, tAttrs, transclude) {
            //Disable ngClick
            tAttrs["ngClick"] = "!("+tAttrs["aDisabled"]+") && ("+tAttrs["ngClick"]+")";

            //Toggle "disabled" to class when aDisabled becomes true
            return function (scope, iElement, iAttrs) {
                scope.$watch(iAttrs["aDisabled"], function(newValue) {
                    if (newValue !== undefined) {
                        iElement.toggleClass("disabled", newValue);
                    }
                });

                //Disable href on click
                iElement.on("click", function(e) {
                    if (scope.$eval(iAttrs["aDisabled"])) {
                        e.preventDefault();
                    }
                });
            };
        }
    };
});


var StopContainerModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.value = value;
    $scope.isLoading = false;
    console.log('in stop container modal stopping ' + value);
    $scope.results = [];

    $scope.ok = function() {
        //$rootScope.$emit('LoadingEvent', { message: "" });
        $scope.isLoading = true;

        var token = $cookieStore.get('cpm_token');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }
        $http.get($cookies.AdminURL + '/admin/stop-pg/' + $scope.value.ID + '.' + token).success(function(data, status, headers, config) {
            $scope.results = data;
            // $rootScope.$emit('DoneLoadingEvent', { message: "" });
            $scope.isLoading = false;
            value.status = 'Completed';
            $modalInstance.close();
            $rootScope.$emit('updateContainerPage', {
                message: $scope.value
            });
        }).error(function(data, status, headers, config) {
            console.log('error:StopContainerModalInstance stop ');
        });

    };

    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    };
};

var StartContainerModalInstanceCtrl = function($rootScope, $scope,
    $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.value = value;
    $scope.isLoading = false;
    console.log('in start container modal stopping ' + value);
    $scope.results = [];

    $scope.ok = function() {

        //$rootScope.$emit('LoadingEvent', { message: "" });
        $scope.isLoading = true;
        var token = $cookieStore.get('cpm_token');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }
        $http.get($cookies.AdminURL + '/admin/start-pg/' + $scope.value.ID + '.' + token).success(function(data, status, headers, config) {
            $scope.results = data;
            //$rootScope.$emit('DoneLoadingEvent', { message: "" });
            $scope.isLoading = false;
            value.status = 'Completed';
            $modalInstance.close();
            $rootScope.$emit('updateContainerPage', {
                message: $scope.value
            });
        }).error(function(data, status, headers, config) {
            console.log('error:StartContainerModalInstance start ');
        });

    };

    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    };
};


var CreateContainerModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore) {

    $scope.alerts = [];
    $scope.servers = null;
    $scope.Image = 'cpm-node';
    $scope.Profile = 'SM';
    $scope.standalone = 'false';
    $scope.isLoading = false;

    var token = $cookieStore.get('cpm_token');
    if (token === void 0) {
        console.log('cookie was undefined');
        alert('login required');
        return;
    }


    $http.get($cookies.AdminURL + '/servers/' + token).
    success(function(data, status, headers, config) {
        $scope.servers = data;
        console.log('got servers len=' + data.length);
        $scope.myServer = $scope.servers[0];
    }).error(function(data, status, headers, config) {
        console.log('error in fetch of servers');
    });

    $scope.cancel = function() {
        $modalInstance.close();
    };

    $scope.ok = function() {

        console.log('in CreateContainerModalInstanceCtrl');
        console.log('    with Name =' + this.Name);
        console.log('    with Image =' + this.Image);
        console.log('    with Image =' + this.Profile);
        console.log('    with standalone flag =' + this.standalone);
        console.log('    with server =' + this.myServer.ID);

        $scope.isLoading = true;
        //$rootScope.$emit('LoadingEvent', { message: "" });
        //
        var token = $cookieStore.get('cpm_token');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }

        $http.get($cookies.AdminURL + '/provision/' + this.Profile + '.' + this.Image + '.' + this.myServer.ID + '.' + this.Name + "." + this.standalone + "." + token).success(function(data, status, headers, config) {
            $scope.results = data;
            console.log('success in provision');
            $rootScope.$emit('createContainerEvent', {
                message: "hi"
            });
            console.log('after emit in provision');
            //$rootScope.$emit('DoneLoadingEvent', { message: "" });
            $scope.isLoading = false;
            $modalInstance.close();
        }).error(function(data, status, headers, config) {
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
            console.log(data.Error);
        });

    };
};


var DeleteContainerModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.value = value;
    $scope.results = [];
    $scope.alerts = [];

    $scope.ok = function() {
        console.log('in DeleteContainerModalInstanceCtrl with value ' + value);
        $scope.isLoading = true;
        var token = $cookieStore.get('cpm_token');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }

        $http.get($cookies.AdminURL + '/deletenode/' + $scope.value.ID + "." + token).success(function(data, status, headers, config) {
            $scope.results = data;
            console.log('success in delete container modal');
            $scope.isLoading = false;
            $rootScope.$emit('deleteContainerEvent', {
                message: "hi"
            });
            $modalInstance.close();
        }).error(function(data, status, headers, config) {
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
            console.log(error.message);
        });

    };

    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    };
};


//controller gor Get All Containers 
cpmApp.controller('GetAllContainersController', function($rootScope, $scope, $http, $modal, $cookies, $cookieStore) {
    $scope.results = [];
    $scope.tab = 1;

    $scope.isSelected2 = function(checkTab) {
        return $scope.tab === checkTab;
    };

    $scope.selectTab2 = function(setTab) {
        console.log('setting tab to ' + setTab.ID);
        $scope.tab = setTab.ID;
        $scope.activeClass = setTab.ID;
        $scope.currentContainer = setTab;
        $rootScope.$emit('setContainer', {
            message: setTab
        });
    }

    function postit() {
        console.log('in GetAllContainers postit');

        var token = $cookieStore.get('cpm_token');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }

        $http.get($cookies.AdminURL + '/nodes/' + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
        }).
        error(function(data, status, headers, config) {
            console.log('error:UpdateContainerPage.http.get');
        });
    }

    var init = function() {
        console.log('GetAllContainers init called');
        var token = $cookieStore.get('cpm_token');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }
        $http.get($cookieStore.get('AdminURL') + '/nodes/' + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
            console.log('containers has ' + $scope.results.length);
            if ($scope.results.length > 0) {
                console.log('first container is ' + $scope.results[0].Name);
                $rootScope.$emit('setContainer', {
                    message: $scope.results[0]
                });
                $scope.activeClass = $scope.results[0].ID;
		$scope.selectTab2($scope.results[0]);
            }
        }).
        error(function(data, status, headers, config) {
            console.log('error:GetAllContainers:init');
        });
    };

    $rootScope.$on('updateContainerPageTarget', function(event, args) {
        console.log('here in GetAllContainersCtrl target ' + args.message);
        $scope.message = args.message;
        $scope.entryID = $scope.message.ID;
        postit();
    });
    $rootScope.$on('createContainerTarget', function(event, args) {
        console.log('GetAllController createContainerTarget recv ' + args.message);
        init();
    });
    $rootScope.$on('deleteContainerTarget', function(event, args) {
        console.log('GetAllController deleteContainerTarget received ');
        init();
    });
    if ($cookies.AdminURL) {
        init();
    } else {
        alert('CPM AdminURL setting is NOT defined, please update on the Settings page before using CPM');
    }
});



cpmApp.controller('GetContainerController', function($scope, $http, $rootScope, $modal, $cookies, $cookieStore, $routeParams) {
    $scope.isRunning;
    $scope.isFound;
    $scope.results = [];
    $scope.clusters = [];
    $scope.myCluster = [];
    $scope.myServer = [];
    $scope.servers = [];
    $scope.currentContainer = [];

    $scope.handleStart = function() {
        console.log('starting db on ' + $scope.currentContainer.Name);
        var modalInstance = $modal.open({
            templateUrl: 'pages/containerstart.html',
            controller: StartContainerModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.currentContainer;
                }
            }
        });
        modalInstance.result.then(function(response) {});
    };
    $scope.handleStop = function() {
        console.log('stopping db on ' + $scope.currentContainer.Name);
        var modalInstance = $modal.open({
            templateUrl: 'pages/containerstop.html',
            controller: StopContainerModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.currentContainer;
                }
            }
        });
        modalInstance.result.then(function(response) {});
    };

    $scope.handleBackupClick = function() {
        console.log('hi from handleBackupClick id=' + $scope.currentContainer.ID);
        var popupWindow = window.open('pages/backups.html');
        console.log('in app.js setting child container=' + $scope.currentContainer.ID);
        popupWindow.containerid = $scope.currentContainer.ID;
    };

    $scope.handleOfflineClick = function(msg) {
        console.log('user wants ' + msg.Name + ' to go offline');
        var modalInstance = $modal.open({
            templateUrl: 'pages/stopcontainermodal.html',
            controller: StopContainerModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.msg;
                }
            }
        });
    };

    $scope.handleMinusClick = function() {
        console.log('hi from handleMinusClick id=' + $scope.currentContainer.ID);
        var modalInstance = $modal.open({
            templateUrl: 'pages/deletecontainermodal.html',
            controller: DeleteContainerModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.currentContainer;
                }
            }
        });
    };

    $scope.handleMonitorClick = function() {
        console.log('hi from handleMonitorClick id=' + $scope.currentContainer.ID);
        var popupWindow = window.open('pages/containermonitor.html');
        console.log('in app.js setting child container=' + $scope.currentContainer.ID);
        popupWindow.containerid = $scope.currentContainer.ID;
        popupWindow.slidervalue = $scope.slidervalue;
    };

    $scope.handlePlusClick = function() {
        console.log('hi from handlePlusClick');
        var modalInstance = $modal.open({
            templateUrl: 'pages/createcontainermodal.html',
            controller: CreateContainerModalInstanceCtrl
        });
        modalInstance.result.then(function(response) {
            $rootScope.$emit('createContainerEvent', {
                message: 'hi'
            });
        });
    };

    function postit(container) {
        $scope.currentContainer = container;
        console.log('in GetContainer postit with containerid=' + container.ID);
        var token = $cookieStore.get('cpm_token');
        $http.get($cookies.AdminURL + '/node/' + container.ID + "." + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
            console.log('here is ServerID=' + $scope.results.ServerID);
            console.log('here is ClusterID=' + $scope.results.ClusterID);
            console.log('here is Status=' + $scope.results.Status);
	    if ($scope.results.Status == 'RUNNING') {
		$scope.isRunning = true;
	    } else {
		$scope.isRunning = false;
	    }
	    if ($scope.results.Status == 'CONTAINER NOT FOUND') {
		$scope.isFound = false;
	    } else {
		$scope.isFound = true;
	    }

            if ($scope.results.ClusterID == -1) {
                console.log('setting myCluster to unassigne value');
                $scope.myCluster.Name = 'unassigned';
            } else {

                if (token === void 0) {
                    console.log('cookie was undefined');
                    alert('login required');
                    return;
                }
                $http.get($cookies.AdminURL + '/cluster/' + $scope.results.ClusterID + "." + token).success(function(data2, status, headers, config) {
                    $scope.myCluster = data2;
                }).error(function(data, status, headers, config) {
                    console.log('error:GetContainer.postit');
                });
            }

            if (token === void 0) {
                console.log('cookie was undefined');
                alert('login required');
                return;
            }
            $http.get($cookies.AdminURL + '/server/' + $scope.results.ServerID + "." + token).success(function(data, status, headers, config) {
                $scope.myServer = data;
            }).error(function(data, status, headers, config) {
                console.log('error:GetContainerController:http.get2');
            });
        }).
        error(function(data, status, headers, config) {
            console.log('error:GetContainerController.http.geta');
        });



        $rootScope.$broadcast('CurrentContainer', $scope.results);

    }

    $rootScope.$on('setContainerTarget', function(event, args) {
        console.log('GetContainerController setContainerTarget ' + args.message.ID);
        postit(args.message);
    });

    $rootScope.$on('updateContainerPageTarget', function(event, args) {
        console.log('here in GetContain ' + args.message.ID);
        postit(args.message);
    });

});
