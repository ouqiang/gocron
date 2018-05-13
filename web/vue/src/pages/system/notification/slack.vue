<template>
  <el-container>
    <system-sidebar></system-sidebar>
    <el-main>
      <notification-tab></notification-tab>
      <el-form ref="form" :model="form" :rules="formRules" label-width="180px" style="width: 700px;">
        <el-form-item label="Slack Webhook URL" prop="url">
          <el-input v-model="form.url"></el-input>
        </el-form-item>
        <el-form-item label="模板" prop="template">
          <el-input
            type="textarea"
            :rows="8"
            placeholder=""
            size="medium"
            v-model="form.template">
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit">保存</el-button>
        </el-form-item>
        <h3>Channel</h3>
        <el-button type="primary" @click="createChannel">新增Channel</el-button> <br><br>
        <el-tag
          v-for="item in channels"
          :key="item.id"
          closable
          @close="deleteChannel(item)"
        >
          {{item.name}}
        </el-tag>
      </el-form>
      <el-dialog
        title=""
        :visible.sync="dialogVisible"
        width="30%">
        <el-form :model="form">
          <el-form-item label="Channel名称" >
            <el-input v-model.trim="channel" v-focus></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveChannel">确 定</el-button>
          </el-form-item>
        </el-form>
      </el-dialog>
    </el-main>
  </el-container>
</template>

<script>
import systemSidebar from '../sidebar'
import notificationTab from './tab'
import notificationService from '../../../api/notification'
export default {
  name: 'notification-slack',
  data () {
    return {
      dialogVisible: false,
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
      },
      channels: [],
      channel: ''
    }
  },
  components: {notificationTab, systemSidebar},
  created () {
    this.init()
  },
  methods: {
    createChannel () {
      this.dialogVisible = true
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
      notificationService.updateSlack(this.form, () => {
        this.$message.success('更新成功')
        this.init()
      })
    },
    saveChannel () {
      if (this.channel === '') {
        this.$message.error('请输入Channel名称')
        return
      }
      notificationService.createSlackChannel(this.channel, () => {
        this.dialogVisible = false
        this.init()
      })
    },
    deleteChannel (item) {
      notificationService.removeSlackChannel(item.id, () => {
        this.init()
      })
    },
    init () {
      this.channel = ''
      notificationService.slack((data) => {
        this.form.url = data.url
        this.form.template = data.template
        this.channels = data.channels
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
