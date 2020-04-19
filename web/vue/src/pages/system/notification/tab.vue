<template>
  <div style="font-size:0;margin-top:1px;display: flex;flex-flow:row;background: #fff;">
    <el-menu :default-active="activeName" @select="changeTab" class="sidebar">
      <el-menu-item index="email">
        <i class="el-icon-message"></i>
        <span slot="title">邮件</span>
      </el-menu-item>
      <el-menu-item index="slack">
        <i class="el-icon-chat-line-round"></i>
        <span slot="title">Slack</span>
      </el-menu-item>
      <el-menu-item index="webhook">
        <i class="el-icon-link"></i>
        <span slot="title">Webhook</span>
      </el-menu-item>
    </el-menu>
    <!-- 路由 -->
    <div style="display:inline-block;flex:1" id="notiy">
      <keep-alive>
        <router-view></router-view>
      </keep-alive>
    </div>
    <div style="display:inline-block;width:300px;padding:30px;background:#fff">
    <el-card style="font-size:16px">
      <pre>
    <strong>通知模板支持的变量</strong>

    TaskId    任务ID
    TaskName  任务名称
    Status    任务执行结果状态
    Result    任务执行输出
      </pre>
    </el-card>
    </div>
  </div>
</template>

<script>
export default {
  name: 'notification-tab',
  data () {
    return {}
  },
  computed: {
    activeName () {
      const segments = this.$route.path.split('/')
      return segments[3]
    }
  },
  mounted () {
    // 动态设置高度
    this.dynamicHeightSetter()
    window.addEventListener('resize', this.dynamicHeightSetter)
  },
  methods: {
    changeTab (item) {
      // console.log(item)
      this.$router.push(`/system/notification/${item}`)
    },
    // 动态设置高度
    dynamicHeightSetter () {
      const bodyHeight = window.document.documentElement.clientHeight
      document.getElementById('notiy').style.minHeight = `${bodyHeight - 61}px`
    }
  }
}
</script>

<style>
.sidebar {
  min-height: 100%;
  width:200px;
  display:inline-block;
  vertical-align: top;
}
</style>
