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
        addComment:function(data){
            return $http.post(MainService.host + '/comments/add', data, apiHeader);
        },
    };
});

/* User related parameters */
app.factory('UserAuth', function($q, $http, MainService, $timeout) {
    this.logged = false;
    this.stateAfterLogin = undefined;
    return {
        isLogged : function () { 
            return this.logged;
        },
        stateAfterLogin : function () { 
            return this.stateAfterLogin;
        },
        setLogged : function (val) { 
            this.logged = val;
        },
        setStateAfterLogin : function (val) { 
            this.stateAfterLogin = val;
        },
        // API to add a new user (sign-up)
        addUser:function(data){
            return $http.post(MainService.host + '/users/add', data, apiHeader);
                    // .then(function(result) {
                    //     if (result.data.success) {
                    //         logged = true;
                    //         resolve(result.data);
                    //     } else {
                    //         reject(result.data);
                    //     }
                    // });
            // });
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

app.service('LoginModal', function($ionicModal, API, UserAuth) {

    var service = this;
    // var toStateName = null;

    // Form user data
    this.user = {
        password : undefined,
        username : undefined,
    };

    this.show = function() {
        // service.toStateName = toStateName;
        // console.log("toState: ", toStateName);

        $ionicModal.fromTemplateUrl('login_modal.html', {
            scope: null,
            animation: 'slide-in-up',
            controller: 'LoginModalCtrl'
        }).then(function(modal) {
            service.modal = modal;
            service.modal.show();
        });
    };

    this.hide = function() {
        // service.modal.remove();
        service.modal.hide();
    };

    this.loginUser = function() {
        var loggedUserJSON = JSON.stringify(this.user);
        // Call UserAuth
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
      url: "/user/signup",
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
      // resolve: {
      //     load: function($q, LoginModal, User){
      //         var defer = $q.defer();
      //         if(User.isLogged()){
      //             defer.resolve();
      //         } else {
      //             defer.reject("not_logged_in");
      //             LoginModal.show("menu.add_word");
      //         }
      //         return defer.promise;
      //     }
      // },
      views: {
        'menuContent': {
          templateUrl: "add_word.html",
          controller: 'AddWordPageCtrl',
          // resolve: { authenticate: authenticate }
          // onEnter: function($state, User){
          //     console.log("USER LOGGED: ", User.logged);
          //     if(!User.logged){
          //         $state.go('menu.login');
          //     }
          // }
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

  // $locationProvider.html5Mode({
  //     enabled: true,
  //     requireBase: false
  // });

//   $httpProvider.interceptors.push(function($q, $location) {
//       return {
//           'responseError': function(response) {
//               if(response.status === 401 ||
//                  response.status === 403 ||
//                  response.status === 500) {
//                   // LoginModal.show();
//                   // LoginModal(toState);
//                   $location.path('/user/login');
//                   // $state.go('menu.login');
//               }
//               return $q.reject(response);
//           }
//       };
//   });

});

// Check if the user is authentified for some pages
app.run(function ($rootScope, $state, $location, UserAuth) {
    $rootScope.$on('$stateChangeStart', function (event, toState, toParams) {
        var logged = UserAuth.isLogged();
        if (toState.requireLogin && !logged) {
            UserAuth.setStateAfterLogin(toState.name);
            $state.transitionTo("menu.login")
            // stop state change
            event.preventDefault();
        }


        // log on / sign in...
        // if (!isLogged && requireLogin) {
        //     $state.go("menu.login");
        // }
    });
});

app.controller('AppCtrl', function($rootScope, API, $translate) {

    ionic.Platform.ready(function() {
    });

});
