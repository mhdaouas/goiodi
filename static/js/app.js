/* App configuration */

'use strict';

var app = angular.module('starter', ['pascalprecht.translate', 'ionic']);

/* App constants */
app.constant('AppConfigConst', {
    // LocalStorage keys
    // Application language
    APP_LOCALE : 'app_locale',
    // Dictionary language
    DICT_LOCALE : 'dict_locale'
});

/* App configuration parameters */
app.factory('AppConfig', function() {
  return {
      // Dictionary language (default: English)
      appLocale : 'en-us',
      // Dictionary language (default: English)
      dictLocale : 'en-us'
  };
});

/* Language related configuration */

// Set translation parameters
app.config(function ($translateProvider, AppConfigConst) {
    // Get available translations
    for(var key in app_locales) {
        var locale = app_locales[key];
        $translateProvider.translations(locale, app_locale_str[locale]);
    }

    // Set preferred language
    var storedLocale = localStorage.getItem(AppConfigConst.APP_LOCALE); 
    if(storedLocale !== undefined && storedLocale !== '') {
        $translateProvider.preferredLanguage(storedLocale);
    }

    // Set fall-back language
    $translateProvider.fallbackLanguage('en-us');

    $translateProvider.useSanitizeValueStrategy('sanitizeParameters');
});

/* Routing related configuration */

app.factory('MainService', function() {
  return {
      host : 'https://localhost:8083'
  };
});

var apiHdr = {
    headers:{
        'Access-Control-Allow-Headers':'Content-Type',
        'Access-Control-Allow-Origin':'*',
        'Content-Type':'application/json'
    }
};

// Define application API
app.factory('API', ['$http', 'MainService' ,function($http, MainService){
    return {
        // API to get word information (definition, creation time, rating, etc.)
        getWordInfo:function(searchedWord){
            return $http.get(MainService.host + '/word/' + searchedWord, apiHdr);
        },
        // API to get all the dictionary words
        getWords:function(){
            return $http.get(MainService.host + '/words', apiHdr);
        },
        // API to get all the dictionary words containing a specific string
        getWordsIncl:function(data){
            return $http.post(MainService.host + '/words/incl', data, apiHdr);
        },
        // API to add a new word to the dictionary
        addWord:function(data){
            return $http.post(MainService.host + '/words/add', data, apiHdr);
        },
        // API to add a new comment for a word in the dictionary
        addComment:function(data){
            return $http.post(MainService.host + '/comments/add', data, apiHdr);
        },
    }
}]);

// Define URL routes
app.config(function ($stateProvider, $urlRouterProvider) {

  $stateProvider
    .state('menu', {
      url: "/menu",
      abstract: true,
      templateUrl: "menu.html",
      controller: 'MenuCtrl'
    })
    // Word list page
    .state('menu.word_list', {
      url: "/word_list",
      views: {
        'menuContent': {
          templateUrl: "word_list.html",
          controller: 'WordListPageCtrl'
        }
      }
    })
    // Word information page
    .state('menu.word', {
      url: "/word/{searchedWord}",
      views: {
        'menuContent': {
          templateUrl: "word_info.html",
          controller: 'WordInfoPageCtrl'
        }
      }
    })
    // Word list page
    .state('menu.add_word', {
      url: "/add_word",
      views: {
        'menuContent': {
          templateUrl: "add_word.html",
          controller: 'AddWordPageCtrl'
        }
      }
    })
    // Application settings page
    .state('menu.settings', {
      url: "/settings",
      views: {
        'menuContent': {
          templateUrl: "settings.html",
          controller: 'SettingsPageCtrl'
        }
      }
    })
    .state('menu.slidebox', {
      url: "/slidebox",
      views: {
        'menuContent': {
          templateUrl: "slidebox.html",
          controller: 'SlideboxCtrl'
        }
      }
    })
    .state('menu.about', {
      url: "/about",
      views: {
        'menuContent': {
          templateUrl: "about.html"
        }
      }
    });

  $urlRouterProvider.otherwise("menu/word_list");

});

app.controller('SlideboxCtrl', function($scope, $ionicSlideBoxDelegate) {
  $scope.nextSlide = function() {
    $ionicSlideBoxDelegate.next();
  }
});

app.controller('MenuCtrl', function($scope, $ionicSideMenuDelegate, $ionicModal) {
    $ionicModal.fromTemplateUrl('modal.html', function (modal) {
        $scope.modal = modal;
    }, {
        animation: 'slide-in-up'
    });
});

app.controller('AppCtrl', function($scope, $translate) {

    ionic.Platform.ready(function() {
    });

});
