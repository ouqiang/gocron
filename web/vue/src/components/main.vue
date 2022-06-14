<template>
  <el-scrollbar style="width: 100%">
    <el-header class="header-container">
      <Header/>
    </el-header>
    <el-main>
      <router-view></router-view>
    </el-main>
  </el-scrollbar>
</template>
<script>
import Header from './Header'

export default {
  name: 'MainContainer',
  components: {Header},
  methods: {}
}
</script>
<style scoped>
.header-container {
  border-bottom: 1px solid #eee;
  background-color: #fff;
}
</style>