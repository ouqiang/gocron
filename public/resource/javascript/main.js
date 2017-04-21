/**
 * Created by ouqiang on 2017/4/1.
 */

function Util() {
    var util = {};
    var SUCCESS = 0;     // 操作成功
    var FAILURE_MESSAGE = '操作失败';
    util.ajaxSuccess = function(response, callback) {
        if (response.code === undefined) {
            swal(FAILURE_MESSAGE, '服务端返回值无法解析', 'error');
            return;
        }
        if (response.code != SUCCESS) {
            swal(FAILURE_MESSAGE, response.message ,'error');
            return;
        }
        if (callback !== undefined) {
            callback(response.code, response.message, response.data);
        }
    };
    util.ajaxFailure = function() {
        // todo 错误处理
        swal(FAILURE_MESSAGE, '未知错误', 'error');
    };
    util.get = function(url, callback) {
        var SUCCESS = 0;     // 操作成功
        var FAILURE_MESSAGE = '操作失败';
        $.get(
            url,
            function(response) {
                util.ajaxSuccess(response, callback);
            },
            'json'
        ).error(util.ajaxFailure);
    };
    util.post = function(url, params, callback) {
        $.post(
            url,
            params,
            function(response) {
                util.ajaxSuccess(response, callback);
            },
            'json'
        ).error(util.ajaxFailure);
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