<template>
  <el-container>
    <system-sidebar></system-sidebar>
    <el-main>
      <notification-tab></notification-tab>
      <el-form ref="form" :model="form" label-width="100px" style="width: 700px;">
        <el-form-item label="URL">
          <el-input v-model="form.url"></el-input>
        </el-form-item>
        <el-form-item label="模板">
          <el-input
            type="textarea"
            :rows="5"
            placeholder=""
            v-model="form.template">
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
export default {
  name: 'notification-webhook',
  data () {
    return {
      form: {
        url: '',
        template: ''
      }
    }
  },
  components: {notificationTab, systemSidebar},
  methods: {
    submit () {

    }
  }
}
</script>
