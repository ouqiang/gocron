<template>
  <el-dialog
    title="您访问的页面不存在"
    :visible.sync="dialogVisible"
    :close-on-click-modal="false"
    :show-close="false"
    :close-on-press-escape="false">
    <el-button type="primary" @click="jump">确定</el-button>
  </el-dialog>
</template>

<script>
export default {
  name: 'not-found',
  data () {
    return {
      dialogVisible: true
    }
  },
  methods: {
    jump () {
      this.$router.push('/')
    }
  }
}
</script>
