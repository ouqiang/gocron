<template>
  <el-container>
    <!-- <task-sidebar></task-sidebar> -->
    <div style="width:100%;padding:40px">
      <el-form :inline="true" style="padding: 20px 40px 5px 40px;background: #fff;border-radius: 5px">
        <el-form-item label="任务ID">
          <el-input size="small" v-model.trim="searchParams.task_id"></el-input>
        </el-form-item>
        <el-form-item label="执行方式">
          <el-select size="small" v-model.trim="searchParams.protocol" placeholder="执行方式" clearable>
            <el-option label="全部" value=""></el-option>
            <el-option
            v-for="item in protocolList"
            :key="item.value"
            :label="item.label"
            :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select size="small" v-model.trim="searchParams.status" clearable>
            <el-option label="全部" value=""></el-option>
            <el-option
              v-for="item in statusList"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button size="small" type="primary" @click="search()">搜索</el-button>
        </el-form-item>
        <el-form-item>
          <el-button size="small" type="info" @click="reset()">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 功能 -->
      <div style="margin-top: 20px; padding: 20px 40px;background: #fff;text-align: right;border-radius: 5px;">
        <h1 style="float: left;height:42px;line-height: 42px;margin: 0;display: inline-block;">
          任务日志
        </h1>
        <el-button size="small" type="danger" v-if="this.$store.getters.user.isAdmin" @click="clearLog">清空日志</el-button>

        <!-- icon btn -->
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
          width="150"
        >
          <el-checkbox-group v-model="checkList">
            <el-checkbox
              v-for="col in checkListOrigin"
              :label="col"
              :key="col"
              style="margin-bottom:10px"
              >
            </el-checkbox>
          </el-checkbox-group>
          <el-button slot="reference" type="text" style="margin-left: 10px;color: rgba(0,0,0,.75);font-size: 16px">
            <el-tooltip class="item" effect="dark" content="列设置" placement="top">
              <i class="el-icon-s-tools"></i>
            </el-tooltip>
          </el-button>
        </el-popover>
      </div>

      <!-- 表格 -->
      <el-table
        :data="logs"
        border
        :size="tableSize"
        v-loading="tableLoading"
        style="width: 100%">
        <el-table-column type="expand">
          <template slot-scope="scope">
            <el-form label-position="left">
              <el-form-item>
                  重试次数: {{scope.row.retry_times}} <br>
                  cron表达式: {{scope.row.spec}} <br>
                  命令: {{scope.row.command}}
              </el-form-item>
            </el-form>
          </template>
        </el-table-column>
        <el-table-column
          prop="id"
          label="ID"
          align="center">
        </el-table-column>
        <el-table-column
          prop="task_id"
          label="任务ID"
          align="center">
        </el-table-column>
        <el-table-column
          prop="name"
          label="任务名称"
          width="180"
          align="center"
          v-if="checkList.includes('任务名称')">
        </el-table-column>
        <el-table-column
          prop="protocol"
          label="执行方式"
          align="center"
          :formatter="formatProtocol"
          v-if="checkList.includes('执行方式')">
        </el-table-column>
        <el-table-column
          label="任务节点"
          width="150"
          align="center"
          v-if="checkList.includes('任务节点')">
          <template slot-scope="scope">
            <div v-html="scope.row.hostname">{{scope.row.hostname}}</div>
          </template>
        </el-table-column>
        <el-table-column
          label="执行时长"
          width="250"
          v-if="checkList.includes('执行时长')">
          <template slot-scope="scope">
            执行时长: {{scope.row.total_time > 0 ? scope.row.total_time : 1}}秒<br>
            开始时间: {{scope.row.start_time | formatTime}}<br>
            <span v-if="scope.row.status !== 1">结束时间: {{scope.row.end_time | formatTime}}</span>
          </template>
        </el-table-column>
        <el-table-column
          label="状态"
          align="center"
          v-if="checkList.includes('状态')">
          <template slot-scope="scope">
            <span style="color:red" v-if="scope.row.status === 0">失败</span>
            <span style="color:green" v-else-if="scope.row.status === 1">执行中</span>
            <span v-else-if="scope.row.status === 2">成功</span>
            <span style="color:#4499EE" v-else-if="scope.row.status === 3">取消</span>
          </template>
        </el-table-column>
        <el-table-column
          label="执行结果"
          width="120"
          align="center">
          <template slot-scope="scope">
            <el-button :size="tableSize"
              type="success"
              v-if="scope.row.status === 2"
              @click="showTaskResult(scope.row)">查看结果</el-button>
            <el-button :size="tableSize"
              type="warning"
              v-if="scope.row.status === 0"
              @click="showTaskResult(scope.row)" >查看结果</el-button>
            <el-button :size="tableSize"
              type="danger"
              v-if="this.isAdmin && scope.row.status === 1 && scope.row.protocol === 2"
              @click="stopTask(scope.row)">停止任务
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div style="padding: 20px;background: #fff;overflow: hidden">
        <el-pagination
          style="float:right;"
          background
          layout="prev, pager, next, sizes, total"
          :total="logTotal"
          :page-size="20"
          @size-change="changePageSize"
          @current-change="changePage"
          @prev-click="changePage"
          @next-click="changePage">
        </el-pagination>
      </div>

      <el-dialog title="任务执行结果" :visible.sync="dialogVisible" width="70%" top="5vh">
        <div>
          <pre>{{currentTaskResult.command}}</pre>
        </div>
        <div>
          <pre>{{currentTaskResult.result}}</pre>
        </div>
      </el-dialog>
    </div>
  </el-container>
</template>

<script>
import taskLogService from '../../api/taskLog'

export default {
  name: 'task-log',
  data () {
    return {
      logs: [],
      logTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        task_id: '',
        protocol: '',
        status: ''
      },
      tableSize: 'mini',
      tableLoading: false,
      checkListOrigin: ['任务名称', '执行方式', '任务节点', '执行时长', '状态'],
      checkList: ['ID', '任务ID', '任务名称', '执行方式', '任务节点', '执行时长', '状态'],
      isAdmin: this.$store.getters.user.isAdmin,
      dialogVisible: false,
      currentTaskResult: {
        command: '',
        result: ''
      },
      protocolList: [
        {
          value: '1',
          label: 'http'
        },
        {
          value: '2',
          label: 'shell'
        }
      ],
      statusList: [
        {
          value: '1',
          label: '失败'
        },
        {
          value: '2',
          label: '执行中'
        },
        {
          value: '3',
          label: '成功'
        },
        {
          value: '4',
          label: '取消'
        }
      ]
    }
  },
  created () {
    if (this.$route.query.task_id) {
      this.searchParams.task_id = this.$route.query.task_id
    }
    this.search()
  },
  methods: {
    formatProtocol (row, col) {
      if (row[col.property] === 1) {
        return 'http'
      }
      return 'shell'
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
      this.tableLoading = true
      taskLogService.list(this.searchParams, (data) => {
        this.logs = data.data
        this.logTotal = data.total
        this.tableLoading = false
        if (callback) {
          callback()
        }
      })
    },
    reset () {
      this.searchParams = {
        page_size: 20,
        page: 1,
        task_id: '',
        protocol: '',
        status: ''
      }
      this.search()
    },
    clearLog () {
      this.$appConfirm(() => {
        taskLogService.clear(() => {
          this.searchParams.page = 1
          this.search()
        })
      }, '警告', '确定清空日志？')
    },
    stopTask (item) {
      taskLogService.stop(item.id, item.task_id, () => {
        this.search()
      })
    },
    showTaskResult (item) {
      this.dialogVisible = true
      this.currentTaskResult.command = item.command
      this.currentTaskResult.result = item.result
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    }
  }
}
</script>
<style scoped>
  pre {
    white-space: pre-wrap;
    word-wrap: break-word;
    padding:10px;
    background-color: #4C4C4C;
    color: white;
  }
  .el-main {
    margin: 40px !important;
  }
</style>
