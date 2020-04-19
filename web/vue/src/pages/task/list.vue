<template>
<el-container style="background: #f0f2f5">
  <!-- <task-sidebar></task-sidebar> -->
  <div style="width:100%;padding:40px">
    <!-- 搜索表单 -->
    <el-form :inline="true" style="padding: 20px 40px 5px 40px;background: #fff;border-radius: 5px">
      <el-row>
        <el-form-item label="任务名称">
          <el-input size="small" v-model.trim="searchParams.name"></el-input>
        </el-form-item>
        <el-form-item label="任务ID">
          <el-input size="small" v-model.trim="searchParams.id"></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select size="small" v-model.trim="searchParams.status">
            <el-option label="全部" value=""></el-option>
            <el-option
              v-for="item in statusList"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="执行方式">
          <el-select size="small" v-model.trim="searchParams.protocol">
            <el-option label="全部" value=""></el-option>
            <el-option
              v-for="item in protocolList"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <!-- button -->
        <el-form-item>
          <el-button size="small" type="primary" @click="search()">搜索</el-button>
        </el-form-item>
        <el-form-item>
          <el-button size="small" type="info" @click="reset()">重置</el-button>
        </el-form-item>
        <el-form-item>
          <el-button
            type="text"
            @click="moreParams = !moreParams"
            v-if="!moreParams"
          >
            <i class="el-icon el-icon-arrow-down"></i> 展开
          </el-button>
          <el-button
            type="text"
            @click="moreParams = !moreParams"
            v-else
          >
            <i class="el-icon el-icon-arrow-up"></i> 收起
          </el-button>
        </el-form-item>
      </el-row>
      <!-- 次行 -->
      <el-row v-if="moreParams">
        <el-form-item label="任务节点">
          <el-select size="small" v-model.trim="searchParams.host_id">
            <el-option label="全部" value=""></el-option>
            <el-option
              v-for="item in hosts"
              :key="item.id"
              :label="item.alias + ' - ' + item.name + ':' + item.port "
              :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-input size="small" v-model.trim="searchParams.tag"></el-input>
        </el-form-item>
      </el-row>
    </el-form>

    <!-- 功能按钮组 -->
    <div style="margin-top: 20px; padding: 20px 40px;background: #fff;text-align: right;border-radius: 5px;">
      <h1 style="float: left;height:42px;line-height: 42px;margin: 0;display: inline-block;">
        任务列表
      </h1>
      <el-button type="primary" size="small" @click="toEdit(null)" v-if="this.$store.getters.user.isAdmin">
        新增
      </el-button>
      <el-button type="success" size="small" @click="batchTrigger" :disabled="!ifMulSelected">批量手动执行</el-button>
      <el-button type="danger" size="small" @click="batchDelete" :disabled="!ifMulSelected">批量删除</el-button>
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
        width="150"
      >
        <el-checkbox-group v-model="checkList">
          <el-checkbox label="任务ID" style="margin-bottom:10px"></el-checkbox>
          <el-checkbox label="任务名称" style="margin-bottom:10px"></el-checkbox>
          <el-checkbox label="标签" style="margin-bottom:10px"></el-checkbox>
          <el-checkbox label="cron表达式" style="margin-bottom:10px"></el-checkbox>
          <el-checkbox label="下次执行时间" style="margin-bottom:10px"></el-checkbox>
          <el-checkbox label="执行方式" style="margin-bottom:10px"></el-checkbox>
          <el-checkbox label="状态"></el-checkbox>
        </el-checkbox-group>
        <el-button slot="reference" type="text" style="margin-left: 10px;color: rgba(0,0,0,.75);font-size: 16px">
          <el-tooltip class="item" effect="dark" content="列设置" placement="top">
            <i class="el-icon-s-tools"></i>
          </el-tooltip>
        </el-button>
      </el-popover>
    </div>

    <!-- 列表 -->
    <el-table
      :data="tasks"
      tooltip-effect="dark"
      border
      :size="tableSize"
      v-loading="tableLoading"
      @selection-change="handleSelectionChange"
      style="width: 100%;">
      <el-table-column type="expand">
        <template slot-scope="scope">
          <el-form label-position="left" inline class="demo-table-expand">
            <el-form-item label="任务创建时间:">
              {{scope.row.created | formatTime}} <br>
            </el-form-item>
            <el-form-item label="任务类型：">
              {{scope.row.level | formatLevel}} <br>
            </el-form-item>
            <el-form-item label="单实例运行：">
               {{scope.row.multi | formatMulti}} <br>
            </el-form-item>
            <el-form-item label="超时时间：">
              {{scope.row.timeout | formatTimeout}} <br>
            </el-form-item>
            <el-form-item label="重试次数:">
              {{scope.row.retry_times}} <br>
            </el-form-item>
            <el-form-item label="重试间隔：">
              {{scope.row.retry_interval | formatRetryTimesInterval}}
            </el-form-item> <br>
            <el-form-item label="任务节点：">
              <span v-if="!(scope.row.hosts && scope.row.hosts.length)">无节点</span>
              <div v-else v-for="item in scope.row.hosts" :key="item.host_id">
                {{item.alias}} - {{item.name}}:{{item.port}} <br>
              </div>
            </el-form-item> <br>
            <el-form-item label="命令：" style="width: 100%">
              {{scope.row.command}}
            </el-form-item> <br>
            <el-form-item label="备注：" style="width: 100%">
              {{scope.row.remark ? scope.row.remark : '--'}}
            </el-form-item>
          </el-form>
        </template>
      </el-table-column>
      <el-table-column
        type="selection"
        width="50"
        align="center">
      </el-table-column>
      <el-table-column
        prop="id"
        label="任务ID"
        align="center"
        v-if="checkList.includes('任务ID')">
      </el-table-column>
      <el-table-column
        prop="name"
        label="任务名称"
        width="150"
        align="center"
        v-if="checkList.includes('任务名称')">
      </el-table-column>
      <el-table-column
        prop="tag"
        label="标签"
        align="center"
        v-if="checkList.includes('标签')">
      </el-table-column>
      <el-table-column
        prop="spec"
        label="cron表达式"
        width="120"
        align="center"
        v-if="checkList.includes('cron表达式')">
      </el-table-column>
      <el-table-column
        label="下次执行时间"
        width="160"
        align="center"
        v-if="checkList.includes('下次执行时间')">
        <template slot-scope="scope">
          {{scope.row.next_run_time | formatTime}}
        </template>
      </el-table-column>
      <el-table-column
        prop="protocol"
        :formatter="formatProtocol"
        label="执行方式"
        align="center"
        v-if="checkList.includes('执行方式')">
      </el-table-column>
      <el-table-column
        label="状态"
        v-if="this.isAdmin && checkList.includes('状态')"
        align="center">
          <template slot-scope="scope">
            <el-switch
              v-if="scope.row.level === 1"
              v-model="scope.row.status"
              :active-value="1"
              :inactive-vlaue="0"
              active-color="#13ce66"
              @change="changeStatus(scope.row)"
              inactive-color="#ff4949">
            </el-switch>
          </template>
      </el-table-column>
      <el-table-column
        label="状态"
        v-if="!this.isAdmin && checkList.includes('状态')"
        align="center">
        <template slot-scope="scope">
          <el-switch
            v-if="scope.row.level === 1"
            v-model="scope.row.status"
            :active-value="1"
            :inactive-vlaue="0"
            active-color="#13ce66"
            :disabled="true"
            inactive-color="#ff4949">
          </el-switch>
        </template>
      </el-table-column>
      <el-table-column label="操作" min-width="400" v-if="this.isAdmin" align="center">
        <template slot-scope="scope">
          <el-button :size="tableSize" type="primary" @click="toEdit(scope.row)">编辑</el-button>
          <el-button :size="tableSize" type="success" @click="runTask(scope.row)">手动执行</el-button>
          <el-button :size="tableSize" type="info" @click="jumpToLog(scope.row)">查看日志</el-button>
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
        :total="taskTotal"
        :page-size="20"
        @size-change="changePageSize"
        @current-change="changePage"
        @prev-click="changePage"
        @next-click="changePage">
      </el-pagination>
    </div>
  </div>

  <div style="overflow:scroll">
    <el-drawer
      :visible="editDrawerVisible"
      direction="rtl"
      destroy-on-close
      :with-header="false"
      @close="editDrawerVisible = false"
      size="70%">
      <task-edit :taskid="taskid" @complete="completeEdit" />
    </el-drawer>
  </div>
</el-container>
</template>

<script>
import taskService from '../../api/task'
import taskEdit from './edit'

export default {
  name: 'task-list',
  data () {
    return {
      tableSize: 'mini',
      checkList: ['任务ID', '任务名称', '标签', 'cron表达式', '下次执行时间', '执行方式', '状态'],
      tasks: [],
      hosts: [],
      taskTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        id: '',
        protocol: '',
        name: '',
        tag: '',
        host_id: '',
        status: ''
      },
      tableLoading: false,
      moreParams: false,
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
      ],
      multipleSelection: [],
      editDrawerVisible: false,
      taskid: ''
    }
  },
  components: {
    taskEdit
  },
  created () {
    const hostId = this.$route.query.host_id
    if (hostId) {
      this.searchParams.host_id = hostId
    }

    this.search()
  },
  filters: {
    formatLevel (value) {
      if (value === 1) {
        return '主任务'
      }
      return '子任务'
    },
    formatTimeout (value) {
      if (value > 0) {
        return value + '秒'
      }
      return '不限制'
    },
    formatRetryTimesInterval (value) {
      if (value > 0) {
        return value + '秒'
      }
      return '系统默认'
    },
    formatMulti (value) {
      if (value > 0) {
        return '否'
      }
      return '是'
    }
  },
  computed: {
    ifMulSelected () {
      return this.multipleSelection.length > 0
    }
  },
  methods: {
    changeStatus (item) {
      if (item.status) {
        taskService.enable(item.id)
      } else {
        taskService.disable(item.id)
      }
    },
    formatProtocol (row, col) {
      if (row[col.property] === 2) {
        return 'shell'
      }
      if (row.http_method === 1) {
        return 'http-get'
      }
      return 'http-post'
    },
    changePage (page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize (pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    handleSelectionChange (val) {
      this.multipleSelection = val
    },
    search (callback = null) {
      this.tableLoading = true
      taskService.list(this.searchParams, (tasks, hosts) => {
        this.tableLoading = false
        this.tasks = tasks.data
        this.taskTotal = tasks.total
        this.hosts = hosts
        if (callback) {
          callback()
        }
      })
    },
    // 重置
    reset () {
      this.searchParams = {
        page_size: 20,
        page: 1,
        id: '',
        protocol: '',
        name: '',
        tag: '',
        host_id: '',
        status: ''
      }
      this.search()
    },
    // 批量手动
    batchTrigger () {
      this.multipleSelection.forEach(element => {
        this.runTask(element)
      })
    },
    // 批量删除
    batchDelete () {
      this.multipleSelection.forEach(element => {
        this.remove(element)
      })
    },
    runTask (item) {
      this.$appConfirm(() => {
        taskService.run(item.id, () => {
          this.$message.success('任务已开始执行')
        })
      }, true)
    },
    remove (item) {
      this.$appConfirm(() => {
        taskService.remove(item.id, () => {
          this.refresh()
        })
      })
    },
    jumpToLog (item) {
      this.$router.push(`/task-log?task_id=${item.id}`)
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    toEdit (item) {
      if (item === null) {
        this.taskid = undefined
      } else {
        this.taskid = item.id
      }
      this.editDrawerVisible = true
    },
    completeEdit () {
      console.log('fuck')
      this.editDrawerVisible = false
      this.search()
    }
  }
}
</script>
<style scoped>
  .demo-table-expand {
    font-size: 0;
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
  .el-drawer .rtl {
    overflow: scroll;
  }
</style>
