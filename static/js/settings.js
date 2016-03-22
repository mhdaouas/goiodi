'use strict';

app.controller('SettingsPageCtrl', function ($scope, $translate, AppConfigConst, AppConfig) {

    // Load actual locales
    var actualAppLocale = localStorage.getItem(AppConfigConst.APP_LOCALE);
    $scope.appLocale = actualAppLocale;

    var actualDictLocale = localStorage.getItem(AppConfigConst.DICT_LOCALE);
    $scope.dictLocale = actualDictLocale;

    // Updates locales when selected from the settings menu
    $scope.changeAppLocale = function (locale) {
        AppConfig.dictLocale = locale;
        $translate.use(locale);
        localStorage.setItem(AppConfigConst.APP_LOCALE, locale);
    };

    $scope.changeDictLocale = function (locale) {
        AppConfig.dictLocale = locale;
        localStorage.setItem(AppConfigConst.DICT_LOCALE, locale);
    };

});
