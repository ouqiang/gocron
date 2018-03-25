/**
 * Created by ouqiang on 2017/4/1.
 */

function Util() {
    var util = {};
    var SUCCESS = 0;     // 操作成功
    var FAILURE_MESSAGE = '操作失败';
    util.alertSuccess = function() {
        swal("操作成功", '保存成功', 'success');
    };
    // ajax成功处理
    util.ajaxSuccess = function(response, callback, failureCallback) {
        if (response.code === undefined) {
            swal(FAILURE_MESSAGE, '服务端返回值无法解析', 'error');
            return;
        }
        if (response.code != SUCCESS) {
            swal(FAILURE_MESSAGE, response.message ,'error');
            if (failureCallback !== undefined) {
                failureCallback(response.code, response.message);
            }
            return;
        }
        if (callback !== undefined) {
            callback(response.code, response.message, response.data);
        }
    };
    // ajax错误处理
    util.ajaxFailure = function() {
        // todo 错误处理
        swal(FAILURE_MESSAGE, '操作失败', 'error');
    };
    // get请求
    util.get = function(url, callback) {
        $.get(
            url,
            function(response) {
                util.ajaxSuccess(response, callback);
            },
            'json'
        ).error(util.ajaxFailure);
    };
    // post请求
    util.post = function(url, params, callback, failureCallback) {
        $.post(
            url,
            util.objectTrim(params),
            function(response) {
                util.ajaxSuccess(response, callback, failureCallback);
            },
            'json'
        ).error(util.ajaxFailure);
    };
    // 弹出确认框
    util.confirm = function(message, callback) {
        swal({
                title: '操作确认',
                text: message,
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#3085d6',
                confirmButtonText: '确定',
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
    // 删除确认
    util.removeConfirm = function(url) {
        util.confirm("确定要删除吗?", function () {
            util.post(url, {}, function () {
                location.reload();
            });
        });
    };

    // 剔除对象元素首尾空格
    util.objectTrim = function(fields) {
      for (key in fields) {
          fields[key] = $.trim(fields[key]);
      }

      return fields;
    };

    util.renderTemplate = function($element, data) {
        var template = Handlebars.compile($($element).html());
        var html = template(data);

        return html;
    };

    return util;
}

// 验证select关联字段
function registerSelectFormValidation(type, $form, $select, selectName) {
    $.fn.form.settings.rules[type] = function(value) {
        var success = true;
        var selectedIndex = $($form).form("get value", selectName);
        $($select).find("option").each(function() {
            var value = $(this).val();
            var match = $(this).data("match");
            var validateType = $(this).data("validate-type");
            if (selectedIndex == value && validateType == type && match) {
                var matchValue = $($form).form("get value", match);
                if (!$.trim(matchValue)) {
                    success = false;
                    return false;
                }
            }
        });

        return success;
    };
}

var util = new Util();