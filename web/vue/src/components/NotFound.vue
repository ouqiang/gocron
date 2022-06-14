<template>
  <el-dialog
      title="您访问的页面不存在"
      v-model="dialogVisible"
      :close-on-click-modal="false"
      :show-close="false"
      :close-on-press-escape="false">
    <el-button type="primary" @click="jump">前往首页</el-button>
  </el-dialog>
</template>

<script>
export default {
  name: 'not-found',
  data() {
    return {
      dialogVisible: true
    }
  },
  methods: {
    jump() {
      this.$router.push('/')
    }
  }
}
</script>
Footer
© 2022 GitHub, Inc.