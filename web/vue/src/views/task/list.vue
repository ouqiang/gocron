<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <strong>定时任务</strong>
      </div>
    </template>
    <el-form :inline="true" label-position="right">
      <el-row>
        <el-form-item label="任务ID">
          <el-input v-model.trim="searchParams.id"></el-input>
        </el-form-item>
        <el-form-item label="任务名称">
          <el-input v-model.trim="searchParams.name"></el-input>
        </el-form-item>
        <el-form-item label="任务命令">
          <el-input v-model.trim="searchParams.command"></el-input>
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model.trim="searchParams.tag"></el-input>
        </el-form-item>
        <el-form-item label="执行方式">
          <el-select v-model.trim="searchParams.protocol">
            <el-option label="全部" value=""></el-option>
            <el-option
                v-for="item in protocolList"
                :key="item.value"
                :label="item.label"
                :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="任务节点">
          <el-select v-model="searchParams.host_id">
            <el-option label="全部" value=""></el-option>
            <el-option
                v-for="item in hosts"
                :key="item.id"
                :label="item.alias + ' - ' + item.name + ':' + item.port "
                :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model.trim="searchParams.status">
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
          <el-button type="primary" @click="search()">搜索</el-button>
          <el-button @click="resetSearch()">重置</el-button>
        </el-form-item>
      </el-row>
    </el-form>
    <el-row type="flex" justify="end">
      <el-col :span="2">
        <el-button type="primary" @click="toEdit(null)" v-if="this.$store.getters.user.isAdmin">新增</el-button>
      </el-col>
      <el-col :span="2">
        <el-button type="info" @click="refresh">刷新</el-button>
      </el-col>
    </el-row>
    <el-pagination
        background
        layout="prev, pager, next, sizes, total"
        :total="taskTotal"
        :page-size="20"
        @size-change="changePageSize"
        @current-change="changePage"
        @prev-click="changePage"
        @next-click="changePage">
    </el-pagination>
    <el-table
        :data="tasks"
        tooltip-effect="dark"
        border
        style="width: 100%">
      <el-table-column type="expand">
        <template #default="scope">
          <el-descriptions style="padding: 15px" border>
            <el-descriptions-item label="任务创建时间:">{{ formatTime(scope.row.created) }}</el-descriptions-item>
            <el-descriptions-item label="任务类型:">{{ formatLevel(scope.row.level) }}</el-descriptions-item>
            <el-descriptions-item label="单实例运行:">{{ formatMulti(scope.row.multi) }}</el-descriptions-item>
            <el-descriptions-item label="超时时间:">{{ formatTimeout(scope.row.timeout) }}</el-descriptions-item>
            <el-descriptions-item label="重试次数:"> {{ scope.row.retry_times }}</el-descriptions-item>
            <el-descriptions-item label="重试间隔:">{{ formatRetryTimesInterval(scope.row.retry_interval) }}</el-descriptions-item>
            <el-descriptions-item label="任务节点:">
              <el-tag style="margin: 0 5px 5px 0" v-for="item in scope.row.hosts" :key="item.host_id">
                {{ item.alias }} - {{ item.name }}:{{ item.port }} <br>
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="命令:"> {{ scope.row.command }}</el-descriptions-item>
            <el-descriptions-item label="备注:">{{ scope.row.remark }}</el-descriptions-item>
          </el-descriptions>
        </template>
      </el-table-column>
      <el-table-column
          prop="id"
          label="任务ID">
      </el-table-column>
      <el-table-column
          prop="name"
          label="任务名称"
          width="150">
      </el-table-column>
      <el-table-column label="标签">
        <template #default="scope">
          <el-link type="primary" @click="tagFilter(scope.row.tag)" :underline="false">{{ scope.row.tag }}</el-link>
        </template>

      </el-table-column>
      <el-table-column
          prop="spec"
          label="cron表达式"
          width="120">
      </el-table-column>
      <el-table-column label="下次执行时间" width="160">
        <template #default="scope">
          {{ formatTime(scope.row.next_run_time) }}
        </template>
      </el-table-column>
      <el-table-column
          prop="protocol"
          :formatter="formatProtocol"
          label="执行方式">
      </el-table-column>
      <el-table-column
          label="状态" v-if="this.isAdmin">
        <template #default="scope">
          <el-switch
              v-if="scope.row.level === 1"
              v-model="scope.row.status"
              :width="50"
              :active-value="1"
              :inactive-value="0"
              active-color="#13ce66"
              inactive-color="#ff4949"
              inline-prompt
              active-text="启用"
              inactive-text="禁用"
              @change="changeStatus(scope.row)">
          </el-switch>
        </template>
      </el-table-column>
      <el-table-column label="状态" v-else>
        <template #default="scope">
          <el-switch
              v-if="scope.row.level === 1"
              v-model="scope.row.status"
              :active-value="1"
              :inactive-value="0"
              active-color="#13ce66"
              :disabled="true"
              inactive-color="#ff4949">
          </el-switch>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" v-if="this.isAdmin">
        <template #default="scope">
          <el-row>
            <el-button type="primary" @click="toEdit(scope.row)">编辑</el-button>
            <el-button type="success" @click="runTask(scope.row)">手动执行</el-button>
          </el-row>
          <br>
          <el-row>
            <el-button type="info" @click="jumpToLog(scope.row)">查看日志</el-button>
            <el-button type="danger" @click="remove(scope.row)">删除</el-button>
          </el-row>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
import taskService from '../../api/task'
import format from '@/utils/format'

export default {
  name: 'task-list',
  data() {
    return {
      tasks: [],
      hosts: [],
      taskTotal: 0,
      defaultSearchParams: {
        page_size: 20,
        page: 1,
        id: '',
        protocol: '',
        name: '',
        tag: '',
        host_id: '',
        status: '',
        command: ''
      },
      searchParams: {},
      isAdmin: this.$store.getters.user.isAdmin,
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
          value: '2',
          label: '激活'
        },
        {
          value: '1',
          label: '停止'
        }
      ]
    }
  },
  created() {
    this.searchParams = Object.assign({}, this.defaultSearchParams)
    let params = localStorage.getItem('task:search')
    if (params) {
      params = JSON.parse(params)
      this.searchParams = Object.assign(this.searchParams, params)
    }

    const hostId = this.$route.query.host_id
    if (hostId) {
      this.searchParams.host_id = Number(hostId)
    }
    this.search()
  },
  methods: {
    formatTime: format.formatDatetime,
    formatLevel(value) {
      if (value === 1) {
        return '主任务'
      }
      return '子任务'
    },
    formatTimeout(value) {
      if (value > 0) {
        return value + '秒'
      }
      return '不限制'
    },
    formatRetryTimesInterval(value) {
      if (value > 0) {
        return value + '秒'
      }
      return '系统默认'
    },
    formatMulti(value) {
      if (value > 0) {
        return '否'
      }
      return '是'
    },
    saveSearchParams() {
      localStorage.setItem('task:search', JSON.stringify(this.searchParams))
    },
    changeStatus(item) {
      console.log(item)
      if (item.status) {
        taskService.enable(item.id)
      } else {
        taskService.disable(item.id)
      }
    },
    formatProtocol(row, col) {
      if (row[col.property] === 2) {
        return 'shell'
      }
      if (row.http_method === 1) {
        return 'http-get'
      }
      return 'http-post'
    },
    changePage(page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize(pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    search(callback = null) {
      this.saveSearchParams()
      taskService.list(this.searchParams, (tasks, hosts) => {
        this.tasks = tasks.data
        this.taskTotal = tasks.total
        this.hosts = hosts
        if (callback) {
          callback()
        }
      })
    },
    tagFilter(tag) {
      this.searchParams.tag = tag
      this.search()
    },
    resetSearch() {
      this.searchParams = Object.assign({}, this.defaultSearchParams)
      localStorage.removeItem('task:search')
      this.search()
    },
    runTask(item) {
      this.$appConfirm(() => {
        taskService.run(item.id, () => {
          this.$message.success('任务已开始执行')
        })
      }, true)
    },
    remove(item) {
      this.$appConfirm(() => {
        taskService.remove(item.id, () => {
          this.refresh()
        })
      })
    },
    jumpToLog(item) {
      this.$router.push(`/task/logs?task_id=${item.id}`)
    },
    refresh() {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    toEdit(item) {
      let path = ''
      if (item === null) {
        path = '/task/create'
      } else {
        path = `/task/edit/${item.id}`
      }
      this.$router.push(path)
    }
  }
}
</script>
<style scoped>
.bl-none{
  border-left: none;
}

.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}

.demo-table-expand .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 50%;
}
.el-table__cell{
  padding: 20px 50px;
}
</style>
