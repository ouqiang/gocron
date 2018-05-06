<template>
  <el-aside width="150px">
    <el-menu
      :default-active="currentRoute"
      mode="vertical"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b"
      router>
      <el-menu-item index="/system">通知配置</el-menu-item>
      <el-menu-item index="/system/login-log">登录日志</el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script>
export default {
  name: 'system-sidebar',
  data () {
    return {}
  },
  computed: {
    currentRoute () {
      if (this.$route.path === '/system/login-log') {
        return '/system/login-log'
      }
      return '/system'
    }
  }
}
</script>
