<script setup>
import { ref } from "vue"

let cardTableVisible = $ref(false)
let handleLinkCardTable = () => {
    console.log(cardTableVisible)
    cardTableVisible = !cardTableVisible
}
</script>

<template>
    <el-container class="container" v-if="!cardTableVisible">
        <el-card>
            <h1 class="title">排谷 微店APP搜索苹果喵藏宝阁</h1>
            <el-container class="card-intro">
                <p>
                    默认囤货，并已阅读定制周边除流团以外不退不换。如果真的退坑需要强行全部掉落，只退80%的款（或者您可以自行转单）！
                </p>
                <p>补邮步骤：</p>

            </el-container>
            <ol class="card-ol">
                <li>
                    首先核对苹果喵藏宝阁的
                    <a href="https://docs.qq.com/sheet/DRnNkVmlHdEtHV3Bt?tab=BB08J2" target="_blank">到货表</a>
                </li>
                <li>
                    <ul>
                        <li>
                            <router-link to="/cardTable">
                                <el-button type="primary">快速发货通道</el-button>
                            </router-link>
                        </li>
                        <li>
                            查看自己微店的下单记录，罗列自己的list—&gt;到货表的对应序号+小料名字+数量
                        </li>
                        <li>
                            等待审核，查看
                            <a href="https://docs.qq.com/form/page/DRkFnRFlGbmlpaVlE" target="_blank">补邮收集进度</a>
                        </li>
                        <li>
                            付款
                            <a href="https://k.youshop10.com/q7zRMcKH" target="_blank">交邮费</a>
                            ，等待发货
                        </li>
                    </ul>
                </li>
                <li>转单</li>
            </ol>
        </el-card>
    </el-container>
    <router-view v-if="cardTableVisible"></router-view>
</template>
<style scoped>
.container {
    padding: 20px;
}

.title {
    font-size: 28px;
    /* 调整标题的字体大小 */
    margin-bottom: 10px;
}

.router-link {
    color: #409EFF;
}

.card-ol {
    font-size: large;
}

.card-intro {
    font-size: large;
}
</style>

