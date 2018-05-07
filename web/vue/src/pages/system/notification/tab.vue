<template>
  <el-tabs v-model="activeName" @tab-click="changeTab">
    <el-tab-pane label="邮件" name="email"></el-tab-pane>
    <el-tab-pane label="Slack" name="slack"></el-tab-pane>
    <el-tab-pane label="Webhook" name="webhook"></el-tab-pane>
  </el-tabs>
</template>

<script>
export default {
  name: 'notification-tab',
  data () {
    return {
      activeName: ''
    }
  },
  created () {
    const segments = this.$route.path.split('/')
    if (segments.length !== 4) {
      return 'email'
    }
    this.activeName = segments[3]
  },
  methods: {
    changeTab (item) {
      this.$router.push(`/system/notification/${item.name}`)
    }
  }
}
</script>
