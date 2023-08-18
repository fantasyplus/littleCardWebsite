import { createRouter, createWebHistory, createWebHashHistory } from 'vue-router'

const words =() => import('../components/words.vue')
const cardTable =() => import('../components/cardTable.vue')

const routes = [
    {
        path: "/",
        component: words,
    },
    {
        name: "cardTable",
        path: "/cardTable",
        component: cardTable,
    },
    {
        name: "words",
        path: "/words",
        component: words,
    }
]

const router = createRouter({
    routes,
    history: createWebHashHistory(),
})

export default router