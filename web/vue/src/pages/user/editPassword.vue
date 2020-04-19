<template>
    <div
      style="padding:20px"
    >
      <el-form ref="pform" :model="form" :rules="formRules" label-width="100px">
        <el-form-item prop="id" style="display:none">
          <el-input v-model="form.id" type="hidden"></el-input>
        </el-form-item>
        <el-form-item label="原密码" prop="old_password" v-if="isMe">
          <el-input v-model="form.old_password" type="password"></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="form.new_password" type="password"></el-input>
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirm_new_password">
          <el-input v-model="form.confirm_new_password" type="password"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit()">保存</el-button>
          <el-button @click="cancel">取消</el-button>
        </el-form-item>
      </el-form>
  </div>
</template>

<script>
import userService from '../../api/user'
import userInfo from '../../storage/user'

export default {
  name: 'user-edit-my-password',
  props: {
    userid: Number
  },
  data: function () {
    return {
      form: {
        id: this.userid || '',
        new_password: '',
        confirm_new_password: ''
      },
      formRules: {
        old_password: [
          {required: true, message: '请输入原密码', trigger: 'blur'}
        ],
        new_password: [
          {required: true, message: '请输入新密码', trigger: 'blur'}
        ],
        confirm_new_password: [
          {required: true, message: '请再次输入新密码', trigger: 'blur'}
        ]
      }
    }
  },
  computed: {
    isMe () {
      return +userInfo.getUid() === this.form.id
    }
  },
  methods: {
    submit () {
      this.$refs['pform'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      userService.editMyPassword(this.form, () => {
        this.$message.success('保存成功')
        this.saveBtnLoading = false
        this.$refs.pform.resetFields()
        this.$emit('complete')
      })
    },
    cancel () {
      this.$emit('complete')
    }
  }
}
</script>
