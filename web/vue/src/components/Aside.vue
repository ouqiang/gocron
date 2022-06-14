<template>
  <el-menu :default-active="this.$store.state.currentPath" router class="aside-menu">
    <li class="el-menu-item logo-container">
      <a href="/#/" class="logo">
        <img alt="title" src="https://avatars.githubusercontent.com/u/53442552?v=4">
        <div class="title">gocron</div>
      </a>
    </li>
    <el-menu-item index="/dashboard">
      <el-icon>
        <svg xmlns="http://www.w3.org/2000/svg" aria-hidden="true" role="img" width="1em" height="1em"
             preserveAspectRatio="xMidYMid meet" viewBox="0 0 1024 1024">
          <path fill="currentColor" d="M924.8 385.6a446.7 446.7 0 0 0-96-142.4a446.7 446.7 0 0 0-142.4-96C631.1 123.8 572.5 112 512 112s-119.1 11.8-174.4 35.2a446.7 446.7 0 0 0-142.4 96a446.7 446.7 0 0 0-96 142.4C75.8 440.9 64 499.5 64 560c0 132.7 58.3 257.7 159.9 343.1l1.7 1.4c5.8 4.8 13.1 7.5 20.6 7.5h531.7c7.5 0 14.8-2.7 20.6-7.5l1.7-1.4C901.7 817.7 960 692.7 960 560c0-60.5-11.9-119.1-35.2-174.4zM482 232c0-4.4 3.6-8 8-8h44c4.4 0 8 3.6 8 8v80c0 4.4-3.6 8-8 8h-44c-4.4 0-8-3.6-8-8v-80zM270 582c0 4.4-3.6 8-8 8h-80c-4.4 0-8-3.6-8-8v-44c0-4.4 3.6-8 8-8h80c4.4 0 8 3.6 8 8v44zm90.7-204.5l-31.1 31.1a8.03 8.03 0 0 1-11.3 0L261.7 352a8.03 8.03 0 0 1 0-11.3l31.1-31.1c3.1-3.1 8.2-3.1 11.3 0l56.6 56.6c3.1 3.1 3.1 8.2 0 11.3zm291.1 83.6l-84.5 84.5c5 18.7.2 39.4-14.5 54.1a55.95 55.95 0 0 1-79.2 0a55.95 55.95 0 0 1 0-79.2a55.87 55.87 0 0 1 54.1-14.5l84.5-84.5c3.1-3.1 8.2-3.1 11.3 0l28.3 28.3c3.1 3.1 3.1 8.1 0 11.3zm43-52.4l-31.1-31.1a8.03 8.03 0 0 1 0-11.3l56.6-56.6c3.1-3.1 8.2-3.1 11.3 0l31.1 31.1c3.1 3.1 3.1 8.2 0 11.3l-56.6 56.6a8.03 8.03 0 0 1-11.3 0zM846 582c0 4.4-3.6 8-8 8h-80c-4.4 0-8-3.6-8-8v-44c0-4.4 3.6-8 8-8h80c4.4 0 8 3.6 8 8v44z"/>
        </svg>
      </el-icon>
      <span>首页</span>
    </el-menu-item>
    <el-sub-menu index="/task">
      <template #title>
        <el-icon>
          <Timer/>
        </el-icon>
        <span>任务管理</span>
      </template>
      <el-menu-item index="/task">
        <span>定时任务</span>
      </el-menu-item>
      <el-menu-item index="/task/logs">
        <span>任务日志</span>
      </el-menu-item>
    </el-sub-menu>
    <el-menu-item index="/process/index">
      <el-icon>
        <House/>
      </el-icon>
      <span>进程管理</span>
    </el-menu-item>
    <el-sub-menu index="/project">
      <template #title>
        <el-icon>
          <FolderAdd/>
        </el-icon>
        <span>项目管理</span>
      </template>
      <el-menu-item index="/project/index">
        <span>项目管理</span>
      </el-menu-item>
      <el-menu-item index="/host/index">
        <span>节点管理</span>
      </el-menu-item>
    </el-sub-menu>
    <el-menu-item index="/user/index">
      <el-icon>
        <UserFilled/>
      </el-icon>
      <span>用户管理</span>
    </el-menu-item>
    <el-sub-menu index="/system">
      <template #title>
        <el-icon>
          <Management/>
        </el-icon>
        <span>系统管理</span>
      </template>
      <el-menu-item index="/system/notification/email">
        <span>通知配置</span>
      </el-menu-item>
      <el-menu-item index="/system/login-log">
        <span>登录日志</span>
      </el-menu-item>
      <el-menu-item index="/system/login-setting">
        <span>登录认证</span>
      </el-menu-item>
    </el-sub-menu>
  </el-menu>
</template>
<script>
export default {
  name: 'AsideNav',
  methods: {}
}
</script>
<style>
.logo-container {
  line-height: 1rem;
  border-bottom: 2px solid black;
  height: 60px;
  margin-left: -10px;
}

.logo-container .logo {
  position: relative;
  display: flex;
  text-decoration: none;
  color: white;
}

.logo-container .logo img {
  height: 40px;
  width: 40px;
  border-radius: 20px;
}

.logo-container div.title {
  margin-left: 10px;
  margin-top: 10px;
  font-size: 16px;
  font-weight: 700;
}

.aside-menu {
  background-color: #001529;
  height: 100%;
}

.aside-menu .el-menu--inline {
  background-color: #0f2438;
}

.aside-menu .el-menu-item, .aside-menu .el-sub-menu__title {
  display: flex;
  align-items: center;
  color: #bfcbd9;
  list-style: none;
  cursor: pointer;
  position: relative;
  box-sizing: border-box;
  white-space: nowrap;
}

.aside-menu .el-menu-item:hover, .aside-menu .el-sub-menu__title:hover {
  color: #fff;
  background-color: #001529;
}

.aside-menu .el-menu--inline .el-menu-item:hover {
  background-color: #0f2438;
}

.aside-menu .el-menu-item.is-active {
  background-color: #409eff !important;
  color: #fff;
}
</style>
