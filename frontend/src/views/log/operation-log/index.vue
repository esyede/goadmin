<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="Username">
          <el-input v-model.trim="params.username" clearable placeholder="Username" @clear="search" />
        </el-form-item>
        <el-form-item label="IP">
          <el-input v-model.trim="params.ip" clearable placeholder="IP" @clear="search" />
        </el-form-item>
        <el-form-item label="Path">
          <el-input v-model.trim="params.path" clearable placeholder="Path" @clear="search" />
        </el-form-item>
        <el-form-item label="Status">
          <el-input v-model.trim="params.status" clearable placeholder="Status" @clear="search" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">Search</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger"
            @click="batchDelete">Batch Delete</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%"
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="username" label="Username" />
        <el-table-column show-overflow-tooltip sortable prop="ip" label="IP" />
        <el-table-column show-overflow-tooltip sortable prop="path" label="Path" />
        <el-table-column show-overflow-tooltip sortable prop="status" label="Status" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status | statusTagFilter" disable-transitions>{{ scope.row.status
            }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="startTime" label="Starts at">
          <template slot-scope="scope">
            {{ parseGoTime(scope.row.startTime) }}
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="timeCost" label="Duration (ms)" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.timeCost | timeCostTagFilter" disable-transitions>{{ scope.row.timeCost
            }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="desc" label="Description" />
        <el-table-column fixed="right" label="Action" align="center" width="80">
          <template slot-scope="scope">
            <el-tooltip content="Delete" effect="dark" placement="top">
              <el-popconfirm title="Delete this data?" @onConfirm="singleDelete(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" circle type="danger" />
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination :current-page="params.pageNum" :page-size="params.pageSize" :total="total"
        :page-sizes="[1, 5, 10, 30]" layout="total, prev, pager, next, sizes" background
        style="margin-top: 10px;float:right;margin-bottom: 10px;" @size-change="handleSizeChange"
        @current-change="handleCurrentChange" />
    </el-card>
  </div>
</template>

<script>
import { batchDeleteOperationLogByIds, getOperationLogs } from '@/api/log/operationLog'
import { parseGoTime } from '@/utils/index'

export default {
  name: 'OperationLog',
  filters: {
    statusTagFilter(val) {
      if (val === 200) {
        return 'success'
      } else if (val === 400) {
        return 'warning'
      } else if (val === 401) {
        return 'danger'
      } else if (val === 403) {
        return 'danger'
      } else if (val === 500) {
        return 'danger'
      } else {
        return 'info'
      }
    },
    timeCostTagFilter(val) {
      if (val <= 200) {
        return 'success'
      } else if (val > 200 && val <= 1000) {
        return ''
      } else if (val > 1000 && val <= 2000) {
        return 'warning'
      } else {
        return 'danger'
      }
    }
  },
  data() {
    return {
      params: {
        username: '',
        ip: '',
        path: '',
        status: '',
        pageNum: 1,
        pageSize: 10
      },
      tableData: [],
      total: 0,
      loading: false,
      popoverVisible: false,
      multipleSelection: []
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    parseGoTime,
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getOperationLogs(this.params)
        this.tableData = data.logs
        this.total = data.total
      } finally {
        this.loading = false
      }
    },
    batchDelete() {
      this.$confirm('This cannot be undone. Do you want to continue?', 'Delete', {
        confirmButtonText: 'Yes',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const operationLogIds = []
        this.multipleSelection.forEach(x => {
          operationLogIds.push(x.ID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteOperationLogByIds({ operationLogIds: operationLogIds })
          msg = message
        } finally {
          this.loading = false
        }

        this.getTableData()
        this.$message({
          showClose: true,
          message: msg,
          type: 'success'
        })
      }).catch(() => {
        this.$message({
          showClose: true,
          type: 'info',
          message: 'Restore'
        })
      })
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    async singleDelete(Id) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteOperationLogByIds({ operationLogIds: [Id] })
        msg = message
      } finally {
        this.loading = false
      }

      this.getTableData()
      this.$message({
        showClose: true,
        message: msg,
        type: 'success'
      })
    },
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getTableData()
    }
  }
}
</script>

<style scoped>
.container-card {
  margin: 10px;
}

.delete-popover {
  margin-left: 10px;
}
</style>
