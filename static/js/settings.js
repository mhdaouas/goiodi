'use strict';

app.controller('SettingsPageCtrl', function ($scope, $translate, AppConfigConst, AppConfig) {

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
