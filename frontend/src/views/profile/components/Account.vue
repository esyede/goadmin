<template>
  <div>
    <el-card style="margin-bottom:20px;max-width: 580px;">
      <div slot="header" class="clearfix">
        <span>Change Password</span>
      </div>

      <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="100px">

        <el-form-item label="Current Password" prop="oldPassword">
          <el-input v-model.trim="dialogFormData.oldPassword" autocomplete="on" :type="passwordTypeOld"
            placeholder="Enter current password" />
          <span class="show-pwd" @click="showPwdOld">
            <svg-icon :icon-class="passwordTypeOld === 'password' ? 'eye' : 'eye-open'" />
          </span>
        </el-form-item>

        <el-form-item label="New Password" prop="newPassword">
          <el-input v-model.trim="dialogFormData.newPassword" autocomplete="on" :type="passwordTypeNew"
            placeholder="Enter new password" />
          <span class="show-pwd" @click="showPwdNew">
            <svg-icon :icon-class="passwordTypeNew === 'password' ? 'eye' : 'eye-open'" />
          </span>
        </el-form-item>

        <el-form-item label="Confirm Password" prop="confirmPassword">
          <el-input v-model.trim="dialogFormData.confirmPassword" autocomplete="on" :type="passwordTypeConfirm"
            placeholder="Confirm new password" />
          <span class="show-pwd" @click="showPwdConfirm">
            <svg-icon :icon-class="passwordTypeConfirm === 'password' ? 'eye' : 'eye-open'" />
          </span>
        </el-form-item>

        <el-form-item>
          <el-button :loading="submitLoading" type="primary" @click="submitForm">Save</el-button>
          <el-button @click="cancelForm">Cancel</el-button>
        </el-form-item>

      </el-form>
    </el-card>
  </div>
</template>

<script>
import { changePwd } from '@/api/system/user'
import store from '@/store'
import JSEncrypt from 'jsencrypt'

export default {
  data() {
    const confirmPass = (rule, value, callback) => {
      if (value) {
        if (this.dialogFormData.newPassword !== value) {
          callback(new Error('Password confirmation does not match'))
        } else {
          callback()
        }
      } else {
        callback(new Error('Please enter new password again'))
      }
    }
    return {
      submitLoading: false,
      dialogFormData: {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      dialogFormRules: {
        oldPassword: [
          { required: true, message: 'Please enter current password', trigger: 'blur' },
          { min: 6, max: 30, message: 'Password must be 6 to 30 characters', trigger: 'blur' }
        ],
        newPassword: [
          { required: true, message: '请输入新密码', trigger: 'blur' },
          { min: 6, max: 30, message: 'Password must be 6 to 30 characters', trigger: 'blur' }
        ],
        confirmPassword: [
          { required: true, validator: confirmPass, trigger: 'blur' }
        ]
      },
      publicKey: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDbOYcY8HbDaNM9ooYXoc9s+R5o
R05ZL1BsVKadQBgOVH/kj7PQuD+ABEFVgB6rJNi287fRuZeZR+MCoG72H+AYsAhR
sEaB5SuI7gDEstXuTyjhx5bz0wUujbDK4VMgRfPO6MQo+A0c95OadDEvEQDG3KBQ
wLXapv+ZfsjG7NgdawIDAQAB
-----END PUBLIC KEY-----`,
      passwordTypeOld: 'password',
      passwordTypeNew: 'password',
      passwordTypeConfirm: 'password'
    }
  },
  methods: {
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.dialogFormDataCopy = { ...this.dialogFormData }
          const encryptor = new JSEncrypt()
          encryptor.setPublicKey(this.publicKey)
          const oldPasswd = encryptor.encrypt(this.dialogFormData.oldPassword)
          const newPasswd = encryptor.encrypt(this.dialogFormData.newPassword)
          const confirmPasswd = encryptor.encrypt(this.dialogFormData.confirmPassword)
          this.dialogFormDataCopy.oldPassword = oldPasswd
          this.dialogFormDataCopy.newPassword = newPasswd
          this.dialogFormDataCopy.confirmPassword = confirmPasswd

          this.submitLoading = true
          const { code, message } = await changePwd(this.dialogFormDataCopy)
          this.submitLoading = false

          if (code !== 200) {
            return this.$message({
              showClose: true,
              message: message,
              type: 'error'
            })
          }
          this.resetForm()
          this.$message({
            showClose: true,
            message: 'Password changed successfully, please log in again',
            type: 'success'
          })
          setTimeout(() => {
            store.dispatch('user/logout').then(() => {
              location.reload()
            })
          }, 1500)
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
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      }
    },
    showPwdOld() {
      if (this.passwordTypeOld === 'password') {
        this.passwordTypeOld = ''
      } else {
        this.passwordTypeOld = 'password'
      }
    },
    showPwdNew() {
      if (this.passwordTypeNew === 'password') {
        this.passwordTypeNew = ''
      } else {
        this.passwordTypeNew = 'password'
      }
    },
    showPwdConfirm() {
      if (this.passwordTypeConfirm === 'password') {
        this.passwordTypeConfirm = ''
      } else {
        this.passwordTypeConfirm = 'password'
      }
    }
  }
}
</script>

<style scoped>
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
