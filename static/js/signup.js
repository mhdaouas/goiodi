'use strict';

app.controller('SignupPageCtrl', function (UserAuth, $scope, $state) {

    // Form user data
    $scope.user = {
        email : undefined,
        password : undefined,
        username : undefined,
    };

    // Add user in DB if the "Sign-up" button is clicked
    $scope.addUser = function() {
        var newUserJSON = JSON.stringify($scope.user);
        UserAuth.addUser(newUserJSON).success(function(data){
            // When user is signed up, redirect him to the words page
            $state.go('menu.words');
        });
    };

});
