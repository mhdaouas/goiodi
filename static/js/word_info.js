'use strict';

app.controller('WordInfoPageCtrl', function ($stateParams, $scope, API, AppConfigConst) {

    API.getWordInfo($stateParams.searchedWord).success(function(data){
        // Get API response containing requested word info
        $scope.searchedWord = data.response;

        // Convert Epoch time to a human-readable date
        var date = new Date($scope.searchedWord.creation_time * 1000);
        var actualLocale = localStorage.getItem(AppConfigConst.APP_LOCALE);
        var creationDate = date.toLocaleString(actualLocale);
        $scope.searchedWord.creation_time = creationDate;
        console.log($scope.searchedWord);
    });

    // Spell the word
    $scope.spellWord = function() {
    };
});
