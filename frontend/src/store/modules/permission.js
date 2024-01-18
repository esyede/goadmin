import { getUserMenuTreeByUserId } from '@/api/system/menu'
import Layout from '@/layout'
import { constantRoutes } from '@/router'

export const getRoutesFromMenuTree = (menuTree) => {
  const routes = []
  menuTree.forEach(menu => {
    if (menu.children && menu.children.length > 0) {
      menu.children = getRoutesFromMenuTree(menu.children)
    }
    // else {
    //   // Children need to be cleared here, otherwise the drop-down icon will be displayed on the right
    //   delete menu.children
    // }
    routes.push({
      path: menu.path,
      name: menu.name,
      component: loadComponent(menu.component),
      hidden: menu.hidden === 1,
      redirect: menu.redirect,
      alwaysShow: menu.alwaysShow === 1,
      children: menu.children,
      meta: {
        name: menu.name,
        title: menu.title,
        icon: menu.icon,
        noCache: menu.noCache === 1,
        breadcrumb: menu.breadcrumb === 1,
        activeMenu: menu.activeMenu
      }
    })
  })

  return routes
}

export const loadComponent = (component) => {
  if (component === '' || component === 'Layout') {
    return Layout
  }

  return (resolve) => require([`@/views${component}`], resolve)
}

const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  generateRoutes({ commit }, userinfo) {
    return new Promise((resolve, reject) => {
      let accessedRoutes = []
      getUserMenuTreeByUserId(userinfo.id).then(res => {
        const { data } = res
        const menuTree = data.menuTree
        accessedRoutes = getRoutesFromMenuTree(menuTree)
        commit('SET_ROUTES', accessedRoutes)
        resolve(accessedRoutes)
      }).catch(err => {
        reject(err)
      })
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
