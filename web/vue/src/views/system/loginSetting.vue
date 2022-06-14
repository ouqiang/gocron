<template>
  <el-card>
      <el-tabs v-model="activeName">
        <el-tab-pane label="Ldap登录认证" name="ldap">
          <el-form ref="form" :model="ldapSetting" :rules="ldapSettingRules" label-width="100px"
                   style="width: 700px;padding-top: 15px">
            <el-row>
              <el-col :span="6">
                <el-form-item label="Enable" prop="enable">
                  <el-switch active-value="1" inactive-value="0" v-model="ldapSetting.enable"></el-switch>
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="Url" prop="url">
                  <el-input v-model="ldapSetting.url"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="绑定DN" prop="bindDn">
                  <el-input v-model.number="ldapSetting.bindDn"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="绑定密码" prop="bindPassword">
                  <el-input type="password" v-model.number="ldapSetting.bindPassword"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="筛选范围" prop="baseDn">
                  <el-input v-model.number="ldapSetting.baseDn"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="筛选规则" prop="filterRule">
                  <el-input v-model.number="ldapSetting.filterRule"></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item>
              <el-button type="primary" @click="ldapSettingSubmit()">提交</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <!--        <el-tab-pane label="oauth2.0登录认证" name="sso">
                  sso
                </el-tab-pane>-->
      </el-tabs>
  </el-card>
</template>

<script>
import httpClient from '../../utils/httpClient'

export default {
  name: 'login-setting',
  data () {
    return {
      activeName: 'ldap',
      ldapSetting: {
        enable: '0',
        url: '',
        bindDn: '',
        bindPassword: '',
        baseDn: 'ou=users,dc=example,dc=com',
        filterRule: '(&(cn={#username}))'
      },
      ldapSettingRules: {}
    }
  },
  mounted () {
    this.renderSetting()
  },
  methods: {
    ldapSettingSubmit () {
      let _this = this
      httpClient.post('/system/ldap/update', this.ldapSetting, function (resp) {
        console.log(resp)
        _this.$message.success('设置成功')
      })
      console.log(this.ldapSetting)
    },
    renderSetting () {
      let _this = this
      httpClient.get('/system/ldap', {}, function (data) {
        Object.assign(_this.ldapSetting, data)

        // _this.ldapSetting.enable = Boolean(_this.ldapSetting.enable)
        // console.log(data)
      })
    }
  }
}
</script>
