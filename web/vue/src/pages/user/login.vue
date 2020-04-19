<template>
    <div>
      <el-dialog
        title="用户登录"
        :visible.sync="dialogVisible"
        :close-on-click-modal="false"
        :show-close="false"
        :close-on-press-escape="false"
        width="400px">
        <el-form
          ref="form"
          :model="form"
          label-width="0"
          :rules="formRules"
          style="padding:0 20px"
        >
          <el-form-item prop="username" >
            <!-- <el-col :span="16"> -->
              <el-input v-model.trim="form.username"
                        placeholder="请输入用户名或邮箱">
                         <i slot="prefix" class="el-input__icon el-icon-user"></i>
              </el-input>
            <!-- </el-col> -->
          </el-form-item>
          <el-form-item prop="password">
            <!-- <el-col :span="16"> -->
              <el-input v-model.trim="form.password" type="password" placeholder="请输入密码">
                 <i slot="prefix" class="el-input__icon el-icon-key"></i>
              </el-input>
            <!-- </el-col> -->
          </el-form-item>
          <el-form-item>

            <el-button type="primary" @click="submit" style="width:100%">登录</el-button>
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
          {required: true, message: '请输入用户名', trigger: ['change', 'blur']}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: ['change', 'blur']}
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
