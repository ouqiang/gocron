<template>
  <el-container>
      <!-- 模板配置 -->
      <div style="padding:20px;display:inline-block;background-color:#fff;width:100%;">
        <el-tabs value="first">
          <el-tab-pane label="模板配置" name="first">
            <el-form ref="form" :model="form" :rules="formRules" label-width="180px" style="padding-top:20px">
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
            </el-form>
          </el-tab-pane>

          <!-- channel -->
          <el-tab-pane label="Channel配置" name="second">
            <!-- table  -->
            <el-button size="small" type="primary" @click="createChannel">
              <i class="el-icon el-icon-plus"></i> 新增Channel
            </el-button>
            <el-table
              :data="channels"
              size="mini"
              style="width: 100%">
              <el-table-column
                prop="id"
                label="ID">
              </el-table-column>
              <el-table-column
                prop="name"
                label="channel名称"
                width="180">
              </el-table-column>
              <el-table-column label="操作">
                <template slot-scope="scope">
                  <el-button size="mini" type="danger" @click="deleteChannel(scope.row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>
      <el-dialog
        title="新增channel"
        :visible.sync="dialogVisible"
        width="30%">
        <el-form :model="form">
          <el-form-item label="Channel名称" >
            <el-input v-model.trim="channel" v-focus></el-input>
          </el-form-item>
        </el-form>
        <div slot="footer">
          <el-button type="primary" @click="saveChannel">确 定</el-button>
          <el-button>取消</el-button>
        </div>
      </el-dialog>
  </el-container>
</template>

<script>
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
