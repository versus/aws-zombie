angular.module('app')
    .controller('chats', [
        '$rootScope',
        '$scope',
        '$location',
        function($rootScope, $scope, $location) {
            // console.log($scope);
            $rootScope.showNav = true;
            $location.path('/contacts');
        }
    ]);