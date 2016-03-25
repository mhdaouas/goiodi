'use strict';

app.controller('AddWordPageCtrl', function (API, $scope, $state) {

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
            // After the user adds a word, redirect him to the words page
            $state.go('menu.words');
        });
    };

});
