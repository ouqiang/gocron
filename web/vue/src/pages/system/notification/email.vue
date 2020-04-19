<template>
  <el-container>
      <div style="padding:20px;display:inline-block;background-color:#fff;width:100%;">
        <el-tabs value="setting">
          <el-tab-pane label="邮件服务器配置" name="setting">
            <el-form ref="form" :model="form" :rules="formRules" label-width="150px" style="padding-top:20px;">
              <el-row>
                <el-col :span="12">
                  <el-form-item label="SMTP服务器" prop="host">
                    <el-input v-model="form.host"></el-input>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="端口" prop="port">
                    <el-input v-model.number="form.port"></el-input>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row>
                <el-col :span="12">
                  <el-form-item label="用户名" prop="user">
                    <el-input v-model="form.user"></el-input>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="密码" prop="password">
                    <el-input v-model="form.password" type="password"></el-input>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="">
                <el-alert
                  title="通知模板支持html"
                  type="info"
                  :closable="false"
                  show-icon>
               </el-alert>
              </el-form-item>

              <el-form-item label="模板" prop="template">
                <el-input
                  type="textarea"
                  :rows="6"
                  placeholder=""
                  v-model="form.template">
                </el-input>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="submit()">保存</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <el-tab-pane label="通知用户" name="user">
            <!-- todo table  -->
            <el-button size="small" type="primary" @click="createUser">
              <i class="el-icon el-icon-plus"></i> 新增用户
            </el-button>
            <el-table
              :data="receivers"
              size="mini"
              style="width: 100%">
              <el-table-column
                prop="email"
                label="邮箱">
              </el-table-column>
              <el-table-column
                prop="username"
                label="姓名"
                width="180">
              </el-table-column>
              <el-table-column label="操作">
                <template slot-scope="scope">
                  <el-button size="mini" type="danger" @click="deleteUser(scope.row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>
      <el-dialog
        title="添加用户"
        :visible.sync="dialogVisible"
        width="30%">
        <el-form :model="form">
          <el-form-item label="用户名" >
            <el-input v-model.trim="username"></el-input>
          </el-form-item>
          <el-form-item label="邮箱地址" >
            <el-input v-model.trim="email"></el-input>
          </el-form-item>
        </el-form>
        <div slot="footer">
          <el-button type="primary" @click="saveUser">确 定</el-button>
          <el-button @click="cancel">取消</el-button>
        </div>
      </el-dialog>
  </el-container>
</template>

<script>
import notificationService from '../../../api/notification'
export default {
  name: 'notification-email',
  data () {
    return {
      form: {
        host: '',
        port: 465,
        user: '',
        password: '',
        template: ''
      },
      formRules: {
        host: [
          {required: true, message: '请输入邮件服务器地址', trigger: 'blur'}
        ],
        port: [
          {type: 'number', required: true, message: '请输入有效的端口', trigger: 'blur'}
        ],
        user: [
          {required: true, message: '请输入用户email', trigger: 'blur'}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: 'blur'}
        ],
        template: [
          {required: true, message: '请输入通知模板内容', trigger: 'blur'}
        ]
      },
      receivers: [],
      username: '',
      email: '',
      dialogVisible: false
    }
  },
  created () {
    this.init()
  },
  methods: {
    createUser () {
      this.dialogVisible = true
    },
    saveUser () {
      if (this.username === '' || this.email === '') {
        this.$message.error('参数不完整')
        return
      }
      notificationService.createMailUser({
        username: this.username,
        email: this.email
      }, () => {
        this.$message.success('添加成功')
        this.dialogVisible = false
        this.init()
      })
    },
    deleteUser (item) {
      notificationService.removeMailUser(item.id, () => {
        this.init()
      })
    },
    submit () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      notificationService.updateMail(this.form, () => {
        this.$message.success('更新成功')
        this.init()
      })
    },
    cancel () {

    },
    init () {
      this.username = ''
      this.email = ''
      notificationService.mail((data) => {
        this.form.host = data.host
        if (data.port) {
          this.form.port = data.port
        }
        this.form.user = data.user
        this.form.password = data.password
        this.form.template = data.template
        this.receivers = data.mail_users
      })
    }
  }
}
</script>

<style scoped>
  .el-tag + .el-tag {
    margin-left: 10px;
  }
</style>
