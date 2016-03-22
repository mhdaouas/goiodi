'use strict';

app.controller('WordListPageCtrl', function ($scope, $sce, $ionicScrollDelegate, API) {

    // Get all the dictionary words by default (as soon as the word list page
    // is loaded and no word is in the search filter)
    API.getWords().success(function(data){
        $scope.words = data.response;
    });

    // Filter words based on a user entered string
    $scope.filterWords = function() {
        var filterStr = {
            filter_str: $scope.searchFilter
        };
        var filterStrJSON = JSON.stringify(filterStr);
        API.getWordsIncl(filterStrJSON).success(function(data){
            $scope.words = data.response;
        });

        // Place the scroller at the top
        $scope.scrollTop();
    };

    $scope.scrollTop = function() {
        $ionicScrollDelegate.resize();
    };

});
