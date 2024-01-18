import request from '@/utils/request'

export function getMenuTree() {
  return request({
    url: '/api/menu/tree',
    method: 'get'
  })
}

export function getMenus() {
  return request({
    url: '/api/menu/list',
    method: 'get'
  })
}

export function createMenu(data) {
  return request({
    url: '/api/menu/create',
    method: 'post',
    data
  })
}

export function updateMenuById(Id, data) {
  return request({
    url: '/api/menu/update/' + Id,
    method: 'patch',
    data
  })
}

export function batchDeleteMenuByIds(data) {
  return request({
    url: '/api/menu/delete/batch',
    method: 'delete',
    data
  })
}

export function getUserMenusByUserId(Id) {
  return request({
    url: '/api/menu/access/list/' + Id,
    method: 'get'
  })
}

export function getUserMenuTreeByUserId(Id) {
  return request({
    url: '/api/menu/access/tree/' + Id,
    method: 'get'
  })
}
