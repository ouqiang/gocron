<template>
    <div>
      <el-dialog
        title="用户登录"
        :visible.sync="dialogVisible"
        :close-on-click-modal="false"
        :show-close="false"
        :close-on-press-escape="false"
        width="40%">
        <el-form ref="form" :model="form" label-width="80px"
        :rules="formRules">
          <el-form-item label="用户名" prop="username" >
            <el-col :span="16">
              <el-input v-model.trim="form.username"
                        placeholder="请输入用户名或邮箱">
              </el-input>
            </el-col>
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-col :span="16">
              <el-input v-model.trim="form.password" type="password" placeholder="请输入密码"></el-input>
            </el-col>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submit">登录</el-button>
          </el-form-item>
        </el-form>
      </el-dialog>
    </div>
</template>

<script>
import userServcie from '../../api/user'

export default {
  name: 'login',
  data () {
    return {
      form: {
        username: '',
        password: ''
      },
      formRules: {
        username: [
          {required: true, message: '请输入用户名', trigger: 'blur'}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: 'blur'}
        ]
      },
      dialogVisible: true
    }
  },
  methods: {
    submit () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.login()
      })
    },
    login () {
      userServcie.login(this.form.username, this.form.password, (data) => {
        this.$store.commit('setUser', {
          token: data.token,
          uid: data.uid,
          username: data.username,
          isAdmin: data.is_admin
        })
        this.$router.push(this.$route.query.redirect || '/')
      })
    }
  }
}
</script>
