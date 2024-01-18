import request from '@/utils/request'

export function getOperationLogs(params) {
  return request({
    url: '/api/log/operation/list',
    method: 'get',
    params
  })
}

export function batchDeleteOperationLogByIds(data) {
  return request({
    url: '/api/log/operation/delete/batch',
    method: 'delete',
    data
  })
}

