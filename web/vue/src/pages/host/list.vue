<template>
  <el-container>
    <el-main>
      <el-form :inline="true" >
        <el-row>
          <el-form-item label="节点ID">
            <el-input v-model.trim="searchParams.id"></el-input>
          </el-form-item>
          <el-form-item label="主机名">
            <el-input v-model.trim="searchParams.name"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="search()">搜索</el-button>
          </el-form-item>
        </el-row>
      </el-form>
      <el-row type="flex" justify="end">
        <el-col :span="2">
          <el-button type="primary" v-if="this.$store.getters.user.isAdmin"  @click="toEdit(null)">新增</el-button>
        </el-col>
        <el-col :span="2">
          <el-button type="info" @click="refresh">刷新</el-button>
        </el-col>
      </el-row>
      <el-pagination
        background
        layout="prev, pager, next, sizes, total"
        :total="hostTotal"
        :page-size="20"
        @size-change="changePageSize"
        @current-change="changePage"
        @prev-click="changePage"
        @next-click="changePage">
      </el-pagination>
      <el-table
        :data="hosts"
        tooltip-effect="dark"
        border
        style="width: 100%">
        <el-table-column
          prop="id"
          label="节点ID">
        </el-table-column>
        <el-table-column
          prop="alias"
          label="节点名称">
        </el-table-column>
        <el-table-column
          prop="name"
          label="主机名">
        </el-table-column>
        <el-table-column
          prop="port"
          label="端口">
        </el-table-column>
        <el-table-column label="查看任务">
          <template slot-scope="scope">
            <el-button type="success" @click="toTasks(scope.row)">查看任务</el-button>
          </template>
        </el-table-column>
        <el-table-column
          prop="remark"
          label="备注">
        </el-table-column>
        <el-table-column label="操作" width="300" v-if="this.isAdmin">
          <template slot-scope="scope">
            <el-row>
              <el-button type="primary" @click="toEdit(scope.row)">编辑</el-button>
              <el-button type="info" @click="ping(scope.row)">测试连接</el-button>
              <el-button type="danger" @click="remove(scope.row)">删除</el-button>
            </el-row>
            <br>
          </template>
        </el-table-column>
      </el-table>
    </el-main>
  </el-container>
</template>

<script>
import hostService from '../../api/host'
export default {
  name: 'host-list',
  data () {
    return {
      hosts: [],
      hostTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        id: '',
        name: '',
        alias: ''
      },
      isAdmin: this.$store.getters.user.isAdmin
    }
  },
  created () {
    this.search()
  },
  methods: {
    changePage (page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize (pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    search (callback = null) {
      hostService.list(this.searchParams, (data) => {
        this.hosts = data.data
        this.hostTotal = data.total
        if (callback) {
          callback()
        }
      })
    },
    remove (item) {
      this.$appConfirm(() => {
        hostService.remove(item.id, () => this.refresh())
      })
    },
    ping (item) {
      hostService.ping(item.id, () => {
        this.$message.success('连接成功')
      })
    },
    toEdit (item) {
      let path = ''
      if (item === null) {
        path = '/host/create'
      } else {
        path = `/host/edit/${item.id}`
      }
      this.$router.push(path)
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    toTasks (item) {
      this.$router.push(
        {
          path: '/task',
          query: {
            host_id: item.id
          }
        })
    }
  }
}
</script>
