/**
 * Created by ouqiang on 2017/4/1.
 */

function Util() {
    var util = {};
    util.post = function(url, params, callback) {
        // 用户认证失败
        var SUCCESS = 0;
        var FAILURE = 1;
        var NOT_FOUND = 2;
        var AUTH_ERROR = 3;
        var FAILURE_MESSAGE = '操作失败';
        $.post(
            url,
            params,
            function(response) {
                if (response.code === undefined) {
                    swal(FAILURE_MESSAGE, '服务端返回值无法解析', 'error');
                }
                if (response.code == AUTH_ERROR) {
                    swal(FAILURE_MESSAGE, response.message, 'error');
                    return;
                }
                if (response.code == NOT_FOUND) {
                    swal(FAILURE_MESSAGE, response.message, 'error');
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
    util.confirm = function(message, callback) {
        swal({
                title: '操作确认',
                text: message,
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#3085d6',
                confirmButtonText: '删除',
                cancelButtonColor: '#d33',
                cancelButtonText: "取消",
                closeOnConfirm: false,
                closeOnCancel: true
            },
            function(isConfirm) {
                if (!isConfirm) {
                    return;
                }
                callback();
            }
        );
    };
    util.removeConfirm = function(url) {
        util.confirm("确定要删除吗?", function () {
            util.post(url, {}, function () {
                location.reload();
            });
        });
    };

    return util;
}

var util = new Util();