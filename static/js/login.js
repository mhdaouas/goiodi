'use strict';

app.controller('LoginPageCtrl', function ($scope, UserAuth, $state, $ionicPopup, $ionicHistory, $translate) {

    // Form user data
    $scope.user = {
        password : undefined,
        username : undefined,
    };

    // Log user in if his entered information is correct after he clicks on "Log-in" button
    $scope.loginUser = function(loginForm) {
        if(loginForm.$valid) {
            var loggedUserJSON = JSON.stringify($scope.user);
            UserAuth.loginUser(loggedUserJSON)
                .success(function(data){
                    UserAuth.setLogged(true);
                    UserAuth.setUsername($scope.user.username);
                    console.log(data);
                    $ionicHistory.nextViewOptions({ historyRoot: true }); 
                    var stateAfterLogin = UserAuth.getStateAfterLogin();
                    $state.transitionTo(stateAfterLogin);
                    // $location.path("/#/menu/words");
                })
                .error(function(data){
                    // Show an alert if user couldn't be added
                    var alertPopup = $ionicPopup.alert({
                        title: $translate('LoginProblem'),
                        template: $translate('PleaseCheckEnteredInformation')
                    });
                });
        }
    };

});
