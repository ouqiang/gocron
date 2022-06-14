<template>
  <el-container class="top-container">
    <el-aside width="200px">
      <Aside/>
    </el-aside>
    <el-container>
      <MainContainer/>
    </el-container>
  </el-container>
</template>

<script>

import Aside from './components/Aside'
import MainContainer from './components/main'

export default {
  name: 'App',
  components: {
    Aside,
    MainContainer
  }
}
</script>

<style>
body {
  margin: 0;
  padding: 0;
}

.top-container {
  min-height: calc(100vh);
}

.el-table {
  margin: 15px 0;
  /*--el-table-header-bg-color: #f5f7fa;*/
  /*--el-table-header-text-color: #909399;*/

  /*--el-table-expanded-cell-bg-color:red;*/
}

.el-pagination__sizes {
  margin: 0 16px;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  background-color: #f5f7f9;
}
</style>
