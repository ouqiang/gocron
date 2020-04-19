<template>
  <el-container>
    <div style="width:100%;padding:40px">
      <el-form :inline="true" style="padding: 20px 40px 5px 40px;background: #fff;border-radius: 5px">
        <el-row>
          <el-form-item label="节点ID">
            <el-input size="small" v-model.trim="searchParams.id"></el-input>
          </el-form-item>
          <el-form-item label="主机名">
            <el-input size="small" v-model.trim="searchParams.name"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button size="small" type="primary" @click="search()">搜索</el-button>
          </el-form-item>
          <el-form-item>
            <el-button size="small" type="info" @click="reset()">重置</el-button>
          </el-form-item>
        </el-row>
      </el-form>

      <!-- 功能区 -->
      <div style="margin-top: 20px; padding: 20px 40px;background: #fff;text-align: right;border-radius: 5px;">
        <h1 style="float: left;height:42px;line-height: 42px;margin: 0;display: inline-block;">
          任务节点
        </h1>
        <el-button size="small" type="primary" v-if="this.$store.getters.user.isAdmin"  @click="toEdit(null)">新增</el-button>
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

      <!-- 表格 -->
      <el-table
        :data="hosts"
        tooltip-effect="dark"
        :size="tableSize"
        v-loading="tableLoading"
        border>
        <el-table-column
          prop="id"
          label="节点ID"
          align="center">
        </el-table-column>
        <el-table-column
          prop="alias"
          label="节点名称"
          align="center"
          v-if="checkList.includes('节点名称')">
        </el-table-column>
        <el-table-column
          prop="name"
          label="主机名"
          align="center"
          v-if="checkList.includes('主机名')">
        </el-table-column>
        <el-table-column
          prop="port"
          label="端口"
          align="center"
          v-if="checkList.includes('端口')">
        </el-table-column>
        <el-table-column
          prop="remark"
          label="备注"
          align="center"
          v-if="checkList.includes('备注')">
        </el-table-column>
        <el-table-column label="操作" width="450" v-if="this.isAdmin" align="center">
          <template slot-scope="scope">
            <el-button :size="tableSize" type="success" @click="toTasks(scope.row)">查看任务</el-button>
            <el-button :size="tableSize" type="primary" @click="toEdit(scope.row)">编辑</el-button>
            <el-button :size="tableSize" type="info" @click="ping(scope.row)">测试连接</el-button>
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
          :total="hostTotal"
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
        size="30%">
        <edit-dialog :hostid="hostid" @complete="completeEdit"></edit-dialog>
      </el-drawer>
    </div>
  </el-container>
</template>

<script>
import hostService from '../../api/host'
import editDialog from './edit'

export default {
  name: 'host-list',
  components: {
    editDialog
  },
  data () {
    return {
      hosts: [],
      hostid: undefined,
      hostTotal: 0,
      tableSize: 'mini',
      tableLoading: false,
      editDrawerVisible: false,
      checkListOrigin: ['节点名称', '主机名', '端口', '备注'],
      checkList: ['节点名称', '主机名', '端口', '备注'],
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
      this.tableLoading = true
      hostService.list(this.searchParams, (data) => {
        this.tableLoading = false
        this.hosts = data.data
        this.hostTotal = data.total
        if (callback) {
          callback()
        }
      })
    },
    reset () {
      this.searchParams = {
        page_size: 20,
        page: 1,
        id: '',
        name: '',
        alias: ''
      }
      this.search()
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
      if (item === null) {
        this.hostid = undefined
      } else {
        this.hostid = item.id
      }
      this.editDrawerVisible = true
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
    },
    completeEdit () {
      this.editDrawerVisible = false
      this.search()
    }
  }
}
</script>
