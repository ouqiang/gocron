<template>
  <el-container>
    <!-- <system-sidebar></system-sidebar> -->
    <div style="width:100%;overflow:scroll;padding:30px">
      <el-table
        :data="logs"
        border
        v-loading="loading"
        ref="table"
        size="mini"
      >
        <el-table-column
          prop="id"
          label="ID"
          align="center">
        </el-table-column>
        <el-table-column
          prop="username"
          label="用户名"
          align="center">
        </el-table-column>
        <el-table-column
          prop="ip"
          label="登录IP"
          align="center">
        </el-table-column>
        <el-table-column
          label="登录时间"
          width=""
          align="center">
          <template slot-scope="scope">
            {{scope.row.created | formatTime}}
          </template>
        </el-table-column>
      </el-table>

      <div style="padding: 20px;padding-right:0;background: #fff;overflow: hidden">
        <el-pagination
        style="display:inline-block;float:right"
        background
        layout="prev, pager, next, sizes, total"
        :total="logTotal"
        :page-size="searchParams.page_size"
        @size-change="changePageSize"
        @current-change="changePage"
        @prev-click="changePage"
        @next-click="changePage">
      </el-pagination>
      </div>
    </div>
  </el-container>
</template>

<script>
import systemService from '../../api/system'
export default {
  name: 'login-log',
  data () {
    return {
      logs: [],
      logTotal: 0,
      loading: false,
      searchParams: {
        page_size: 20,
        page: 1
      }
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
    search () {
      this.loading = true
      systemService.loginLogList(this.searchParams, (data) => {
        this.loading = false
        this.logs = data.data
        this.logTotal = data.total
      })
    }
  }
}
</script>

<style scoped>
.el-main {
  margin: 40px;
  /* background: #f0f2f5 */
}
</style>
