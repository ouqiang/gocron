<template>
  <el-row>
    <el-col :span="18">
      <el-breadcrumb separator="/" class="header-item">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item v-for="(item,i) in this.$store.state.breadcrumb" v-bind:key="i">{{ item.name }}</el-breadcrumb-item>
      </el-breadcrumb>
    </el-col>
    <el-col :span="6" style="text-align: right;padding-right: 25px">
      <el-dropdown trigger="hover" class="header-item" v-if="this.$store.getters.user.token" style="cursor: pointer">
        <span class="el-dropdown-link"> {{ this.$store.getters.user.username }} </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="changePwd">修改密码</el-dropdown-item>
            <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </el-col>
  </el-row>
</template>
<script>
export default {
  name: 'nav-headers',
  methods: {
    changePwd() {
      this.$router.push({path: '/user/edit-my-password'})
    },
    logout() {
      this.$store.commit('logout')
      this.$router.push({path: '/user/login'})
    }
  }
}
</script>
<style scoped>
.header-item {
  height: var(--el-header-height);
  line-height: var(--el-header-height);
}
</style>