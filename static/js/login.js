'use strict';

app.controller('LoginPageCtrl', function (UserAuth, $state, $scope, $ionicPopup, $ionicHistory) {

    // Form user data
    $scope.user = {
        password : undefined,
        username : undefined,
    };

    // Log user in if his entered information is correct after he clicks on "Log-in" button
    $scope.loginUser = function() {
        var loggedUserJSON = JSON.stringify($scope.user);
        UserAuth.loginUser(loggedUserJSON)
            .success(function(data){
                UserAuth.setLogged(true);
                console.log(data);
                $ionicHistory.nextViewOptions({ historyRoot: true }); 
                $state.transitionTo(UserAuth.stateAfterLogin);
                // $location.path("/#/menu/words");
            })
            .error(function(data){
                // Show an alert if user couldn't be added
                var alertPopup = $ionicPopup.alert({
                    title: 'User couldn\'t be added!',
                    template: 'Please check entered information'
                });
            });
    };

});
