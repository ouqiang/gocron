<template>
  <el-container>
    <div style="width:100%;padding:40px">
      <div style="padding: 20px 40px;background: #fff;text-align: right;border-radius: 5px;">
        <h1 style="float: left;height:42px;line-height: 42px;margin: 0;display: inline-block;">
          用户列表
        </h1>
        <el-button size="small" type="primary"  @click="toEdit(null)" v-if="this.isAdmin">
          新增
        </el-button>
        <!-- icon btn  -->
        <el-popover
          placement="bottom"
          width="150"
        >
          <el-radio-group v-model="tableSize">
            <el-radio style="margin-bottom:10px" label="mini">mini - 最小</el-radio>
            <el-radio style="margin-bottom:10px" label="small">small - 小</el-radio>
            <el-radio label="medium">medium - 正常</el-radio>
          </el-radio-group>
          <el-button slot="reference" type="text" style="margin-left: 40px;color: rgba(0,0,0,.75);font-size: 16px">
            <el-tooltip class="item" effect="dark" content="表格大小" placement="top">
              <i class="el-icon-rank"></i>
            </el-tooltip>
          </el-button>
        </el-popover>
        <el-button type="text" @click="refresh" style="margin-left: 10px;color: rgba(0,0,0,.75);font-size: 16px">
          <el-tooltip class="item" effect="dark" content="刷新" placement="top">
            <i class="el-icon-refresh"></i>
          </el-tooltip>
        </el-button>
        <el-popover
          placement="bottom"
          width="100"
        >
          <el-checkbox-group v-model="checkList">
            <div v-for="col in checkListOrigin" :key="col">
              <el-checkbox
                :label="col"
                style="margin-top:10px"
                >
              </el-checkbox>
            </div>
          </el-checkbox-group>
          <el-button slot="reference" type="text" style="margin-left: 10px;color: rgba(0,0,0,.75);font-size: 16px">
            <el-tooltip class="item" effect="dark" content="列设置" placement="top">
              <i class="el-icon-s-tools"></i>
            </el-tooltip>
          </el-button>
        </el-popover>
      </div>

      <el-table
        :data="users"
        tooltip-effect="dark"
        border
        v-loading="loading"
        :size="tableSize">
        <el-table-column
          prop="id"
          label="用户id"
          align="center">
        </el-table-column>
        <el-table-column
          prop="name"
          label="用户名"
          align="center"
          v-if="checkList.includes('用户名称')">
        </el-table-column>
        <el-table-column
          prop="email"
          label="邮箱"
          align="center"
          v-if="checkList.includes('邮箱')">
        </el-table-column>
        <el-table-column
          prop="is_admin"
          :formatter="formatRole"
          label="角色"
          align="center"
          v-if="checkList.includes('角色')">
        </el-table-column>
        <el-table-column
          label="状态"
          align="center"
          v-if="checkList.includes('状态')">
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.status"
              :active-value="1"
              :inactive-vlaue="0"
              active-color="#13ce66"
              @change="changeStatus(scope.row)"
              inactive-color="#ff4949">
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" v-if="this.isAdmin" align="center">
          <template slot-scope="scope">
            <el-button :size="tableSize" type="primary" @click="toEdit(scope.row)">编辑</el-button>
            <el-button :size="tableSize" type="success" @click="editPassword(scope.row)">修改密码</el-button>
            <el-button :size="tableSize" type="danger" @click="remove(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div style="padding: 20px;background: #fff;overflow: hidden">
        <el-pagination
          style="float:right;"
          background
          layout="prev, pager, next, sizes, total"
          :total="userTotal"
          :page-size="searchParams.page_size"
          @size-change="changePageSize"
          @current-change="changePage"
          @prev-click="changePage"
          @next-click="changePage">
        </el-pagination>
      </div>

      <!-- edit or add -->
      <div style="overflow:scroll">
        <el-drawer
          :visible="editDrawerVisible"
          direction="rtl"
          destroy-on-close
          :with-header="false"
          @close="editDrawerVisible = false"
          size="40%">
          <edit :userid="userid" @complete="completeEdit" />
        </el-drawer>
      </div>

      <!-- 密码 -->
      <div style="overflow:scroll">
        <el-drawer
          :visible="editPasswordDrawerVisible"
          direction="rtl"
          destroy-on-close
          :with-header="false"
          @close="editPasswordDrawerVisible = false"
          size="40%">
          <edit-password :userid="userid" @complete="completeEdit" />
        </el-drawer>
      </div>
    </div>
  </el-container>
</template>

<script>
import userService from '../../api/user'
import edit from './edit'
import editPassword from './editPassword'
export default {
  name: 'user-list',
  components: {
    edit,
    editPassword
  },
  data () {
    return {
      users: [],
      userTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1
      },
      userid: undefined,
      loading: false,
      editDrawerVisible: false,
      editPasswordDrawerVisible: false,
      isAdmin: this.$store.getters.user.isAdmin,
      tableSize: 'mini',
      checkListOrigin: ['用户名', '邮箱', '角色', '状态'],
      checkList: ['用户名', '邮箱', '角色', '状态']
    }
  },
  created () {
    this.search()
  },
  methods: {
    changeStatus (item) {
      if (item.status) {
        userService.enable(item.id)
      } else {
        userService.disable(item.id)
      }
    },
    formatRole (row, col) {
      if (row[col.property] === 1) {
        return '管理员'
      }
      return '普通用户'
    },
    changePage (page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize (pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    search (callback = null) {
      this.loading = true
      userService.list(this.searchParams, (data) => {
        this.users = data.data
        this.userTotal = data.total
        this.loading = false
        if (callback) {
          callback()
        }
      })
    },
    remove (item) {
      this.$appConfirm(() => {
        userService.remove(item.id, () => {
          this.refresh()
        })
      })
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    toEdit (item) {
      if (item === null) {
        this.userid = undefined
      } else {
        this.userid = item.id
      }
      this.editDrawerVisible = true
    },
    editPassword (item) {
      this.userid = item.id
      this.editPasswordDrawerVisible = true
    },
    completeEdit () {
      this.editDrawerVisible = false
      this.editPasswordDrawerVisible = false
      this.search()
    }
  }
}
</script>
