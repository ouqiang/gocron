<template>
  <el-container>
    <system-sidebar></system-sidebar>
    <el-main>
      <notification-tab></notification-tab>
      <el-form ref="form" :model="form" :rules="formRules" label-width="100px" style="width: 700px;">
        <el-alert
          title="通知内容推送到指定URL, POST请求, 设置Header[ Content-Type: application/json]"
          type="info"
          :closable="false">
        </el-alert><br>
        <el-form-item label="URL" prop="url">
          <el-input v-model.trim="form.url"></el-input>
        </el-form-item>
        <el-form-item label="模板" prop="template">
          <el-input
            type="textarea"
            :rows="8"
            placeholder=""
            v-model.trim="form.template">
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit()">保存</el-button>
        </el-form-item>
      </el-form>
    </el-main>
  </el-container>
</template>

<script>
import systemSidebar from '../sidebar'
import notificationTab from './tab'
import notificationService from '../../../api/notification'
export default {
  name: 'notification-webhook',
  data () {
    return {
      form: {
        url: '',
        template: ''
      },
      formRules: {
        url: [
          {type: 'url', required: true, message: '请输入有效的通知URL', trigger: 'blur'}
        ],
        template: [
          {required: true, message: '请输入通知模板', trigger: 'blur'}
        ]
      }
    }
  },
  components: {notificationTab, systemSidebar},
  created () {
    this.init()
  },
  methods: {
    submit () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      notificationService.updateWebHook(this.form, () => {
        this.$message.success('更新成功')
        this.init()
      })
    },
    init () {
      notificationService.webhook((data) => {
        this.form.url = data.url
        this.form.template = data.template
      })
    }
  }
}
</script>
