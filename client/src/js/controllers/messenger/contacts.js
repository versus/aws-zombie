angular.module('app')
    .controller('contacts', [
        '$rootScope',
        '$scope',
        'app.user',
        function($rootScope, $scope, user) {


            $rootScope.showNav = true;
            $scope.contatcs = [];
            $scope.search = '';
            $scope.phone = window.localStorage.getItem("phone");

            user.all(function(data){
                $scope.users = data;

                user.get({phone: window.localStorage.getItem("phone")}, function(data){
                    console.log(data);
                    $scope.contatcs = data.contact_list;
                });
            });

            $scope.change = function(){

                $scope.contatcs = [];
                if ($scope.search == '') {
                    return
                }

                for (var i=0;i<$scope.users.length;i++) {
                    if (
                        $scope.users[i].first_name.toLowerCase().search($scope.search.toLowerCase()) != -1 ||
                        $scope.users[i].last_name.toLowerCase().search($scope.search.toLowerCase()) != -1 ||
                        $scope.users[i].phone.toLowerCase().search($scope.search.toLowerCase()) != -1
                    ) {
                        $scope.contatcs.push($scope.users[i]);
                    }
                }
            };
        }
    ]);