'use strict';

app.controller('WordListPageCtrl', function ($scope, $sce, $ionicScrollDelegate, API, words) {

    // Get all the dictionary words before the word list page is loaded
    // using a state resolve
    $scope.words = words;

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

    $scope.highlight = function(text, search) {
        if (!search) {
            return $sce.trustAsHtml(text);
        }
        return $sce.trustAsHtml(text.replace(new RegExp(search, 'gi'), '<span class="highlighted-text">$&</span>'));
    };
});
