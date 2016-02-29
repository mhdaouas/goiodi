'use strict';

app.controller('AddWordPageCtrl', function (API, $scope, $window) {

    $scope.newWord = {
        word : undefined,
        definition : undefined,
    };

    $scope.addWord = function() {
        $scope.newWord.word = $scope.newWord.word.toString().toLowerCase();
        var newWordJSON = JSON.stringify($scope.newWord);
        API.addWord(newWordJSON).success(function(data){
            $window.location.href = '/#/menu/word_list';
        });
    };

});
