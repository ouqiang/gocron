<template>
  <div style="padding:0 20px" v-loading="pageLoading">
      <el-form ref="hform" :model="form" :rules="formRules" label-width="100px">
        <el-form-item style="visiblity:hidden">
          <el-input v-model="form.id" type="hidden"></el-input>
        </el-form-item>
        <el-form-item label="节点名称" prop="alias">
          <el-input v-model="form.alias" ref="alias"></el-input>
        </el-form-item>
        <el-form-item label="主机名" prop="name">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input v-model.number="form.port"></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            type="textarea"
            :rows="5"
            size="medium"
            width="100"
            v-model="form.remark">
          </el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" style="margin-top:20px;">
        <el-button type="primary" @click="submit()" style="width:35%;float:right;" :loading="saveBtnLoading">保存</el-button>
        <el-button @click="cancel" style="width:35%;float:right;margin-right:20px">取消</el-button>
      </div>
  </div>
</template>

<script>
import hostService from '../../api/host'
export default {
  name: 'editDialog',
  props: {
    hostid: Number
  },
  data: function () {
    return {
      form: {
        id: '',
        name: '',
        port: 5921,
        alias: '',
        remark: ''
      },
      pageLoading: false,
      saveBtnLoading: false,
      formRules: {
        name: [
          {required: true, message: '请输入主机名', trigger: 'blur'}
        ],
        port: [
          {required: true, message: '请输入端口', trigger: 'blur'},
          {type: 'number', message: '端口无效'}
        ],
        alias: [
          {required: true, message: '请输入节点名称', trigger: 'blur'}
        ]
      }
    }
  },
  mounted () {
    this.$refs.alias.focus()
    this.fetchData()
  },
  methods: {
    fetchData () {
      const id = this.hostid
      if (!id) {
        return
      }
      this.pageLoading = true
      hostService.detail(id, (data) => {
        this.pageLoading = false
        if (!data) {
          this.$message.error('数据不存在')
          this.cancel()
          return
        }
        this.form.id = data.id
        this.form.name = data.name
        this.form.port = data.port
        this.form.alias = data.alias
        this.form.remark = data.remark
      })
    },
    submit () {
      this.$refs['hform'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      this.saveBtnLoading = true
      hostService.update(this.form, () => {
        this.saveBtnLoading = false
        this.$message.success('保存成功')
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
        port: 5921,
        alias: '',
        remark: ''
      }
    }
  },
  watch: {
    hostid (newVal, oldVal) {
      if (newVal !== oldVal) {
        this.fetchData()
      }
    }
  }
}
</script>
