<template>
  <div v-loading="pageLoading">
      <el-form ref="uform" :model="form" :rules="formRules" label-width="100px" style="padding: 0 20px">
        <el-form-item>
          <el-input v-model="form.id" type="hidden"></el-input>
        </el-form-item>
        <el-form-item label="用户名" prop="name">
          <el-input v-model="form.name" ref="uname"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email"></el-input>
        </el-form-item>
        <template v-if="!form.id">
          <el-form-item label="密码" prop="password">
            <el-input v-model="form.password" type="password"></el-input>
          </el-form-item>
          <el-form-item label="确认密码" prop="confirm_password">
            <el-input v-model="form.confirm_password" type="password"></el-input>
          </el-form-item>
        </template>
        <el-form-item label="角色" prop="is_admin">
          <el-radio-group v-model="form.is_admin">
            <el-radio :label="0">普通用户</el-radio>
            <el-radio :label="1">管理员</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit()" style="width:35%;float:right" :loading="saveBtnLoading">保存</el-button>
          <el-button @click="cancel" style="width:35%;float:right;margin-right:20px">取消</el-button>
        </el-form-item>
      </el-form>
  </div>
</template>

<script>
import userService from '../../api/user'
export default {
  name: 'user-edit',
  props: {
    userid: Number
  },
  data: function () {
    return {
      pageLoading: false,
      saveBtnLoading: false,
      form: {
        id: '',
        name: '',
        email: '',
        is_admin: 0,
        password: '',
        confirm_password: '',
        status: 1
      },
      formRules: {
        name: [
          {required: true, message: '请输入用户名', trigger: 'blur'}
        ],
        email: [
          {type: 'email', required: true, message: '请输入有效邮箱地址', trigger: ['blur', 'change']}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: 'blur'}
        ],
        confirm_password: [
          {required: true, message: '请再次输入密码', trigger: 'blur'},
          {
            validator: (rules, value, callback) => {
              if (this.form.password !== this.form.confirm_password) {
                callback(new Error('两次密码必须一致'))
              } else callback()
            },
            trigger: ['blur', 'change']
          }
        ]
      }
    }
  },
  created () {
    const id = this.userid
    if (!id) {
      return
    }
    this.pageLoading = true
    userService.detail(id, (data) => {
      this.pageLoading = false
      if (!data) {
        this.$message.error('数据不存在')
        return
      }
      this.form.id = data.id
      this.form.name = data.name
      this.form.email = data.email
      this.form.is_admin = data.is_admin
      this.form.status = data.status
    })
  },
  mounted () {
    this.$refs.uname.focus()
  },
  methods: {
    submit () {
      this.$refs['uform'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      this.saveBtnLoading = true
      userService.update(this.form, () => {
        this.$message.success('保存成功')
        this.saveBtnLoading = false
        this.resetForm()
        this.$emit('complete')
      })
    },
    cancel () {
      this.resetForm()
      this.$emit('complete')
    },
    resetForm () {
      this.form = {
        id: '',
        name: '',
        email: '',
        is_admin: 0,
        password: '',
        confirm_password: '',
        status: 1
      }
      this.saveBtnLoading = false
    }
  }
}
</script>
