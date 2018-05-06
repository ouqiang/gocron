<template>
  <el-container>
    <el-header>
      <app-header></app-header>
      <app-nav-menu></app-nav-menu>
    </el-header>
    <el-main >
      <div>
        <router-view/>
      </div>
    </el-main>
    <el-footer>
      <app-footer></app-footer>
    </el-footer>
  </el-container>
</template>

<script>

import appHeader from './components/common/header.vue'
import appNavMenu from './components/common/navMenu.vue'
import appFooter from './components/common/footer.vue'

export default {
  name: 'App',
  data () {
    return {}
  },

  components: {
    appHeader,
    appNavMenu,
    appFooter
  }
}
</script>
<style>
  .el-header {
    padding:0;
    margin:0;
  }

</style>
