<template>
  <el-container>
    <el-header>
      <app-header></app-header>
      <app-nav-menu></app-nav-menu>
    </el-header>
    <el-container>
      <el-aside>
        <app-sidebar></app-sidebar>
      </el-aside>
      <el-main>
        <app-main></app-main>
      </el-main>
    </el-container>
    <el-footer>
      <app-footer></app-footer>
    </el-footer>
  </el-container>
</template>

<script>

import appHeader from './components/common/header.vue'
import appNavMenu from './components/common/navMenu.vue'
import appSidebar from './components/common/sidebar.vue'
import appMain from './components/common/main.vue'
import appFooter from './components/common/footer.vue'

export default {
  name: 'App',
  data () {
    return {}
  },
  components: {
    appHeader,
    appNavMenu,
    appSidebar,
    appMain,
    appFooter
  }
}
</script>
