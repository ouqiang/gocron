<template>
<div>
  <el-table
    :data="tasks"
    border
    style="width: 100%">
    <el-table-column
      prop="id"
      label="任务ID"
      width="180">
    </el-table-column>
    <el-table-column
      prop="name"
      label="任务名称"
      width="180">
    </el-table-column>
    <el-table-column
      prop="spec"
      label="cron表达式"
      width="180">
    </el-table-column>
  </el-table>
  <el-pagination
    background
    layout="prev, pager, next"
    :total="taskTotal">
  </el-pagination>
</div>
</template>

<script>
import httpClient from '../../utils/httpClient'

export default {
  name: 'task-list',
  data () {
    return {
      tasks: [],
      taskTotal: 0
    }
  },
  created () {
    httpClient.get('/task', (data) => {
      this.tasks = data.data
      this.taskTotal = data.total
    })

    httpClient.batchGet(['/task', '/user'], (a, b) => {
      console.log(a)
      console.log(b)
    })
  }
}
</script>
