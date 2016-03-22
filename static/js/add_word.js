'use strict';

app.controller('AddWordPageCtrl', function (API, $scope, $window) {

    // Check if the user is authentified before showing the page
    // $scope.initPage = function () {
    //     // If the user is not logged in, re-direct him to log-in page
        // API.checkUserLogged().error(function(){
        //     console.log('User access denied');
        //     event.preventDefault();
        //     $window.location.href = '/#/menu/user/login';
        //     $window.location.reload();
        // });
    // };

    // New word to be entered by the user
    $scope.newWord = {
        word : undefined,
        definition : undefined,
    };

    // Add a new word if the user clicks on the "Add" button
    $scope.addWord = function() {
        $scope.newWord.word = $scope.newWord.word.toString().toLowerCase();
        var newWordJSON = JSON.stringify($scope.newWord);
        API.addWord(newWordJSON).success(function(data){
            $window.location.href = '/#/menu/words';
            $window.location.reload();
        });
    };

});
