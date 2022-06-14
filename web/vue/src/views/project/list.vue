<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <strong>项目管理</strong>
      </div>
    </template>
    <el-form :inline="true">
      <el-row>
        <el-form-item label="名称">
          <el-input v-model.trim="searchForm.name"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button type="default" @click="resetSearch">重置</el-button>
          <el-button type="success" @click="toEdit(null)" v-if="this.$store.getters.user.isAdmin">新增</el-button>
        </el-form-item>
      </el-row>
    </el-form>
    <el-pagination
        style="margin: 5px 0"
        background
        layout="prev, pager, next, sizes, total"
        :total="projectTotal"
        :page-size="searchForm.page_size"
        @size-change="changePageSize"
        @current-change="changePage"
        @prev-click="changePage"
        @next-click="changePage">
    </el-pagination>

    <el-table
        :data="projects"
        tooltip-effect="dark"
        id="project-list"
        border
        style="width: 100%">
      <el-table-column prop="id" label="ID" width="100" align="center"/>
      <el-table-column prop="name" label="名称"/>
      <el-table-column prop="code" label="编码"/>
      <el-table-column prop="remark" label="备注" :show-overflow-tooltip="true"/>
      <el-table-column prop="created_at" :formatter="formatDatetime" label="创建时间"/>
      <el-table-column prop="updated_at" :formatter="formatDatetime" label="更新时间"/>
      <el-table-column label="操作">
        <template #default="scope">
          <el-button type="primary" @click="toEdit(scope.row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>



  <el-dialog v-model="isEdit" :title="editForm.id ?'编辑项目':'新增项目'" :close-on-click-modal="false" :destroy-on-close="true">
    <el-form ref="projectForm" :model="editForm" :rules="rules" label-width="120px" >
      <el-form-item label="项目名称" prop="name" required>
        <el-input v-model="editForm.name"/>
      </el-form-item>
      <el-form-item label="项目编码" prop="code" required>
        <el-input v-model="editForm.code"/>
      </el-form-item>
      <el-form-item label="主机" prop="host_ids">
        <el-select v-model="editForm.host_ids" placeholder="请选择关联节点" multiple style="width:100%">
          <el-option
              v-for="host in hosts"
              :key="host.id"
              :label="host.alias + ' - ' + host.name"
              :value="host.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="备注">
        <el-input
            type="textarea"
            :rows="5"
            width="100"
            v-model="editForm.remark">
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="projectSubmit">提交</el-button>
        <el-button @click="isEdit = false">取消</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script>
import projectService from '@/api/project'
import format from '@/utils/format'
import hostService from "@/api/host";

export default {
  name: 'project-list',
  data() {
    return {
      isEdit: false,
      hosts:[],
      searchForm: {
        page_size: 20,
        page: 1
      },
      projectTotal: 0,
      projects: [],
      editForm: {}
    }
  },
  created() {
    let _this = this
    hostService.all({}, function (data) {
      _this.hosts = data
    })
    this.search()
  },
  methods: {
    search() {
      let _this = this
      projectService.list(this.searchForm, function (data) {
        console.log(data)
        _this.projects = data.projects
        _this.projectTotal = data.total
      })
    },
    formatDatetime: function (row, col) {
      return format.formatDatetime(row[col.property])
    },
    changePage(page) {
      this.searchForm.page = page
      this.search()
    },
    changePageSize(pageSize) {
      this.searchForm.page_size = pageSize
      this.search()
    },
    resetSearch() {

    },
    toEdit(row) {
      if (row !== null) {
        let project = Object.assign({}, row)
        project.host_ids = []
        project.hosts && project.hosts.forEach(h => {
          project.host_ids.push(h.id)
        })
        this.editForm = project
      } else {
        this.editForm = {}
      }
      this.isEdit = true
    },
    projectSubmit() {
      let _this = this
      this.$refs['projectForm'].validate(valid => {
        if (!valid) {
          return false
        }
        this.editForm.host_ids = this.editForm.host_ids.join(',')
        projectService.store(this.editForm, function () {
          _this.$message.success('操作成功')
          _this.isEdit = false
          _this.search()
        })
      })
    }
  }
}
</script>
<style>
#project-list .el-popper {
  max-width: 30%
}
</style>