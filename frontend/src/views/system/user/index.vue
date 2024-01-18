<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="Username">
          <el-input v-model.trim="params.username" clearable placeholder="Username" @clear="search" />
        </el-form-item>
        <el-form-item label="Fullname">
          <el-input v-model.trim="params.nickname" clearable placeholder="Fullname" @clear="search" />
        </el-form-item>
        <el-form-item label="Status">
          <el-select v-model.trim="params.status" clearable placeholder="Status" @change="search" @clear="search">
            <el-option label="On" value="1" />
            <el-option label="Off" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="Phone">
          <el-input v-model.trim="params.mobile" clearable placeholder="Phone" @clear="search" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="primary" @click="search">Search</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="warning" @click="create">Create</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :Off="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger"
            @click="batchDelete">Batch Delete</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%"
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="username" label="Username" />
        <el-table-column show-overflow-tooltip sortable prop="nickname" label="Fullname" />
        <el-table-column show-overflow-tooltip sortable prop="status" label="Status" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status === 1 ? 'success' : 'danger'" disable-transitions>
              {{ scope.row.status === 1 ? 'On' : 'Off' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="mobile" label="Phone" />
        <el-table-column show-overflow-tooltip sortable prop="creator" label="Creator" />
        <el-table-column show-overflow-tooltip sortable prop="introduction" label="Description" />
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

      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="30%">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="100px">
          <el-form-item label="Username" prop="username">
            <el-input ref="password" v-model.trim="dialogFormData.username" placeholder="Username" />
          </el-form-item>
          <el-form-item :label="dialogType === 'create' ? 'New Password' : 'Reset Password'" prop="password">
            <el-input v-model.trim="dialogFormData.password" autocomplete="off" :type="passwordType"
              :placeholder="dialogType === 'create' ? 'New password' : 'Reset password'" />
            <span class="show-pwd" @click="showPwd">
              <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
            </span>
          </el-form-item>
          <el-form-item label="Role" prop="roleIds">
            <el-select v-model.trim="dialogFormData.roleIds" multiple placeholder="Select role" style="width:100%">
              <el-option v-for="item in roles" :key="item.ID" :label="item.name" :value="item.ID" />
            </el-select>
          </el-form-item>
          <el-form-item label="Status" prop="status">
            <el-select v-model.trim="dialogFormData.status" placeholder="Select status" style="width:100%">
              <el-option label="On" :value="1" />
              <el-option label="Off" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item label="Fullname" prop="nickname">
            <el-input v-model.trim="dialogFormData.nickname" placeholder="Fullname" />
          </el-form-item>
          <el-form-item label="Phone" prop="mobile">
            <el-input v-model.trim="dialogFormData.mobile" placeholder="Phone" />
          </el-form-item>
          <el-form-item label="Description" prop="introduction">
            <el-input v-model.trim="dialogFormData.introduction" type="textarea" placeholder="Description" show-word-limit
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
import { getRoles } from '@/api/system/role'
import { batchDeleteUserByIds, createUser, getUsers, updateUserById } from '@/api/system/user'
import JSEncrypt from 'jsencrypt'

export default {
  name: 'User',
  data() {
    var checkPhone = (rule, value, callback) => {
      if (!value) {
        return callback(new Error('Phone cannot be empty'))
      } else {
        const rgx = new RegExp('^\\+[1-9]{1}[0-9]{0,2}-[2-9]{1}[0-9]{2}-[2-9]{1}[0-9]{2}-[0-9]{4}$')
        if (rgx.test(value)) {
          callback()
        } else {
          return callback(new Error('Phone number is invalid'))
        }
      }
    }
    return {
      params: {
        username: '',
        nickname: '',
        status: '',
        mobile: '',
        pageNum: 1,
        pageSize: 10
      },
      tableData: [],
      total: 0,
      loading: false,
      roles: [],
      passwordType: 'password',
      publicKey: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDbOYcY8HbDaNM9ooYXoc9s+R5o
R05ZL1BsVKadQBgOVH/kj7PQuD+ABEFVgB6rJNi287fRuZeZR+MCoG72H+AYsAhR
sEaB5SuI7gDEstXuTyjhx5bz0wUujbDK4VMgRfPO6MQo+A0c95OadDEvEQDG3KBQ
wLXapv+ZfsjG7NgdawIDAQAB
-----END PUBLIC KEY-----`,
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        username: '',
        password: '',
        nickname: '',
        status: 1,
        mobile: '',
        avatar: '',
        introduction: '',
        roleIds: ''
      },
      dialogFormRules: {
        username: [
          { required: true, message: 'Please enter username', trigger: 'blur' },
          { min: 2, max: 20, message: 'Must be between 2 - 20 characters', trigger: 'blur' }
        ],
        password: [
          { required: false, message: 'Please enter password', trigger: 'blur' },
          { min: 6, max: 30, message: 'Must be between 6 - 30 characters', trigger: 'blur' }
        ],
        nickname: [
          { required: false, message: 'Please enter fullname', trigger: 'blur' },
          { min: 2, max: 20, message: 'Must be between 2 - 20 characters', trigger: 'blur' }
        ],
        mobile: [
          { required: true, validator: checkPhone, trigger: 'blur' }
        ],
        status: [
          { required: true, message: 'Please choose status', trigger: 'change' }
        ],
        introduction: [
          { required: false, message: 'Please enter description', trigger: 'blur' },
          { min: 0, max: 100, message: 'Must be less than 100 characters', trigger: 'blur' }
        ]
      },
      popoverVisible: false,
      multipleSelection: []
    }
  },
  created() {
    this.getTableData()
    this.getRoles()
  },
  methods: {
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getUsers(this.params)
        this.tableData = data.users
        this.total = data.total
      } finally {
        this.loading = false
      }
    },
    async getRoles() {
      const res = await getRoles(null)

      this.roles = res.data.roles
    },
    create() {
      this.dialogFormTitle = 'Create'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },
    update(row) {
      this.dialogFormData.ID = row.ID
      this.dialogFormData.username = row.username
      this.dialogFormData.password = ''
      this.dialogFormData.nickname = row.nickname
      this.dialogFormData.status = row.status
      this.dialogFormData.mobile = row.mobile
      this.dialogFormData.introduction = row.introduction
      this.dialogFormData.roleIds = row.roleIds
      this.dialogFormTitle = 'Edit'
      this.dialogType = 'update'
      this.passwordType = 'password'
      this.dialogFormVisible = true
    },
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.submitLoading = true

          this.dialogFormDataCopy = { ...this.dialogFormData }
          if (this.dialogFormData.password !== '') {
            const encryptor = new JSEncrypt()
            encryptor.setPublicKey(this.publicKey)
            const encPassword = encryptor.encrypt(this.dialogFormData.password)
            this.dialogFormDataCopy.password = encPassword
          }
          let msg = ''
          try {
            if (this.dialogType === 'create') {
              const { message } = await createUser(this.dialogFormDataCopy)
              msg = message
            } else {
              const { message } = await updateUserById(this.dialogFormDataCopy.ID, this.dialogFormDataCopy)
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
        username: '',
        password: '',
        nickname: '',
        status: 1,
        mobile: '',
        avatar: '',
        introduction: '',
        roleIds: ''
      }
    },
    batchDelete() {
      this.$confirm('This cannot be undone. Do you want to continue?', 'Delete', {
        confirmButtonText: 'Yes',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const userIds = []
        this.multipleSelection.forEach(x => {
          userIds.push(x.ID)
        })
        let msg = ''
        try {
          const { message } = await batchDeleteUserByIds({ userIds: userIds })
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
        const { message } = await batchDeleteUserByIds({ userIds: [Id] })
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
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
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

.show-pwd {
  position: absolute;
  right: 10px;
  top: 3px;
  font-size: 16px;
  color: #889aa4;
  cursor: pointer;
  user-select: none;
}
</style>
