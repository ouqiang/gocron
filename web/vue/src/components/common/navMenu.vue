<template>
  <div v-cloak>
    <el-menu
      :default-active="currentRoute"
      mode="horizontal"
      router
    >
      <el-menu-item index="/">
        <img :src="`../../static/gocron.png`" style="height:60px">
      </el-menu-item>
      <el-submenu index="/task">
        <template slot="title">任务管理</template>
        <el-menu-item index="/task">任务列表</el-menu-item>
        <el-menu-item index="/task-log">任务日志</el-menu-item>
      </el-submenu>
      <el-menu-item index="/host">任务节点</el-menu-item>
      <el-menu-item v-if="this.$store.getters.user.isAdmin" index="/user">用户管理</el-menu-item>
      <el-submenu v-if="this.$store.getters.user.isAdmin" index="/system/notification/email">
        <template slot="title">系统管理</template>
        <el-menu-item index="/system/notification/email">通知配置</el-menu-item>
        <el-menu-item index="/system/login-log">登录日志</el-menu-item>
      </el-submenu>

      <div style="float:right;overflow:hidden;height:60px;">
        <!-- 用户 -->
        <el-menu-item v-if="this.$store.getters.user.token" style="display:inline-block;height:60px">
          <template slot="title">
            <span>{{this.$store.getters.user.username}}</span>
          </template>
        </el-menu-item>
        <el-menu-item
          v-if="this.$store.getters.user.token"
          style="display:inline-block;height:60px"
          @click="dialogVisible = true">
          <template slot="title">
            <el-button type="text">修改密码</el-button>
          </template>
        </el-menu-item>
        <el-menu-item
          v-if="this.$store.getters.user.token"
          style="display:inline-block;height:60px"
          @click="logout">
          <template slot="title">
            <el-button type="text">退出</el-button>
          </template>
        </el-menu-item>
      </div>
    </el-menu>

    <!-- 修改我的密码 -->
    <el-dialog title="修改密码" :visible.sync="dialogVisible" width="30%">
      <!-- 修改密码组件复用 -->
      <edit-password :userid="userID" @complete="dialogVisible = false" />
    </el-dialog>
  </div>
</template>

<script>
import editPassword from '../../pages/user/editPassword'
import userInfo from '../../storage/user'

export default {
  name: 'app-nav-menu',
  components: {
    editPassword
  },
  data () {
    return {
      dialogVisible: false
    }
  },
  computed: {
    currentRoute () {
      if (this.$route.path === '/') {
        return '/task'
      }
      if (this.$route.path.includes('/system/notification/')) {
        return '/system/notification/email'
      }
      return this.$route.path
    },
    userID () {
      return +userInfo.getUid()
    }
  },
  methods: {
    logout () {
      this.$appConfirm(() => {
        this.$store.commit('logout')
        this.$router.push('/user/login')
      }, '登出', '确定要登出吗？')
    }
  }
}
</script>
