<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="Path">
          <el-input v-model.trim="params.path" clearable placeholder="Path" @clear="search" />
        </el-form-item>
        <el-form-item label="Category">
          <el-input v-model.trim="params.category" clearable placeholder="Category" @clear="search" />
        </el-form-item>
        <el-form-item label="Method">
          <el-select v-model.trim="params.method" clearable placeholder="Method" @change="search" @clear="search">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="PATCH" value="PATCH" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="Creator">
          <el-input v-model.trim="params.creator" clearable placeholder="Creator" @clear="search" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">Search</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">Create</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger"
            @click="batchDelete">Batch Delete</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%"
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="path" label="Path" />
        <el-table-column show-overflow-tooltip sortable prop="category" label="Category" />
        <el-table-column show-overflow-tooltip sortable prop="method" label="Method" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.method | methodTagFilter" disable-transitions>
              {{ scope.row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="creator" label="Creator" />
        <el-table-column show-overflow-tooltip sortable prop="desc" label="Description" />
        <el-table-column fixed="right" label="Action" align="center" width="120">
          <template slot-scope="scope">
            <el-tooltip content="Edit" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip class="delete-popover" content="Delete" effect="dark" placement="top">
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

      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="120px">
          <el-form-item label="Path" prop="path">
            <el-input v-model.trim="dialogFormData.path" placeholder="Path" />
          </el-form-item>
          <el-form-item label="Category" prop="category">
            <el-input v-model.trim="dialogFormData.category" placeholder="Category" />
          </el-form-item>
          <el-form-item label="Method" prop="method">
            <el-select v-model.trim="dialogFormData.method" placeholder="Select one">
              <el-option label="GET" value="GET" />
              <el-option label="POST" value="POST" />
              <el-option label="PUT" value="PUT" />
              <el-option label="PATCH" value="PATCH" />
              <el-option label="DELETE" value="DELETE" />
            </el-select>
          </el-form-item>
          <el-form-item label="Description" prop="desc">
            <el-input v-model.trim="dialogFormData.desc" type="textarea" placeholder="Description" show-word-limit
              maxlength="100" />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">Cancel</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">Save</el-button>
        </div>
      </el-dialog>

    </el-card>
  </div>
</template>

<script>
import { batchDeleteApiByIds, createApi, getApis, updateApiById } from '@/api/system/api'

export default {
  name: 'Api',
  filters: {
    methodTagFilter(val) {
      if (val === 'GET') {
        return ''
      } else if (val === 'POST') {
        return 'success'
      } else if (val === 'PUT') {
        return 'info'
      } else if (val === 'PATCH') {
        return 'warning'
      } else if (val === 'DELETE') {
        return 'danger'
      } else {
        return 'info'
      }
    }
  },
  data() {
    return {
      params: {
        path: '',
        method: '',
        category: '',
        creator: '',
        pageNum: 1,
        pageSize: 10
      },
      tableData: [],
      total: 0,
      loading: false,
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        path: '',
        category: '',
        method: '',
        desc: ''
      },
      dialogFormRules: {
        path: [
          { required: true, message: 'Please enter access path', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
        ],
        category: [
          { required: true, message: 'Please enter category', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        method: [
          { required: true, message: 'Please choose request method', trigger: 'change' }
        ],
        desc: [
          { required: false, message: 'Please enter description', trigger: 'blur' },
          { min: 0, max: 100, message: 'Maximum 100 characters', trigger: 'blur' }
        ]
      },
      popoverVisible: false,
      multipleSelection: []
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getApis(this.params)
        this.tableData = data.apis
        this.total = data.total
      } finally {
        this.loading = false
      }
    },
    create() {
      this.dialogFormTitle = 'Create'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },
    update(row) {
      this.dialogFormData.ID = row.ID
      this.dialogFormData.path = row.path
      this.dialogFormData.category = row.category
      this.dialogFormData.method = row.method
      this.dialogFormData.desc = row.desc

      this.dialogFormTitle = 'Edit'
      this.dialogType = 'update'
      this.dialogFormVisible = true
    },
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          let msg = ''
          this.submitLoading = true
          try {
            if (this.dialogType === 'create') {
              const { message } = await createApi(this.dialogFormData)
              msg = message
            } else {
              const { message } = await updateApiById(this.dialogFormData.ID, this.dialogFormData)
              msg = message
            }
          } finally {
            this.submitLoading = false
          }

          this.resetForm()
          this.getTableData()
          this.$message({
            showClose: true,
            message: msg,
            type: 'success'
          })
        } else {
          this.$message({
            showClose: true,
            message: 'Please check your data',
            type: 'error'
          })
          return false
        }
      })
    },
    cancelForm() {
      this.resetForm()
    },

    resetForm() {
      this.dialogFormVisible = false
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {
        path: '',
        category: '',
        method: '',
        desc: ''
      }
    },
    batchDelete() {
      this.$confirm('This cannot be undone. Do you want to continue?', 'Delete', {
        confirmButtonText: 'Yes',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const apiIds = []
        this.multipleSelection.forEach(x => {
          apiIds.push(x.ID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteApiByIds({ apiIds: apiIds })
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
        const { message } = await batchDeleteApiByIds({ apiIds: [Id] })
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
