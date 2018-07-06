<template>
  <el-container>
    <el-main>
      <el-form ref="form" :model="form" :rules="formRules" label-width="100px" style="width: 700px;">
        <h3>数据库配置</h3>
        <el-form-item label="数据库选择" prop="db_type">
          <el-select v-model.trim="form.db_type" @change="update_port">
            <el-option
              v-for="item in dbList"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="主机名" prop="db_host">
              <el-input v-model="form.db_host"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="端口" prop="db_port">
              <el-input v-model.number="form.db_port"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户名" prop="db_username">
              <el-input v-model="form.db_username"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码" prop="db_password">
              <el-input v-model="form.db_password" type="password"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="数据库名称" prop="db_name">
              <el-input v-model="form.db_name" placeholder="如果数据库不存在, 需提前创建"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="表前缀" prop="db_table_prefix">
              <el-input v-model="form.db_table_prefix"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <h3>管理员账号配置</h3>
        <el-row>
          <el-col :span="12">
            <el-form-item label="账号" prop="admin_username">
              <el-input v-model="form.admin_username"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="admin_email">
              <el-input v-model="form.admin_email"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="密码" prop="admin_password">
              <el-input v-model="form.admin_password" type="password"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="确认密码" prop="confirm_admin_password">
              <el-input v-model="form.confirm_admin_password" type="password"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button type="primary" @click="submit()">安装</el-button>
        </el-form-item>
      </el-form>
    </el-main>
  </el-container>
</template>

<script>
import installService from '../../api/install'
export default {
  name: 'index',
  data () {
    return {
      form: {
        db_type: 'mysql',
        db_host: '127.0.0.1',
        db_port: 3306,
        db_username: '',
        db_password: '',
        db_name: '',
        db_table_prefix: '',
        admin_username: '',
        admin_password: '',
        confirm_admin_password: '',
        admin_email: ''
      },
      formRules: {
        db_type: [
          {required: true, message: '请选择数据库', trigger: 'blur'}
        ],
        db_host: [
          {required: true, message: '请输入数据库主机名', trigger: 'blur'}
        ],
        db_port: [
          {type: 'number', required: true, message: '请输入数据库端口', trigger: 'blur'}
        ],
        db_username: [
          {required: true, message: '请输入数据库用户名', trigger: 'blur'}
        ],
        db_password: [
          {required: true, message: '请输入数据库密码', trigger: 'blur'}
        ],
        db_name: [
          {required: true, message: '请输入数据库名称', trigger: 'blur'}
        ],
        admin_username: [
          {required: true, message: '请输入管理员账号', trigger: 'blur'}
        ],
        admin_email: [
          {type: 'email', required: true, message: '请输入管理员邮箱', trigger: 'blur'}
        ],
        admin_password: [
          {required: true, message: '请输入管理员密码', trigger: 'blur'},
          {min: 6, message: '长度至少6个字符', trigger: 'blur'}
        ],
        confirm_admin_password: [
          {required: true, message: '请再次输入管理员密码', trigger: 'blur'},
          {min: 6, message: '长度至少6个字符', trigger: 'blur'}
        ]
      },
      dbList: [
        {
          value: 'mysql',
          label: 'MySQL'
        },
        {
          value: 'postgres',
          label: 'PostgreSql'
        }
      ],
      default_ports: {
        'mysql': 3306,
        'postgres': 5432
      }
    }
  },
  methods: {
    update_port (dbType) {
      console.log(dbType)
      console.log(this.default_ports[dbType])
      this.form['db_port'] = this.default_ports[dbType]
      console.log(this.form['db_port'])
    },
    submit () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      installService.store(this.form, () => {
        this.$router.push('/')
      })
    }
  }
}
</script>
