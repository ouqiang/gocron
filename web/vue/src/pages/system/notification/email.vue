<template>
  <el-container>
    <system-sidebar></system-sidebar>
    <el-main>
      <notification-tab></notification-tab>
      <el-form ref="form" :model="form" label-width="100px" style="width: 700px;">
        <h3>邮件服务器配置</h3>
        <el-row>
          <el-col :span="12">
            <el-form-item label="SMTP服务器">
              <el-input v-model="form.host"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="端口">
              <el-input v-model="form.port"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户名">
              <el-input v-model="form.user"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码">
              <el-input v-model="form.password"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="模板">
          <el-input
            type="textarea"
            :rows="5"
            placeholder=""
            v-model="form.template">
          </el-input>
        </el-form-item>
        <h3>邮箱用户</h3>
        <el-tag
          v-for="item in form.receivers"
          :key="item.email"
          closable>
          {{item.username}}
        </el-tag>

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
export default {
  name: 'notification-email',
  data () {
    return {
      form: {
        host: '',
        port: 465,
        user: '',
        password: '',
        template: '',
        receivers: [
          {
            username: '欧强',
            email: 'qingqianludao@gmail.com'
          }
        ]
      }
    }
  },
  components: {notificationTab, systemSidebar},
  methods: {
    submit () {}
  }
}
</script>
