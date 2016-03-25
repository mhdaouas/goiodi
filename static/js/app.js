/* App configuration */

'use strict';

var app = angular.module('starter', ['pascalprecht.translate', 'ionic', 'ngMessages']);

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
      // Application default language: English
      appLocale : 'en-us',
      // Dictionary default language: English
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

var apiHeader = {
    headers:{
        'Access-Control-Allow-Headers':'Content-Type',
        'Access-Control-Allow-Origin':'*',
        'Content-Type':'application/json'
    }
};

// Define application API
app.factory('API', function($http, MainService){
    return {
        /* Word related API */

        // API to get word information (definition, creation time, rating, etc.)
        getWordInfo:function(searchedword){
            return $http.get(MainService.host + '/word/' + searchedword, apiHeader);
        },
        // API to get all the dictionary words
        getWords:function(){
            return $http.get(MainService.host + '/words', apiHeader);
        },
        // API to get all the dictionary words containing a specific string
        getWordsIncl:function(data){
            return $http.post(MainService.host + '/words/incl', data, apiHeader);
        },
        // API to add a new word to the dictionary
        addWord:function(data){
            return $http.post(MainService.host + '/words/add', data, apiHeader);
        },
        // API to add a new comment for a word in the dictionary
        addWordComment:function(data){
            return $http.post(MainService.host + '/comments/add', data, apiHeader);
        },
    };
});

/* User related parameters */
app.factory('UserAuth', function($q, $http, MainService, $timeout) {

    // User properties
    this.logged = false;
    this.username = undefined;
    this.stateAfterLogin = undefined;

    return {
        isLogged : function () { 
            return this.logged;
        },
        getStateAfterLogin : function () { 
            return this.stateAfterLogin;
        },
        setLogged : function (val) { 
            this.logged = val;
        },
        setStateAfterLogin : function (val) { 
            this.stateAfterLogin = val;
        },
        setUsername : function (val) { 
            this.username = val;
        },
        // API to add a new user (sign-up)
        addUser:function(data){
            return $http.post(MainService.host + '/users/add', data, apiHeader);
        },
        // API to log a user in
        loginUser:function(data){
            return $http.post(MainService.host + '/user/login', data, apiHeader);
        },
        // Variable to check if the user is authentified
        // checkLogged : function() {
        //     return $http.get(MainService.host + '/user/login/check', apiHeader);
        // },
    };
});

// Define URL routes
app.config(function ($httpProvider, $stateProvider, $urlRouterProvider) {

  $stateProvider
    // General application menu
    .state('menu', {
      url: "/menu",
      abstract: true,
      templateUrl: "menu.html",
      // controller: 'MenuCtrl'
    })
    // Sign-up page
    .state('menu.signup', {
      url: "/signup",
      requireLogin: false,
      views: {
        'menuContent': {
          templateUrl: "signup.html",
          controller: 'SignupPageCtrl'
        }
      }
    })
    // Login page
    .state('menu.login', {
      url: "/login",
      requireLogin: false,
      views: {
        'menuContent': {
          templateUrl: "login.html",
          controller: 'LoginPageCtrl'
        }
      }
    })
    // Word list page
    .state('menu.words', {
      url: "/words",
      requireLogin: false,
      cache: false,
      views: {
        'menuContent': {
          templateUrl: "word_list.html",
          controller: 'WordListPageCtrl',
          resolve: {
             words: function(WordService) {
                 return WordService.getWords();
             }
          }
        }
      }
    })
    // Word information page
    .state('menu.word', {
      url: "/word/{searchedWord}",
      requireLogin: false,
      views: {
        'menuContent': {
          templateUrl: "word_info.html",
          controller: 'WordInfoPageCtrl'
        }
      }
    })
    // Word creation page
    .state('menu.add_word', {
      url: "/words/add",
      requireLogin: true,
      views: {
        'menuContent': {
          templateUrl: "add_word.html",
          controller: 'AddWordPageCtrl',
        }
      }
    })
    // Application settings page
    .state('menu.settings', {
      url: "/settings",
      requireLogin: false,
      views: {
        'menuContent': {
          templateUrl: "settings.html",
          controller: 'SettingsPageCtrl'
        }
      }
    });

  // Default page
  $urlRouterProvider.otherwise("menu/words");

});

// Check if the user is authentified for some pages
app.run(function ($rootScope, $state, UserAuth) {
    $rootScope.$on('$stateChangeStart', function (event, toState, toParams) {
        var logged = UserAuth.isLogged();
        if (toState.requireLogin && !logged) {
            UserAuth.setStateAfterLogin(toState.name);
            $state.transitionTo("menu.login")
            // stop state change
            event.preventDefault();
        }
    });
});

app.service('WordService', function($q, API) {
  return {
    getWords: function() {
      var dfd = $q.defer()

      API.getWords().then(function(result) {
          dfd.resolve(result.data.response)
      });

      return dfd.promise
    }
  }
});

app.controller('AppCtrl', function($rootScope, API, $translate) {

    ionic.Platform.ready(function() {
    });

});
