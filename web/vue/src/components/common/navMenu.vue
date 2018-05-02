<template>
  <el-menu
    :default-active="currentRoute"
    mode="horizontal"
    background-color="#545c64"
    text-color="#fff"
    active-text-color="#ffd04b"
    router>
    <el-menu-item index="/task">任务管理</el-menu-item>
    <el-menu-item index="/host">任务节点</el-menu-item>
    <el-menu-item index="/user">用户管理</el-menu-item>
    <el-menu-item index="/system">系统管理</el-menu-item>
  </el-menu>
</template>

<script>

export default {
  name: 'app-nav-menu',
  data () {
    return {}
  },
  computed: {
    currentRoute () {
      if (this.$route.path === '/') {
        return '/task'
      }
      const segments = this.$route.path.split('/')
      return `/${segments[1]}`
    }
  }
}
</script>

<style>
  .el-header {
    padding:0;
  }
</style>
