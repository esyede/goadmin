import request from '@/utils/request'

export function getRoles(params) {
  return request({
    url: '/api/role/list',
    method: 'get',
    params
  })
}

export function createRole(data) {
  return request({
    url: '/api/role/create',
    method: 'post',
    data
  })
}

export function updateRoleById(roleId, data) {
  return request({
    url: '/api/role/update/' + roleId,
    method: 'patch',
    data
  })
}

export function getRoleMenusById(roleId) {
  return request({
    url: '/api/role/menus/get/' + roleId,
    method: 'get'
  })
}

export function updateRoleMenusById(roleId, data) {
  return request({
    url: '/api/role/menus/update/' + roleId,
    method: 'patch',
    data
  })
}

export function getRoleApisById(roleId) {
  return request({
    url: '/api/role/apis/get/' + roleId,
    method: 'get'
  })
}

export function updateRoleApisById(roleId, data) {
  return request({
    url: '/api/role/apis/update/' + roleId,
    method: 'patch',
    data
  })
}

export function batchDeleteRoleByIds(data) {
  return request({
    url: '/api/role/delete/batch',
    method: 'delete',
    data
  })
}
