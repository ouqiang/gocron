/**
 * Created by ouqiang on 2017/4/1.
 */

function Util() {
    var util = {};
    util.post = function(url, params, callback) {
        // 用户认证失败
        var AUTH_ERROR = -1;
        var FAILURE = 1;
        var SUCCESS = 0;
        var FAILURE_MESSAGE = '操作失败';
        $.post(
            url,
            params,
            function(response) {
                if (!response) {

                }
                if (response.code === undefined) {
                    swal(FAILURE_MESSAGE, '服务端返回值无法解析', 'error');
                }
                if (response.code == AUTH_ERROR) {
                    swal(FAILURE_MESSAGE, '请登录后操作', 'error');
                    return;
                }
                if (response.code == FAILURE) {
                    swal(FAILURE_MESSAGE, response.message ,'error');
                    return;
                }
                callback(response.code, response.message, response.data);
            },
            'json'
        )
    };

    return util;
}