<template>
    <el-card>
      <template #header>
        <div class="card-header">
          <strong>进程管理</strong>
        </div>
      </template>
      <el-form :inline="true" >
        <el-row>
          <el-form-item label="进程Id">
            <el-input v-model.trim="searchParams.id"></el-input>
          </el-form-item>
          <el-form-item label="名称">
            <el-input v-model.trim="searchParams.name"></el-input>
          </el-form-item>
          <el-form-item label="命令">
            <el-input v-model.trim="searchParams.command"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="search()">搜索</el-button>
            <el-button type="default" @click="resetSearch()">重置</el-button>
            <el-button type="success" @click="toEdit(null)" v-if="this.$store.getters.user.isAdmin">新增</el-button>
          </el-form-item>
        </el-row>
      </el-form>
      <el-pagination
        style="margin: 5px 0"
        background
        layout="prev, pager, next, sizes, total"
        :total="processTotal"
        :page-size="20"
        @size-change="changePageSize"
        @current-change="changePage"
        @prev-click="changePage"
        @next-click="changePage">
      </el-pagination>
      <el-table
        :data="processList"
        tooltip-effect="dark"
        row-key="id"
        border
        style="width: 100%">
        <el-table-column type="expand">
          <template #default="scope">
            <el-container>
              <el-table :data="scope.row.workers" border>
                <el-table-column label="节点" width="240">
                  <template #default="scope">
                    <el-tag v-if="scope.row.host_id === 0" type="info">待定</el-tag>
                    <el-tag v-else :title="hosts[scope.row.host_id]">{{ hosts[scope.row.host_id] }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="pid" label="PID" align="center" width="100"/>
                <el-table-column label="状态" align="center" width="120">
                  <template #default="scope">
                    <el-tag v-if="scope.row.state === 0" type="info">
                      <el-icon>
                        <VideoPause/>
                      </el-icon>
                      挂起
                    </el-tag>
                    <el-tag v-else-if="scope.row.state === 1" type="success">
                      <el-icon>
                        <Loading/>
                      </el-icon>
                      启动中
                    </el-tag>
                    <el-tag v-else-if="scope.row.state === 2" type="success">
                      <el-icon>
                        <Loading/>
                      </el-icon>
                      运行中
                    </el-tag>
                    <el-tag v-else-if="scope.row.state === 3">
                      <el-icon>
                        <QuestionFilled/>
                      </el-icon>
                      停止中
                    </el-tag>
                    <el-tag v-else-if="scope.row.state === 4" type="info">
                      <el-icon>
                        <WarningFilled/>
                      </el-icon>
                      已停止
                    </el-tag>
                    <el-tag v-else-if="scope.row.state === 5" type="danger">
                      <el-icon>
                        <CircleCloseFilled/>
                      </el-icon>
                      异常退出
                    </el-tag>
                    <el-tag v-else>未知状态 {{ scope.row.state }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="启动时间">
                  <template #default="scope">
                    {{ formatDatetime(scope.row.start_at) }}
                  </template>
                </el-table-column>
              </el-table>
            </el-container>
          </template>
        </el-table-column>
        <el-table-column  prop="id" label="ID" width="100" align="center"/>
        <el-table-column  prop="name" label="名称" />
        <el-table-column  prop="command" label="命令" />
        <el-table-column  prop="num_proc" label="进程数" width="80" align="center" />
        <el-table-column  prop="tag" label="标签" />
        <el-table-column  label="状态">
          <template #default="scope">
            <el-tag v-if="scope.row.status === 1" type="success">已启动</el-tag>
            <el-tag v-else-if="scope.row.status === 2" type="info">已停止</el-tag>
            <el-tag v-else>未知状态 {{scope.row.status}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column  label="是否启用" width="100" align="center">
          <template #default="scope">
            <el-switch
                v-model="scope.row.enable"
                :active-value="true"
                :inactive-value="false"
                @change="changeStatus(scope.row)"
                active-color="#13ce66"
                inactive-color="#ff4949">
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="primary" @click="toEdit(scope.row.id)">编辑</el-button>
            <el-button type="success" v-if="scope.row.status !== 1" @click="start(scope.row.id)">启动</el-button>
            <el-button type="danger" v-if="scope.row.status === 1" @click="stop(scope.row.id)">停止</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
</template>

<script>
import processService from '../../api/process'
import hostService from '../../api/host'
import format from '@/utils/format'
import {SUCCESS_CODE} from '@/utils/httpClient'

export default {
  name: 'process-list',
  data () {
    return {
      searchParams: {},
      processTotal: 0,
      processList: [],
      hosts: {},
    }
  },
  created () {
    let _this = this
    this.search()

    hostService.all({}, function (data) {
      data.forEach(host => {
        _this.hosts[host.id] = host.alias + ' - ' + host.name + ':' + host.port
      })
      console.log('hosts',data)
    })
  },
  methods: {
    formatDatetime: format.formatDatetime,
    search () {
      let _this = this
      processService.list(this.searchParams, (resp) => {
        _this.processTotal = resp.total
        _this.processList = resp.data
      })
    },
    changePage (page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize (pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    resetSearch () {
      this.searchParams = {}
      this.search()
    },
    changeStatus(row) {
      let _this = this
      if (row.enable) {
        processService.enable(row.id, function (resp) {
          if (resp.code === SUCCESS_CODE) {
            _this.$message.success({message: "启用成功"})
          } else {
            _this.processList.forEach(process => {
              if (process.id === row.id) {
                process.enable = false
              }
            })
            _this.$message.error({message: resp.message})
          }
        })
      } else {
        processService.disable(row.id,function(resp){
          if (resp.code === SUCCESS_CODE) {
            _this.$message.success({message: "禁用成功"})
          } else {
            _this.processList.forEach(process => {
              if (process.id === row.id) {
                process.enable = true
              }
            })
            _this.$message.error({message: resp.message})
          }
        })
      }
    },
    start (id) {
      processService.start(id, (data) => {
        console.log(data)
        this.$message.success('程序已启动')
        this.search()
      })
    },
    stop (id) {
      processService.stop(id, (data) => {
        console.log(data)
        this.$message.success('程序已停止')
      })
    },
    restart (id) {
      processService.restart(id, (data) => {
        console.log(data)
      })
    },
    toEdit (id) {
      let path = ''
      if (id === null) {
        path = '/process/create'
      } else {
        path = `/process/edit/${id}`
      }
      this.$router.push(path)
    },
    refresh () {

    }
  }
}
</script>
<style scoped>

</style>
