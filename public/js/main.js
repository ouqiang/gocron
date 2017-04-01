/**
 * Created by ouqiang on 2017/4/1.
 */

function Util() {
    var util = {}
    util.post = function(url, params, callback) {
        // 用户认证失败
        var AUTH_ERROR = -1;
        var FAILURE = 1;
        var SUCCESS = 0;
        $.post(
            url,
            params,
            function(response) {
                if (!response) {

                }
                if (response.code === undefined) {

                }
                if (response.code == AUTH_ERROR) {
                    location.href = '/login';
                    return;
                }
                if (response.code == FAILURE) {
                    return;
                }
                callback(response.code, response.message, response.data);
            },
            'json'
        )
    };

    return util;
}