'use strict';

app.controller('AddWordPageCtrl', ['$scope', 'API', function ($scope, API) {
    // var newWord;
    $scope.newWord = {
        word : undefined,
        definition : undefined,
    };

    $scope.addWord = function() {
        // if (newWordForm.$valid) {
            var newWordJSON = JSON.stringify($scope.newWord);
            API.addWord(newWordJSON).success(function(data){
            });
        // }
    };

}]);
