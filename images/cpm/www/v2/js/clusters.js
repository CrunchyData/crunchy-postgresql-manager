// create the module and name it cpmApp
var cpmApp = angular.module('cpmApp.clusters', ['ngRoute', 'ngTable', 'ngCookies']);

cpmApp.controller('clustersController', function($scope, $cookies) {
    console.log('hi from clusters controller');
    $scope.message = 'clusters page.';
    if ($cookies.AdminURL) {} else {
        alert('CPM AdminURL setting is NOT defined, please update on the Settings page before using CPM');
    }

});

cpmApp.run(function($rootScope) {
    $rootScope.$on('LoadingEvent', function(event, args) {
        $rootScope.$broadcast('LoadingEventTarget', args);
    });
    $rootScope.$on('DoneLoadingEvent', function(event, args) {
        $rootScope.$broadcast('DoneLoadingEventTarget', args);
    });
    $rootScope.$on('updateClusterPage', function(event, args) {
        $rootScope.$broadcast('updateClusterPageTarget', args);
    });
    $rootScope.$on('createClusterEvent', function(event, args) {
        $rootScope.$broadcast('createClusterTarget', args);
    });
    $rootScope.$on('deleteClusterEvent', function(event, args) {
        $rootScope.$broadcast('deleteClusterTarget', args);
    });
    $rootScope.$on('configureClusterEvent', function(event, args) {
        $rootScope.$broadcast('configureClusterTarget', args);
    });
});




  var FailoverModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

      $scope.value = value;
      $scope.results = [];

      $scope.ok = function() {
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }

         console.log('in FailoverModalInstanceCtrl with container ' + value.Name + " ID=" + value.ID);
        $http.get($cookies.AdminURL + '/admin/failover/' + $scope.value.ID + "." + token).success(function(data, status, headers, config) {
            $scope.results = data;
         }).
         error(function(data, status, headers, config) {
            console.log('error happended');
         });

         value.status = 'Completed';
         $modalInstance.close();
         $rootScope.$emit('deleteClusterEvent', {
            message: ""
         });
      };

      $scope.cancel = function() {
         $modalInstance.dismiss('cancel');
      };
   };


var CreateClusterModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore) {

    $scope.ClusterType = 'asynchronous';
    $scope.alerts = [];

    $scope.cancel = function() {
        $modalInstance.close();
    };

    $scope.ok = function() {
        console.log('in CreateClusterModalInstanceCtrl');
        console.log('    with Name =' + this.Name);
        console.log('    with ClusterType =' + this.ClusterType);
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }

        $http.post($cookies.AdminURL + '/cluster', {
            'Name': this.Name,
            'Status': 'uninitialized',
            'ClusterType': this.ClusterType,
            'Token': token
        }).success(function(data, status, headers, config) {
            $scope.results = data;
            $rootScope.$emit('createClusterEvent', {
                message: ""
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

var AutoClusterModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore) {

    $scope.Name = '';
    $scope.ClusterType = 'asynchronous';
    $scope.ClusterProfile = 'SM';
    $scope.alerts = [];

    $scope.cancel = function() {
        $modalInstance.close();
    };

    $scope.ok = function() {
        console.log('in AutoClusterModalInstanceCtrl');
        console.log('    with Name =' + this.Name);
        console.log('    with ClusterType =' + this.ClusterType);
        console.log('    with ClusterProfile =' + this.ClusterProfile);

        $scope.isLoading = true;
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }

        $http.post($cookies.AdminURL + '/autocluster', {
            'Name': this.Name,
            'ClusterType': this.ClusterType,
            'ClusterProfile': this.ClusterProfile,
            'Token': token
        }).success(function(data, status, headers, config) {
            $scope.results = data;
            $scope.isLoading = false;
            $rootScope.$emit('createClusterEvent', {
                message: ""
            });
            $modalInstance.close();
        }).error(function(data, status, headers, config) {
            console.log('error in Auto Cluster Modal Instance Ctrl');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

    }; //ok function
};

var DeleteClusterModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.value = value;
    $scope.results = [];
    $scope.alerts = [];

    $scope.ok = function() {
        console.log('in DeleteClusterModalInstanceCtrl with value ' + value);

        //$rootScope.$emit('LoadingEvent', { message: "" });
        $scope.isLoading = true;
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }
        $http.get($cookies.AdminURL + '/cluster/delete/' + $scope.value.ID + "." + token).success(function(data, status, headers, config) {
            $scope.results = data;
            $rootScope.$emit('deleteClusterEvent', {
                message: ""
            });
            //$rootScope.$emit('DoneLoadingEvent', { message: "" });
            $scope.isLoading = false;
            $modalInstance.close();
        }).error(function(data, status, headers, config) {
            console.log('error:DeleteClusterModal ');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

        //value.status = 'Completed';

    };

    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    };
};


var ConfigureClusterModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.value = value;
    $scope.results = [];
    $scope.alerts = [];
    $scope.isLoading = false;

    $scope.ok = function() {
        $scope.isLoading = true;
        console.log('in ConfigureClusterModalInstanceCtrl with value ' + value);
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }

        $http.get($cookies.AdminURL + '/cluster/configure/' + $scope.value.ID + "." + token).success(function(data, status, headers, config) {
            $scope.results = data;
            $scope.isLoading = false;
            $modalInstance.close();
        }).error(function(data, status, headers, config) {
            console.log('error:ConfigureClusterModalInstance delete cluster');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

        value.status = 'Completed';
        $rootScope.$emit('configureClusterEvent', {
            message: ""
        });
    };

    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    };
};




//controller gor Get All Clusters
cpmApp.controller('GACController', function($rootScope, $scope, $http, $modal, $cookies, $cookieStore) {
    $scope.results = [];
    $scope.tab = 1;

    console.log('inside GACController!!!');

    $scope.isSelected3 = function(checkTab) {
        return $scope.tab === checkTab;
    };

    $scope.selectTab3 = function(setTab) {
        console.log('setting tab to ' + setTab.ID);
        $scope.tab = setTab.ID;
        $scope.currentCluster = setTab;
        $rootScope.$emit('updateClusterPage', {
            message: setTab
        });
    }

    $scope.handleClick = function(msg) {
        $scope.activeClass = msg.ID;
        $rootScope.$emit('updateClusterPage', {
            message: msg
        });
        console.log(msg.ID);
    };



    function postit() {
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }

        $http.get($cookies.AdminURL + '/clusters/' + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
        }).
        error(function(data, status, headers, config) {
            console.log('error happended');
        });
    };

    var init = function() {
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }
        $http.get($cookies.AdminURL + '/clusters/' + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
            console.log('clusters has ' + $scope.results.length);
            if ($scope.results.length > 0) {
                console.log('first cluster is ' + $scope.results[0].Name);
                $rootScope.$emit('updateClusterPage', {
                    message: $scope.results[0]
                });
                $scope.activeClass = $scope.results[0].ID;
		$scope.selectTab3($scope.results[0]);
            }
        }).
        error(function(data, status, headers, config) {
            console.log('error happended');
        });
    };

    $rootScope.$on('updateClusterPageTarget', function(event, args) {
        console.log('here in GACController ' + args.message.Name);
        postit();
    });

    $rootScope.$on('createClusterTarget', function(event, args) {
        console.log('GACController createClusterTarget received ' + args.message);
        init();
    });
    $rootScope.$on('deleteClusterTarget', function(event, args) {
        console.log('GACController deleteClusterTarget received ');
        init();
    });
    $rootScope.$on('configureClusterTarget', function(event, args) {
        console.log('GACController configureClusterTarget received ');
        init();
    });
    $rootScope.$on('LoadingEventTarget', function(event, args) {
        console.log('GACController LoadingEventTarget received ');
        $scope.isLoading = true;
    });
    $rootScope.$on('DoneLoadingEventTarget', function(event, args) {
        console.log('GACController DoneLoadingEventTarget received ');
        $scope.isLoading = false;
    });

    init();
});



var AddClusterContainerModalInstanceCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value, ngTableParams) {

    $scope.value = value;
    $scope.errorText = '';
    $scope.containers = [];
    $scope.results = [];
    $scope.currentCluster = value;
    $scope.currentMasterID = '';
    $scope.myContainer = [];
    $scope.checkboxes = {
        'checked': true,
        items: {}
    };
    //var containers = [];

    var token = $cookieStore.get('cpmsession');
    if (token === void 0) {
        console.log('cookie was undefined');
        alert('login required');
        return;
    }
    $http.get($cookies.AdminURL + '/nodes/nocluster/' + token).
    success(function(data, status, headers, config) {
        $scope.containers = data;
        console.log('got containers len=' + data.length);
    }).error(function(data, status, headers, config) {
        console.log('error:JoinController.http.geth');
    });


    $scope.tableParams = new ngTableParams({
        page: 1, // show first page
        count: 5, // count per page
        sorting: {
            name: 'asc' // initial sorting
        }
    }, {
        total: 0, // length of data
        getData: function($defer, params) {
            var token = $cookieStore.get('cpmsession');
            if (token === void 0) {
                console.log('cookie was undefined');
                alert('login required');
                return;
            }
            $http.get($cookies.AdminURL + '/nodes/nocluster/' + token).
            success(function(data, status, headers, config) {
                $scope.containers = data;
                console.log('got containers len=' + data.length);
                params.total(data.length);

                data = data.slice((params.page() - 1) * params.count(), params.page() * params.count());
                $defer.resolve(data);
            }).error(function(data, status, headers, config) {
                console.log('error:JoinController.http.geth');
            });
        }
    });


    // watch for check all checkbox
    $scope.$watch('checkboxes.checked', function(value) {
        console.log('checkboxes.checked value=' + value);
        angular.forEach($scope.containers, function(item) {
            if (angular.isDefined(item.ID)) {
                $scope.checkboxes.items[item.ID] = value;
            }
        });
    });

    // watch for data checkboxes
    $scope.$watch('checkboxes.items', function(values) {
        if (!$scope.containers) {
            console.log('here');
            return;
        }
        var checked = 0,
            unchecked = 0,
            total = $scope.containers.length;
        angular.forEach($scope.containers, function(item) {
            checked += ($scope.checkboxes.items[item.ID]) || 0;
            unchecked += (!$scope.checkboxes.items[item.ID]) || 0;
        });
        if ((unchecked == 0) || (checked == 0)) {
            $scope.checkboxes.checked = (checked == total);
        }
        // grayed checkbox
        angular.element(document.getElementById("select_all")).prop("indeterminate", (checked != 0 && unchecked != 0));
    }, true);

    $scope.OnMasterClick = function(value) {
        value.$edit = false;
        $scope.currentMasterID = '';
        //alert('master clicked id=' + value.id);
    };
    $scope.OnStandbyClick = function(value) {
        angular.forEach($scope.containers, function(item) {
            if (angular.isDefined(item.ID)) {
                item.$edit = false;
            }
        });
        value.$edit = true;
        $scope.checkboxes.items[value.ID] = true;
        $scope.currentMasterID = value.ID;
        //alert('standby clicked id=' + value.id);
    };


    $scope.OnSubmitClick = function() {
        var names = '';
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }
        angular.forEach($scope.containers, function(item) {
            if (angular.isDefined(item.ID)) {
                if ($scope.checkboxes.items[item.ID]) {
                    if (item.ID != $scope.currentMasterID) {
                        names += item.ID + "_";
                    }
                }
            }
        });

        $scope.errorText = '';

        if (names == '') {
            $scope.errorText = 'no containers are selected';
            return;
        }

        if ($scope.value.Status == 'initialized') {
            $scope.currentMasterID = '-1';
        }

        if ($scope.currentMasterID == '' && $scope.value.Status != 'initialized') {
            $scope.errorText = 'master is required to be selected cluster status is ' + $scope.value.Status;
            console.log('no master defined');
        } else {
            console.log(names + ' current master=' + $scope.currentMasterID);
            $http.get($cookies.AdminURL + '/event/join-cluster/' + names + '.' + $scope.currentMasterID + '.' + $scope.currentCluster.ID + '.' + token).then(function(result) {
                $scope.results = result;
                console.log('success in join-cluster');
                $rootScope.$emit('updateClusterPage', {
                    message: $scope.currentCluster
                });
                $rootScope.$emit('DoneLoadingEvent', {
                    message: ""
                });
            }, function(result) {
                console.log('error:AddClusterContainerModal.http.get');
                console.log(error.message);
            });

            $modalInstance.close();
        }
    };

    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    };
};


cpmApp.controller('GetClusterController', function($rootScope, $scope, $http, $rootScope, $modal, $cookies, $cookieStore) {
    $scope.results = [];
    $scope.allstatus = [{
        "Value": "initialized"
    }, {
        "Value": "uninitialized"
    }];

    $scope.handleMinusClick = function() {
        console.log('hi from handleMinusClick id=' + $scope.results.ID);
        var modalInstance = $modal.open({
            templateUrl: 'pages/deleteclustermodal.html',
            controller: DeleteClusterModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.results;
                }
            }
        });
    };
    $scope.handleConfigureClick = function() {
        console.log('hi from handleConfigureClick id=' + $scope.results.ID);
        var modalInstance = $modal.open({
            templateUrl: 'pages/configureclustermodal.html',
            controller: ConfigureClusterModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.results;
                }
            }
        });
    };

    $scope.addNewContainer = function() {
        console.log('update cluster add container clicked');
        var modalInstance = $modal.open({
            templateUrl: 'pages/addclustercontainermodal.html',
            controller: AddClusterContainerModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.results;
                }
            }
        });
    };


    $scope.handleAutoClick = function() {
        console.log('hi from handleAutoClick');
        var modalInstance = $modal.open({
            templateUrl: 'pages/autoclustermodal.html',
            controller: AutoClusterModalInstanceCtrl
        });
        modalInstance.result.then(function(response) {});
    };
    $scope.handlePlusClick = function() {
        console.log('hi from handlePlusClick');
        var modalInstance = $modal.open({
            templateUrl: 'pages/createclustermodal.html',
            controller: CreateClusterModalInstanceCtrl
        });
        modalInstance.result.then(function(response) {});
    };

    function postit(clusterid) {
        console.log('in junkit id=' + clusterid);
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }
        $http.get($cookies.AdminURL + '/cluster/' + clusterid + "." + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
            if (data.Status == 'initialized') {
                $scope.results.Status = $scope.allstatus[0].Value;
            } else {
                $scope.results.Status = $scope.allstatus[1].Value;
            }
            $rootScope.$broadcast('CurrentCluster', data);
        }).
        error(function(data, status, headers, config) {
            console.log('error:GetClusterController:http.get');
        });
    }
    $rootScope.$on('updateClusterPageTarget', function(event, args) {
        console.log('here in target ' + args.message);
        $scope.message = args.message;
        $scope.entryID = $scope.message.ID;
        postit(args.message.ID);
    });
});


//controller gor Get All Containers for a given cluster
cpmApp.controller('GetAllContainersForClusterController', function($rootScope, $scope, $http, $modal, $cookies, $cookieStore, $filter, ngTableParams) {
    $scope.selectedContainers = [];
    $scope.results = [];

    var getData = function() {
        return $scope.results;
    };


    $scope.tableParams = new ngTableParams({
        page: 1, // show first page
        count: 10 // count per page
    }, {
        total: $scope.results.length, // length of data
        getData: function($defer, params) {
            $defer.resolve($scope.results.slice((params.page() - 1) * params.count(), params.page() * params.count()));
        },
        $scope: {
            $data: {}
        }
    });


    $scope.failover = function(container) {
        console.log('failover clicked for container ' + container.Name);
        $scope.selectedContainers = container;
        var modalInstance = $modal.open({
            templateUrl: 'pages/failovermodal.html',
            controller: FailoverModalInstanceCtrl,
            resolve: {
                value: function() {
                    return $scope.selectedContainers;
                }
            }
        });
    };

    function postit(v) {
        console.log('in GetAllContainersForCluster postit');
        var token = $cookieStore.get('cpmsession');
        if (token === void 0) {
            console.log('cookie was undefined');
            alert('login required');
            return;
        }
        $http.get($cookies.AdminURL + '/clusternodes/' + v + "." + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
            //console.log('calling tableParams.reload');
            $scope.tableParams.reload();
        }).
        error(function(data, status, headers, config) {
            console.log('error:GetAllContainersForCluster:htt.get');
        });
    }
    $rootScope.$on('updateClusterPageTarget', function(event, args) {
        console.log('here in GetAllContainersForClusterCtrl target ' + args.message.Name);
        $scope.message = args.message;
        $scope.entryID = $scope.message.ID;
        postit(args.message.ID);
    });
});
