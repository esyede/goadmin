import request from '@/utils/request'

export function getInfo() {
  return request({
    url: '/api/user/info',
    method: 'post'
  })
}

export function getUsers(params) {
  return request({
    url: '/api/user/list',
    method: 'get',
    params
  })
}

export function changePwd(data) {
  return request({
    url: '/api/user/changePwd',
    method: 'put',
    data
  })
}

export function createUser(data) {
  return request({
    url: '/api/user/create',
    method: 'post',
    data
  })
}

export function updateUserById(id, data) {
  return request({
    url: '/api/user/update/' + id,
    method: 'patch',
    data
  })
}

export function batchDeleteUserByIds(data) {
  return request({
    url: '/api/user/delete/batch',
    method: 'delete',
    data
  })
}

