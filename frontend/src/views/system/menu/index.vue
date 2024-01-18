<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" class="demo-form-inline">
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">Create</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger"
            @click="batchDelete">Batch Delete</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" row-key="ID"
        :data="tableData" border stripe style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip prop="title" label="Title" width="150" />
        <el-table-column show-overflow-tooltip prop="name" label="Name" />
        <el-table-column show-overflow-tooltip prop="icon" label="Icon" />
        <el-table-column show-overflow-tooltip prop="path" label="Path" />
        <el-table-column show-overflow-tooltip prop="component" label="Component Path" />
        <el-table-column show-overflow-tooltip prop="redirect" label="Redirect" />
        <el-table-column show-overflow-tooltip prop="sort" label="Sort" align="center" width="80" />
        <el-table-column show-overflow-tooltip prop="status" label="Status" align="center" width="80">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? 'Off' : 'On' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip prop="hidden" label="Hidden" align="center" width="80">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.hidden === 1 ? 'danger' : 'success'">
              {{ scope.row.hidden === 1 ? 'Yes' : 'No' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip prop="noCache" label="Cache" align="center" width="80">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.noCache === 1 ? 'danger' : 'success'">
              {{ scope.row.noCache === 1 ? 'No' : 'Yes' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip prop="activeMenu" label="Highlight menu" />
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

      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="580px">
        <el-form ref="dialogForm" :inline="true" size="small" :model="dialogFormData" :rules="dialogFormRules"
          label-width="80px">
          <el-form-item label="Menu Title" prop="title">
            <el-input v-model.trim="dialogFormData.title" placeholder="Menu title" style="width: 440px" />
          </el-form-item>
          <el-form-item label="Name" prop="name">
            <el-input v-model.trim="dialogFormData.name" placeholder="Menu name" style="width: 220px" />
          </el-form-item>
          <el-form-item label="Sort" prop="sort">
            <el-input-number v-model.number="dialogFormData.sort" controls-position="right" :min="1" :max="999" />
          </el-form-item>
          <el-form-item label="Icon" prop="icon">
            <el-popover placement="bottom-start" width="450" trigger="click" @show="$refs['iconSelect'].reset()">
              <IconSelect ref="iconSelect" @selected="selected" />
              <el-input slot="reference" v-model="dialogFormData.icon" style="width: 440px;" placeholder="Select icon"
                readonly>
                <svg-icon v-if="dialogFormData.icon" slot="prefix" :icon-class="dialogFormData.icon"
                  class="el-input__icon" style="height: 32px;width: 16px;" />
                <i v-else slot="prefix" class="el-icon-search el-input__icon" />
              </el-input>
            </el-popover>
          </el-form-item>
          <el-form-item label="Path" prop="path">
            <el-input v-model.trim="dialogFormData.path" placeholder="Route path" style="width: 440px" />
          </el-form-item>
          <el-form-item label="Component" prop="component">
            <el-input v-model.trim="dialogFormData.component" placeholder="Component path" style="width: 440px" />
          </el-form-item>
          <el-form-item label="Redirect" prop="redirect">
            <el-input v-model.trim="dialogFormData.redirect" placeholder="Redirect path" style="width: 440px" />
          </el-form-item>
          <el-form-item label="Status" prop="status">
            <el-radio-group v-model="dialogFormData.status">
              <el-radio-button label="On" />
              <el-radio-button label="Off" />
            </el-radio-group>
          </el-form-item>
          <el-form-item label="Hidden" prop="hidden">
            <el-radio-group v-model="dialogFormData.hidden">
              <el-radio-button label="Yes" />
              <el-radio-button label="No" />
            </el-radio-group>
          </el-form-item>
          <el-form-item label="Cache" prop="noCache">
            <el-radio-group v-model="dialogFormData.noCache">
              <el-radio-button label="On" />
              <el-radio-button label="Off" />
            </el-radio-group>
          </el-form-item>
          <el-form-item label="Highlight Menu" prop="activeMenu">
            <el-input v-model.trim="dialogFormData.activeMenu" placeholder="Highlight menu" style="width: 440px" />
          </el-form-item>
          <el-form-item label="Parent ID" prop="parentId">
            <!-- <el-cascader v-model="dialogFormData.parentId" :show-all-levels="false" :options="treeselectData"
              :props="{ checkStrictly: true, label: 'title', value: 'ID', emitPath: false }" clearable filterable /> -->
            <treeselect v-model="dialogFormData.parentId" :options="treeselectData" :normalizer="normalizer"
              style="width:440px" @input="treeselectInput" />
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
import { batchDeleteMenuByIds, createMenu, getMenuTree, updateMenuById } from '@/api/system/menu'
import IconSelect from '@/components/IconSelect'
import Treeselect from '@riophae/vue-treeselect'
import '@riophae/vue-treeselect/dist/vue-treeselect.css'

export default {
  name: 'Menu',
  components: {
    IconSelect,
    Treeselect
  },
  data() {
    return {
      tableData: [],
      loading: false,
      treeselectData: [],
      treeselectValue: 0,
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        title: '',
        name: '',
        icon: '',
        path: '',
        component: 'Layout',
        redirect: '',
        sort: 999,
        status: '否',
        hidden: '否',
        noCache: '是',
        alwaysShow: 2,
        breadcrumb: 1,
        activeMenu: '',
        parentId: 0
      },
      dialogFormRules: {
        title: [
          { required: true, message: 'Please enter title', trigger: 'blur' },
          { min: 1, max: 50, message: 'Must be between 1 - 50 characters', trigger: 'blur' }
        ],
        name: [
          { required: true, message: 'Please enter name', trigger: 'blur' },
          { min: 1, max: 100, message: 'Must be between 1 - 100 characters', trigger: 'blur' }
        ],
        path: [
          { required: true, message: 'Please enter route path', trigger: 'blur' },
          { min: 1, max: 100, message: 'Must be between 1 - 100 characters', trigger: 'blur' }
        ],
        component: [
          { required: false, message: 'Please enter component path', trigger: 'blur' },
          { min: 0, max: 100, message: 'Must be less than 100 characters', trigger: 'blur' }
        ],
        redirect: [
          { required: false, message: 'Please enter redirect path', trigger: 'blur' },
          { min: 0, max: 100, message: 'Must be less than 100 characters', trigger: 'blur' }
        ],
        activeMenu: [
          { required: false, message: 'Please enter highlight menu', trigger: 'blur' },
          { min: 0, max: 100, message: 'Must be less than 100 characters', trigger: 'blur' }
        ],
        parentId: [
          { required: true, message: 'Please enter parent ID', trigger: 'change' }
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
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getMenuTree()
        this.tableData = data.menuTree
        this.treeselectData = [{ ID: 0, title: 'Root', children: data.menuTree }]
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
      this.dialogFormData.title = row.title
      this.dialogFormData.name = row.name
      this.dialogFormData.icon = row.icon
      this.dialogFormData.path = row.path
      this.dialogFormData.component = row.component
      this.dialogFormData.redirect = row.redirect
      this.dialogFormData.sort = row.sort
      this.dialogFormData.status = row.status === 1 ? 'Off' : 'On'
      this.dialogFormData.hidden = row.hidden === 1 ? 'Yes' : 'No'
      this.dialogFormData.noCache = row.noCache === 1 ? 'No' : 'Yes'
      this.dialogFormData.activeMenu = row.activeMenu
      this.dialogFormData.parentId = row.parentId
      this.dialogFormTitle = 'Edit'
      this.dialogType = 'update'
      this.dialogFormVisible = true
    },
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.submitLoading = true

          if (this.dialogFormData.ID === this.dialogFormData.parentId) {
            return this.$message({
              showClose: true,
              message: 'Cannot select self as parent',
              type: 'error'
            })
          }

          if (this.dialogFormData.component === '') {
            this.dialogFormData.component = 'Layout'
          }

          this.dialogFormData.status = this.dialogFormData.status === 'On' ? 2 : 1
          this.dialogFormData.hidden = this.dialogFormData.hidden === 'Yes' ? 1 : 2
          this.dialogFormData.noCache = this.dialogFormData.noCache === 'Yes' ? 2 : 1

          const dialogFormDataCopy = { ...this.dialogFormData, parentId: this.treeselectValue }
          let msg = ''
          try {
            if (this.dialogType === 'create') {
              const { message } = await createMenu(dialogFormDataCopy)
              msg = message
            } else {
              const { message } = await updateMenuById(dialogFormDataCopy.ID, dialogFormDataCopy)
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
        title: '',
        name: '',
        icon: '',
        path: '',
        component: 'Layout',
        redirect: '',
        sort: 999,
        status: 'Off',
        hidden: 'No',
        noCache: 'Yes',
        alwaysShow: 2,
        breadcrumb: 1,
        activeMenu: '',
        parentId: 0
      }
    },
    batchDelete() {
      this.$confirm('This cannot be undone. Do you want to continue?', 'Delete', {
        confirmButtonText: 'Yes',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const menuIds = []
        this.multipleSelection.forEach(x => {
          menuIds.push(x.ID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteMenuByIds({ menuIds: menuIds })
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
    async singleDelete(Id) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await batchDeleteMenuByIds({ menuIds: [Id] })
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
    selected(name) {
      this.dialogFormData.icon = name
    },
    normalizer(node) {
      return {
        id: node.ID,
        label: node.title,
        children: node.children
      }
    },
    treeselectInput(value) {
      this.treeselectValue = value
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
