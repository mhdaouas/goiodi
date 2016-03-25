'use strict';

app.controller('WordInfoPageCtrl', function ($stateParams, $sce, $scope, UserAuth, API, AppConfigConst) {

    // New comment
    $scope.newComment = {
        word : undefined,
        creator : UserAuth.username,
        content : undefined,
    };

    // Get word info
    $scope.getWordInfo = function() {
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
    };
    $scope.getWordInfo();

    $scope.epochToReadable = function(time) {
        var date = new Date(time * 1000);
        var actualLocale = localStorage.getItem(AppConfigConst.APP_LOCALE);
        date = date.toLocaleString(actualLocale);
        return $sce.trustAsHtml(date);
    };

    // Add comment
    $scope.addComment = function() {
        $scope.newComment.word = $scope.searchedWord.word;
        var newCommentJSON = JSON.stringify($scope.newComment);
        API.addWordComment(newCommentJSON).success(function(data){
            $scope.getWordInfo();
        });
    };

    // Check user is logged (to show comments)
    $scope.userIsLogged = function() {
        return UserAuth.isLogged();
    };
});
