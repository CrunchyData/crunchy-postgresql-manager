// create the module and name it cpmApp
var cpmApp = angular.module('cpmApp.settings', ['ngRoute', 'ngCookies']);

cpmApp.run(function($rootScope) {
    $rootScope.$on('deleteRoleEvent', function(event, args) {
        $rootScope.$broadcast('deleteRoleTarget', args);
    });
    $rootScope.$on('deleteUserEvent', function(event, args) {
        $rootScope.$broadcast('deleteUserTarget', args);
    });
});

cpmApp.controller('settingsController', function($rootScope, $scope, $http, $cookies, $modal, $cookieStore) {
    $scope.AdminURL = 'http://cpm-admin.crunchy.lab:8080';
    if ($cookieStore.get('AdminURL')) {} else {
        alert('AdminURL setting is NOT defined, please update before using CPM');
    }

    $scope.alerts = [];
    $scope.items = ['item1', 'item2'];
    $scope.AdminURL = $cookieStore.get('AdminURL');
    $scope.DockerRegistry = 'registry:5000';
    $scope.PGPort = '5432';
    $scope.DomainName = 'crunchy.lab';
    $scope.settings = [];

    $scope.largeCPU = '';
    $scope.largeMEM = '';
    $scope.mediumCPU = '';
    $scope.mediumMEM = '';
    $scope.smallCPU = '';
    $scope.smallMEM = '';

    $scope.CPsmCount = '';
    $scope.CPsmAlgo = '';
    $scope.CPsmMProfile = '';
    $scope.CPsmSProfile = '';
    $scope.CPsmMServer = '';
    $scope.CPsmSServer = '';
    $scope.CPmedCount = '';
    $scope.CPmedAlgo = '';
    $scope.CPmedMProfile = '';
    $scope.CPmedSProfile = '';
    $scope.CPmedMServer = '';
    $scope.CPmedSServer = '';
    $scope.CPlgCount = '';
    $scope.CPlgAlgo = '';
    $scope.CPlgMProfile = '';
    $scope.CPlgSProfile = '';
    $scope.CPlgMServer = '';
    $scope.CPlgSServer = '';

    $scope.roles = [{
        'Name': 'Read-Only Role'
    }, {
        'Name': 'SuperUser Role'
    }];
    $scope.roleName = '';
    $scope.roleIndex = 0;
    $scope.perm1 = '';
    $scope.perm2 = '';
    $scope.perm3 = '';
    $scope.perm4 = '';
    $scope.users = [{
        'Name': 'Jeff'
    }, {
        'Name': 'Bob'
    }];
    $scope.userName = '';
    $scope.userIndex = 0;

    console.log('hi from settings controlle23r');
    this.tab = 1;

    $rootScope.$on('deleteRoleTarget', function(event, args) {
        console.log('role was deleted....reloading roles');
        $scope.getRoles();
        $scope.getUsers();
    });
    $rootScope.$on('deleteUserTarget', function(event, args) {
        console.log('user was deleted....reloading user');
        $scope.getUsers();
    });

    this.selectRole = function(roleIndex) {
        console.log(' set role to ' + roleIndex);
        $scope.roleIndex = roleIndex;
        this.tab = 6;
    };
    this.selectUser = function(userIndex) {
        console.log(' set user to ' + userIndex);
        $scope.userIndex = userIndex;
        this.tab = 7;
    };

    this.selectTab4 = function(setTab) {
        console.log(' set tab to ' + setTab);
	$scope.alerts = [];
        this.tab = setTab;
    };
    $scope.saveSettings = function() {
        console.log(' save settings called');
        var token = $cookieStore.get('cpm_token');

        $http.post($cookieStore.get('AdminURL') + '/savesettings', {
            'DockerRegistry': this.DockerRegistry,
            'PGPort': this.PGPort,
            'DomainName': this.DomainName,
            'Token': token
        }).success(function(data, status, headers, config) {
            console.log('Save settings success');
            $scope.alerts = [{
                type: 'success',
                msg: 'General Settings saved.'
            }];
        }).error(function(data, status, headers, config) {
            console.log('Error in saving settings.');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
    };
    $scope.saveProfiles = function() {
        console.log(' save Profiles called');

        var token = $cookieStore.get('cpm_token');

        $http.post($cookieStore.get('AdminURL') + '/saveprofiles', {
            'SmallCPU': this.smallCPU,
            'SmallMEM': this.smallMEM,
            'MediumCPU': this.mediumCPU,
            'MediumMEM': this.mediumMEM,
            'LargeCPU': this.largeCPU,
            'LargeMEM': this.largeMEM,
            'Token': token
        }).success(function(data, status, headers, config) {
            console.log('save profiles success');
            $scope.alerts = [{
                type: 'success',
                msg: 'Docker profiles saved.'
            }];
        }).error(function(data, status, headers, config) {
            console.log('error in Save Profiles');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

    };

    $scope.closeAlert = function(index) {
        $scope.alerts.splice(index, 1);
    };

    this.isSelected = function(checkTab) {
        //console.log('settings.js isSelected with ' + checkTab);
        if (checkTab == 6) $scope.roleName = 'Role 1';
        if (checkTab == 7) $scope.roleName = 'Role 2';

        return this.tab === checkTab;
    };

    $scope.getUsers = function() {
        console.log(' get users ');
        var token = $cookieStore.get('cpm_token');

        $http.get($cookieStore.get('AdminURL') + '/sec/getusers/' + token).
        success(function(data, status, headers, config) {
            //console.log('recv users len=' + data.length);
            $scope.users = data;
            //console.log('jeff=' + JSON.stringify($scope.users[0]));
        }).error(function(data, status, headers, config) {
            console.log('error:settingsController.getUsers');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
    };

    $scope.getRoles = function() {
        console.log(' get roles ');
        var token = $cookieStore.get('cpm_token');

        $http.get($cookieStore.get('AdminURL') + '/sec/getroles/' + token).
        success(function(data, status, headers, config) {
            //console.log('recv roles len=' + data.length);
            $scope.roles = data;
            //console.log('jeff=' + JSON.stringify($scope.roles[0]));
        }).error(function(data, status, headers, config) {
            console.log('error:settingsController.getRoles');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
    };
    $scope.saveRole = function() {
        console.log(' save role ');
        console.log(' role is ' + $scope.roles[$scope.roleIndex].Name);
        console.log($scope.roles[$scope.roleIndex]);
        var token = $cookieStore.get('cpm_token');

        $scope.roles[$scope.roleIndex].Token = token;

        $http.post($cookieStore.get('AdminURL') + '/sec/updaterole',
            $scope.roles[$scope.roleIndex]
        ).success(function(data, status, headers, config) {
            //console.log('recv users len=' + data.length);
            //$scope.users = data;
            $scope.alerts = [{
                type: 'success',
                msg: 'Role saved.'
            }];
        }).error(function(data, status, headers, config) {
            console.log('error:settingsController.updateRole');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
    };

    $scope.saveUser = function() {
        console.log(' save users ');
        console.log(' user is ' + $scope.users[$scope.userIndex].Name);
        console.log($scope.users[$scope.userIndex]);

        var token = $cookieStore.get('cpm_token');

        $scope.users[$scope.userIndex].Token = token;

        $http.post($cookieStore.get('AdminURL') + '/sec/updateuser',
            $scope.users[$scope.userIndex]
        ).success(function(data, status, headers, config) {
            //console.log('recv users len=' + data.length);
            //$scope.users = data;
            $scope.alerts = [{
                type: 'success',
                msg: 'User saved.'
            }];
        }).error(function(data, status, headers, config) {
            console.log('error:settingsController.getUsers');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
    };

    $scope.deleteUser = function() {
        console.log(' delete user ');
        console.log(' user is ' + $scope.users[$scope.userIndex].Name);
        var modalInstance = $modal.open({
            size: 'sm',
            templateUrl: 'pages/deleteuser.html',
            controller: DeleteUserController,
            resolve: {
                value: function() {
                    return $scope.users[$scope.userIndex].Name;
                }
            }
        });
    };
    $scope.deleteRole = function() {
        console.log(' delete role ');
        console.log(' role is ' + $scope.roles[$scope.roleIndex].Name);
        var modalInstance = $modal.open({
            size: 'sm',
            templateUrl: 'pages/deleterole.html',
            controller: DeleteRoleController,
            resolve: {
                value: function() {
                    return $scope.roles[$scope.roleIndex].Name;
                }
            }
        });
    };

    $scope.getSettings = function() {
        console.log(' get settings ');

        var token = $cookieStore.get('cpm_token');

        $http.get($cookieStore.get('AdminURL') + '/settings/' + token).
        success(function(data, status, headers, config) {
            console.log('recv settings len=' + data.length);
            $scope.settings = data;
            getAllSettings(data);
        }).error(function(data, status, headers, config) {
            console.log('error:settingsController.get.settings');
        });
    };

    var getAllSettings = function(values) {
        console.log(values.length + ' is the settings len');
        for (i = 0; i < values.length; i++) {
            key = values[i].Name;
            value = values[i].Value;
            switch (key) {
                case 'S-DOCKER-PROFILE-CPU':
                    console.log('key=' + key + ' value=' + value);
                    $scope.smallCPU = value;
                    smallCPU = value;
                    break;
                case 'S-DOCKER-PROFILE-MEM':
                    $scope.smallMEM = value;
                    break;
                case 'M-DOCKER-PROFILE-CPU':
                    $scope.mediumCPU = value;
                    break;
                case 'M-DOCKER-PROFILE-MEM':
                    $scope.mediumMEM = value;
                    break;
                case 'L-DOCKER-PROFILE-CPU':
                    $scope.largeCPU = value;
                    break;
                case 'L-DOCKER-PROFILE-MEM':
                    $scope.largeMEM = value;
                    break;
                case 'DOCKER-REGISTRY':
                    $scope.DockerRegistry = value;
                    break;
                case 'PG-PORT':
                    $scope.PGPort = value;
                    break;
                case 'DOMAIN-NAME':
                    $scope.DomainName = value;
                    break;
                case 'CP-SM-COUNT':
                    $scope.CPsmCount = value;
                    break;
                case 'CP-SM-ALGO':
                    $scope.CPsmAlgo = value;
                    break;
                case 'CP-SM-M-PROFILE':
                    $scope.CPsmMProfile = value;
                    break;
                case 'CP-SM-S-PROFILE':
                    $scope.CPsmSProfile = value;
                    break;
                case 'CP-SM-M-SERVER':
                    $scope.CPsmMServer = value;
                    break;
                case 'CP-SM-S-SERVER':
                    $scope.CPsmSServer = value;
                    break;
                case 'CP-MED-COUNT':
                    $scope.CPmedCount = value;
                    break;
                case 'CP-MED-ALGO':
                    $scope.CPmedAlgo = value;
                    break;
                case 'CP-MED-M-PROFILE':
                    $scope.CPmedMProfile = value;
                    break;
                case 'CP-MED-S-PROFILE':
                    $scope.CPmedSProfile = value;
                    break;
                case 'CP-MED-M-SERVER':
                    $scope.CPmedMServer = value;
                    break;
                case 'CP-MED-S-SERVER':
                    $scope.CPmedSServer = value;
                    break;
                case 'CP-LG-COUNT':
                    $scope.CPlgCount = value;
                    break;
                case 'CP-LG-M-PROFILE':
                    $scope.CPlgMProfile = value;
                    break;
                case 'CP-LG-S-PROFILE':
                    $scope.CPlgSProfile = value;
                    break;
                case 'CP-LG-M-SERVER':
                    $scope.CPlgMServer = value;
                    break;
                case 'CP-LG-S-SERVER':
                    $scope.CPlgSServer = value;
                    break;
                case 'CP-LG-ALGO':
                    $scope.CPlgAlgo = value;
                    break;
                default:
                    console.log('error:settings default reached  ' + $scope.settings[i].Name);
                    break;
            }

        }

    };

    $scope.saveSmallClusterProfiles = function() {
        console.log(' save SmallClusterProfiles called');
        var token = $cookieStore.get('cpm_token');
        $http.post($cookieStore.get('AdminURL') + '/saveclusterprofiles', {
            'Size': 'SM',
            'Count': this.CPsmCount,
            'Algo': this.CPsmAlgo,
            'MasterProfile': this.CPsmMProfile,
            'StandbyProfile': this.CPsmSProfile,
            'MasterServer': this.CPsmMServer,
            'StandbyServer': this.CPsmSServer,
            'Token': token
        }).success(function(data, status, headers, config) {
            console.log('save profiles success');
            $scope.alerts = [{
                type: 'success',
                msg: 'Cluster profiles saved.'
            }];
        }).error(function(data, status, headers, config) {
            console.log('error in Save Profiles');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

    };
    $scope.saveMediumClusterProfiles = function() {
        console.log(' save MediumClusterProfiles called');
        var token = $cookieStore.get('cpm_token');
        $http.post($cookieStore.get('AdminURL') + '/saveclusterprofiles', {
            'Size': 'MED',
            'Count': this.CPmedCount,
            'Algo': this.CPmedAlgo,
            'MasterProfile': this.CPmedMProfile,
            'StandbyProfile': this.CPmedSProfile,
            'MasterServer': this.CPmedMServer,
            'StandbyServer': this.CPmedSServer,
            'Token': token
        }).success(function(data, status, headers, config) {
            console.log('save profiles success');
            $scope.alerts = [{
                type: 'success',
                msg: 'Cluster profiles saved.'
            }];
        }).error(function(data, status, headers, config) {
            console.log('error in Save Profiles');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

    };
    $scope.saveLargeClusterProfiles = function() {
        console.log(' save LargeClusterProfiles called');
        var token = $cookieStore.get('cpm_token');
        $http.post($cookieStore.get('AdminURL') + '/saveclusterprofiles', {
            'Size': 'LG',
            'Count': this.CPlgCount,
            'Algo': this.CPlgAlgo,
            'MasterProfile': this.CPlgMProfile,
            'StandbyProfile': this.CPlgSProfile,
            'MasterServer': this.CPlgMServer,
            'StandbyServer': this.CPlgSServer,
            'Token': token
        }).success(function(data, status, headers, config) {
            console.log('save profiles success');
            $scope.alerts = [{
                type: 'success',
                msg: 'Cluster profiles saved.'
            }];
        }).error(function(data, status, headers, config) {
            console.log('error in Save Profiles');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });

    };

    this.addUser = function() {
        console.log(' addUser called');
        var modalInstance = $modal.open({
            size: 'sm',
            templateUrl: 'pages/adduser.html',
            controller: AddUserController,
            resolve: {
                items: function() {
                    return $scope.items;
                }
            }
        });

        modalInstance.result.then(function(alerts) {
            $scope.alerts = alerts;
        }, function() {
            console.log('modal dismissed');
        });

    }
    this.changePassword = function() {
        console.log(' changePassword called');
        var modalInstance = $modal.open({
            size: 'sm',
            templateUrl: 'pages/chpsw.html',
            controller: ChangePasswordController,
            resolve: {
                value: function() {
                    return $scope.users[$scope.userIndex].Name;
                }
            }
        });

        modalInstance.result.then(function(alerts) {
            $scope.alerts = alerts;
        }, function() {
            console.log('modal dismissed');
        });

    }

    this.addRole = function() {
        console.log(' addRole called');
        var modalInstance = $modal.open({
            size: 'sm',
            templateUrl: 'pages/addrole.html',
            controller: AddRoleController,
            resolve: {
                items: function() {
                    return $scope.items;
                }
            }
        });

        modalInstance.result.then(function(alerts) {
            $scope.alerts = alerts;
        }, function() {
            console.log('modal dismissed');
        });

    }


    $scope.getSettings();
    $scope.getUsers();
    $scope.getRoles();

});

var AddUserController = function($rootScope, $scope, $cookies, $cookieStore, $http, $modalInstance) {
    $scope.ID = '';
    $scope.Password = '';
    $scope.results = [];

    console.log('AddUserController called');
    $scope.doSomething = function() {
        console.log(' doSomething called');
    }
    $scope.ok = function() {
        var token = $cookieStore.get('cpm_token');
        console.log(' login ok called id=' + $scope.ID + ' psw=' + $scope.Password);
        $http.post($cookieStore.get('AdminURL') + '/sec/adduser', {
            'Name': $scope.ID,
            'Password': $scope.Password,
            'Token': token
        }).success(function(data, status, headers, config) {
            console.log('add user success');
            $modalInstance.close([{
                type: 'success',
                msg: 'Added user.'
            }]);
            $rootScope.$emit('deleteUserEvent', {
                message: ""
            });
        }).error(function(data, status, headers, config) {
            console.log('Error in adding user.');
            $modalInstance.close([{
                type: 'danger',
                msg: data.Error
            }]);
        });
    }
    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    }

};

var DeleteUserController = function($rootScope, $http, $cookies, $scope, $modalInstance, $cookieStore, value) {
    console.log('DeleteUserController called user=' + value);
    $scope.value = value;
    $scope.ok = function() {
        console.log(' delete user name=' + value)
        var token = $cookieStore.get('cpm_token');
        $http.get($cookieStore.get('AdminURL') + '/sec/deleteuser/' + value + '.' + token).success(function(data, status, headers, config) {

            console.log('user was deleted');
            $scope.alerts = [{
                type: 'success',
                msg: 'User deleted.'
            }];
            $rootScope.$emit('deleteUserEvent', {
                message: ""
            });
        }).error(function(data, status, headers, config) {
            console.log('error:deleteUser');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
        $modalInstance.close('');
    }
    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    }

};


var AddRoleController = function($rootScope, $scope, $cookies, $cookieStore, $http, $modalInstance) {
    $scope.Name = '';
    $scope.results = [];

    console.log('AddRoleController called');
    $scope.ok = function() {
        console.log(' add role name=' + $scope.Name);
        var token = $cookieStore.get('cpm_token');
        $http.post($cookieStore.get('AdminURL') + '/sec/addrole', {
            'Name': $scope.Name,
            'Token': token
        }).success(function(data, status, headers, config) {
            console.log('add role success');
            $scope.alerts = [{
                type: 'success',
                msg: 'Added role.'
            }];
            $rootScope.$emit('deleteRoleEvent', {
                message: ""
            });
        }).error(function(data, status, headers, config) {
            console.log('Error in adding role.');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
        $modalInstance.close('');
    }
    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    }

};

var DeleteRoleController = function($rootScope, $http, $cookies, $scope, $modalInstance, $cookieStore, value) {
    console.log('DeleteRoleController called Name=' + value);
    $scope.value = value;
    $scope.ok = function() {
        var token = $cookieStore.get('cpm_token');
        console.log(' delete role name=' + value)
        $http.get($cookieStore.get('AdminURL') + '/sec/deleterole/' + value + '.' + token).success(function(data, status, headers, config) {

            console.log('role was deleted');
            $scope.alerts = [{
                type: 'success',
                msg: 'Role deleted.'
            }];
            $rootScope.$emit('deleteRoleEvent', {
                message: ""
            });

        }).error(function(data, status, headers, config) {
            console.log('error:deleteRole');
            $scope.alerts = [{
                type: 'danger',
                msg: data.Error
            }];
        });
        $modalInstance.close('');
    }
    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    }

};

var ChangePasswordController = function($rootScope, $scope, $cookies, $cookieStore, $http, $modalInstance, value) {
    $scope.ID = '';
    $scope.Password = '';
    $scope.value = value;
    $scope.ConfirmPassword = '';
    $scope.results = [];

    console.log('ChangePasswordController called');
    $scope.ok = function() {
        var token = $cookieStore.get('cpm_token');

        console.log(' psw=' + $scope.Password + ' psw2=' + $scope.ConfirmPassword);
        if ($scope.Password != $scope.ConfirmPassword) {
            console.log(' passwords did not match ');
            $scope.alerts = [{
                type: 'danger',
                msg: 'passwords did not match'
            }];
            return;
        }

        $http.post($cookieStore.get('AdminURL') + '/sec/cp', {
            'Username': $scope.value,
            'Password': $scope.Password,
            'Token': token
        }).success(function(data, status, headers, config) {
            console.log('chg psw success');
            $modalInstance.close([{
                type: 'success',
                msg: 'Changed password.'
            }]);
        }).error(function(data, status, headers, config) {
            console.log('Error in chg psw.');
            $modalInstance.close([{
                type: 'danger',
                msg: data.Error
            }]);
        });
        //$modalInstance.close('');
    }
    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    }

};
