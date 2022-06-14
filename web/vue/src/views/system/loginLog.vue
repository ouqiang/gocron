<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <strong>登录日志</strong>
      </div>
    </template>
      <el-pagination
        background
        layout="prev, pager, next, sizes, total"
        :total="logTotal"
        :page-size="searchParams.page_size"
        @size-change="changePageSize"
        @current-change="changePage"
        @prev-click="changePage"
        @next-click="changePage">
      </el-pagination>
      <el-table
        :data="logs"
        border
        ref="table"
        style="width: 100%">
        <el-table-column
          prop="id"
          label="ID">
        </el-table-column>
        <el-table-column
          prop="username"
          label="用户名">
        </el-table-column>
        <el-table-column
          prop="ip"
          label="登录IP">
        </el-table-column>
        <el-table-column
          label="登录时间"
          width="">
          <template #default="scope">
            {{ format.formatDatetime(scope.row.created) }}
          </template>
        </el-table-column>
      </el-table>

  </el-card>
</template>

<script>
import systemService from '../../api/system'
import format from '@/utils/format'

export default {
  name: 'login-log',
  data () {
    return {
      logs: [],
      logTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1
      }
    }
  },
  created () {
    this.format = format
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
      systemService.loginLogList(this.searchParams, (data) => {
        this.logs = data.data
        this.logTotal = data.total
      })
    }
  }
}
</script>
