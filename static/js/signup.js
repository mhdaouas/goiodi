'use strict';

app.controller('SignupPageCtrl', function (API, $scope, $window) {

    // Form user data
    $scope.user = {
        email : undefined,
        password : undefined,
        username : undefined,
    };

    // Add user in DB if the "Sign-up" button is clicked
    $scope.addUser = function() {
        var newUserJSON = JSON.stringify($scope.user);
        API.addUser(newUserJSON).success(function(data){
            console.log(data);
            $window.location.href = '/#/menu/words';
            $window.location.reload();
        });
    };

});
