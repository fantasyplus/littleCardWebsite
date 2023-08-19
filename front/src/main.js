import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'amfe-flexible'
import ElTableInfiniteScroll from "el-table-infinite-scroll";
import router from './utils/router'
import App from './App.vue'

const app = createApp(App)

app.use(ElementPlus)
app.use(ElTableInfiniteScroll);
app.use(router)

app.mount('#app')