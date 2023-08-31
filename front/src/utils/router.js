import { createRouter, createWebHistory, createWebHashHistory } from 'vue-router'

const words = () => import('../components/words.vue')
const searchPage = () => import('../components/searchPage.vue')
const manager = () => import('../components/manager.vue')
const test = () => import('../components/test.vue')
const routes = [
    {
        path: "/",
        component: manager,
    },
    {
        name: "searchPage",
        path: "/searchPage",
        component: searchPage,
    },
    {
        name: "words",
        path: "/words",
        component: words,
    },
    {
        name: "manager",
        path: "/manager",
        component: manager
    }
]

const router = createRouter({
    routes,
    history: createWebHashHistory(),
})

export default router