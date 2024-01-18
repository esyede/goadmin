<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="Role">
          <el-input v-model.trim="params.name" clearable placeholder="Role name" @clear="search" />
        </el-form-item>
        <el-form-item label="Keyword">
          <el-input v-model.trim="params.keyword" clearable placeholder="Keyword" @clear="search" />
        </el-form-item>
        <el-form-item label="Status">
          <el-select v-model.trim="params.status" clearable placeholder="Status" @change="search" @clear="search">
            <el-option label="Enabled" :value="1" />
            <el-option label="Disabled" :value="2" />
          </el-select>
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
        <el-table-column show-overflow-tooltip sortable prop="name" label="Name" />
        <el-table-column show-overflow-tooltip sortable prop="keyword" label="Keyword" />
        <el-table-column show-overflow-tooltip sortable prop="sort" label="Sort" />
        <el-table-column show-overflow-tooltip sortable prop="status" label="Status" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status === 1 ? 'success' : 'danger'" disable-transitions>
              {{ scope.row.status === 1 ? 'Enabled' : 'Disabled' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="creator" label="Creator" />
        <el-table-column show-overflow-tooltip sortable prop="desc" label="Description" />
        <el-table-column fixed="right" label="Action" align="center" width="140">
          <template slot-scope="scope">
            <el-tooltip content="Edit" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" circle type="primary" @click="update(scope.row)" />
            </el-tooltip>
            <el-tooltip content="Permission" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-key" circle type="warning" @click="updatePermission(scope.row.ID)" />
            </el-tooltip>
            <el-tooltip content="Delete" effect="dark" placement="top">
              <el-popconfirm style="margin-left:10px" title="Delete this data?" @onConfirm="singleDelete(scope.row.ID)">
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

      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="580px">
        <el-form ref="dialogForm" :inline="true" size="small" :model="dialogFormData" :rules="dialogFormRules"
          label-width="100px">
          <el-form-item label="Role Name" prop="name">
            <el-input v-model.trim="dialogFormData.name" placeholder="Role name" style="width: 420px" />
          </el-form-item>
          <el-form-item label="Keyword" prop="keyword">
            <el-input v-model.trim="dialogFormData.keyword" placeholder="Keyword" style="width: 420px" />
          </el-form-item>
          <el-form-item label="Status" prop="status">
            <el-select v-model.trim="dialogFormData.status" placeholder="Status" style="width: 180px">
              <el-option label="Enabled" :value="1" />
              <el-option label="Disabled" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item label="Sort (1 highest)" prop="sort">
            <el-input-number v-model.number="dialogFormData.sort" controls-position="right" :min="1" :max="999" />
          </el-form-item>
          <el-form-item label="Description" prop="desc">
            <el-input v-model.trim="dialogFormData.desc" style="width: 420px" type="textarea" placeholder="Description"
              show-word-limit maxlength="100" />
          </el-form-item>
        </el-form>
        <div slot="footer">
          <el-button size="mini" @click="cancelForm()">Cancel</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">Save</el-button>
        </div>
      </el-dialog>

      <el-dialog title="Permission" :visible.sync="permsDialogVisible" width="580px" custom-class="perms-dialog">
        <el-tabs>
          <el-tab-pane>
            <span slot="label"><svg-icon icon-class="menu1" class-name="role-menu" />Role</span>
            <el-tree ref="roleMenuTree" v-loading="menuTreeLoading" :props="{ children: 'children', label: 'title' }"
              :data="menuTree" show-checkbox node-key="ID" check-strictly
              :default-checked-keys="defaultCheckedRoleMenu" />

          </el-tab-pane>

          <el-tab-pane>
            <span slot="label"><svg-icon icon-class="api1" class-name="role-menu" />Permission</span>
            <el-tree ref="roleApiTree" v-loading="apiTreeLoading" :props="{ children: 'children', label: 'desc' }"
              :data="apiTree" show-checkbox node-key="ID" :default-checked-keys="defaultCheckedRoleApi" />

          </el-tab-pane>
        </el-tabs>
        <div slot="footer">
          <el-button size="mini" :loading="permissionLoading" @click="cancelPermissionForm()">Cancel</el-button>
          <el-button size="mini" type="primary" @click="submitPermissionForm()">Save</el-button>
        </div>
      </el-dialog>

    </el-card>
  </div>
</template>

<script>
import { getApiTree } from '@/api/system/api'
import { getMenuTree } from '@/api/system/menu'
import { batchDeleteRoleByIds, createRole, getRoleApisById, getRoleMenusById, getRoles, updateRoleApisById, updateRoleById, updateRoleMenusById } from '@/api/system/role'

export default {
  name: 'Role',
  data() {
    return {
      params: {
        name: '',
        keyword: '',
        status: '',
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
        name: '',
        keyword: '',
        status: 1,
        sort: 999,
        desc: ''
      },
      dialogFormRules: {
        name: [
          { required: true, message: 'Please enter name', trigger: 'blur' },
          { min: 1, max: 20, message: 'Must be between 1 - 20 characters', trigger: 'blur' }
        ],
        keyword: [
          { required: true, message: 'Please enter keyword', trigger: 'blur' },
          { min: 1, max: 20, message: 'Must be between 1 - 20 characters', trigger: 'blur' }
        ],
        status: [
          { required: true, message: 'Please choose status', trigger: 'change' }
        ],
        desc: [
          { required: false, message: 'Please enter description', trigger: 'blur' },
          { min: 0, max: 100, message: 'Must be less than 100 characters', trigger: 'blur' }
        ]
      },
      popoverVisible: false,
      multipleSelection: [],
      permsDialogVisible: false,
      permissionLoading: false,
      menuTree: [],
      defaultCheckedRoleMenu: [],
      apiTree: [],
      defaultCheckedRoleApi: [],
      roleId: 0
    }
  },
  created() {
    this.getTableData()
    this.getMenuTree()
    this.getApiTree()
  },
  methods: {
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getRoles(this.params)
        this.tableData = data.roles
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
      this.dialogFormData.name = row.name
      this.dialogFormData.keyword = row.keyword
      this.dialogFormData.sort = row.sort
      this.dialogFormData.status = row.status
      this.dialogFormData.desc = row.desc
      this.dialogFormTitle = 'Edit'
      this.dialogType = 'update'
      this.dialogFormVisible = true
    },
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.submitLoading = true
          let msg = ''
          try {
            if (this.dialogType === 'create') {
              const { message } = await createRole(this.dialogFormData)
              msg = message
            } else {
              const { message } = await updateRoleById(this.dialogFormData.ID, this.dialogFormData)
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
        name: '',
        keyword: '',
        status: 1,
        sort: 999,
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
        const roleIds = []
        this.multipleSelection.forEach(x => {
          roleIds.push(x.ID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteRoleByIds({ roleIds: roleIds })
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
          type: 'info',
          message: 'Restore'
        })
      })
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    async singleDelete(id) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteRoleByIds({ roleIds: [id] })
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
    async updatePermission(roleId) {
      this.roleId = roleId
      this.permsDialogVisible = true
      this.getMenuTree()
      this.getApiTree()
      this.getRoleMenusById(roleId)
      this.getRoleApisById(roleId)
    },
    async getMenuTree() {
      this.menuTreeLoading = true
      try {
        const { data } = await getMenuTree()
        this.menuTree = data.menuTree
      } finally {
        this.menuTreeLoading = false
      }
    },
    async getApiTree() {
      this.apiTreeLoading = true
      try {
        const { data } = await getApiTree()
        this.apiTree = data.apiTree
      } finally {
        this.apiTreeLoading = false
      }
    },
    async getRoleMenusById(roleId) {
      this.permissionLoading = true
      let rseData = []
      try {
        const { data } = await getRoleMenusById(roleId)
        rseData = data
      } finally {
        this.permissionLoading = false
      }

      const menus = rseData.menus
      const ids = []
      menus.forEach(x => { ids.push(x.ID) })
      this.defaultCheckedRoleMenu = ids
      this.$refs.roleMenuTree.setCheckedKeys(this.defaultCheckedRoleMenu)
    },
    async getRoleApisById(roleId) {
      this.permissionLoading = true
      let resData = []
      try {
        const { data } = await getRoleApisById(roleId)
        resData = data
      } finally {
        this.permissionLoading = false
      }

      const apis = resData.apis
      const ids = []
      apis.forEach(x => { ids.push(x.ID) })
      this.defaultCheckedRoleApi = ids
      this.$refs.roleApiTree.setCheckedKeys(this.defaultCheckedRoleApi)
    },
    async updateRoleMenusById() {
      this.permissionLoading = true
      let ids = this.$refs.roleMenuTree.getCheckedKeys()
      const idsHalf = this.$refs.roleMenuTree.getHalfCheckedKeys()
      ids = ids.concat(idsHalf)
      ids = [...new Set(ids)]
      try {
        await updateRoleMenusById(this.roleId, { menuIds: ids })
      } finally {
        this.permissionLoading = false
      }

      this.permsDialogVisible = false
      this.$message({
        showClose: true,
        message: 'Updated successfully',
        type: 'success'
      })
    },
    async updateRoleApisById() {
      this.permissionLoading = true
      const ids = this.$refs.roleApiTree.getCheckedKeys(true)
      try {
        await updateRoleApisById(this.roleId, { apiIds: ids })
      } finally {
        this.permissionLoading = false
      }

      this.permsDialogVisible = false
      this.$message({
        showClose: true,
        message: 'Updated successfully',
        type: 'success'
      })
    },
    submitPermissionForm() {
      this.updateRoleMenusById()
      this.updateRoleApisById()
    },
    cancelPermissionForm() {
      this.permsDialogVisible = false
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

<style scoped >
.container-card {
  margin: 10px;
}

.role-menu {
  font-size: 15px;
}
</style>

<style lang="scss">
.perms-dialog>.el-dialog__body {
  padding-top: 0;
  padding-bottom: 15px;
}
</style>

