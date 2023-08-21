<script setup>
import { ref } from "vue"

let cardTableVisible = $ref(false)
let handleLinkCardTable = () => {
    console.log(cardTableVisible)
    cardTableVisible = !cardTableVisible
}
</script>

<template>
    <div class="container" v-if="!cardTableVisible">
        <v-card>
            <h1 class="title">排谷 微店APP搜索苹果喵藏宝阁</h1>
            <div class="card-intro">
                <p>
                    默认囤货，并已阅读定制周边除流团以外不退不换。如果真的退坑需要强行全部掉落，只退80%的款（或者您可以自行转单）！
                </p>
                <p>补邮步骤：</p>

            </div>
            <v-list>
                <v-list-item class="card-ol">
                    1. 首先核对苹果喵藏宝阁通知群的
                    <a href="https://docs.qq.com/sheet/DRnNkVmlHdEtHV3Bt?tab=BB08J2" target="_blank">
                        <v-btn color="blue-lighten-1">到货表</v-btn>
                    </a>
                </v-list-item>
                <v-list-item class="card-ol">
                    2.
                    <li>
                        <a href="https://docs.qq.com/form/page/DRkFnRFlGbmlpaVlE" target="_blank">
                            <v-btn color="green-lighten-1">快速发货通道登记</v-btn>
                        </a>
                        <router-link to="/cardTable" @click="handleLinkCardTable">
                            <v-btn color="amber-lighten-2">快速发货通道</v-btn>
                        </router-link>
                    </li>
                    <li>
                        <a href="https://docs.qq.com/form/page/DRnJKcUdPR09XSXpT" target="_blank">
                            <v-btn color="red-lighten-1">普通发货通道登记</v-btn>
                        </a>
                        查看自己微店的下单记录，罗列自己的list—&gt;到货表的对应序号+小料名字+数量
                    </li>
                    <li>
                        等待审核，查看
                    </li>
                    <li>
                        付款
                        <a href="https://k.youshop10.com/q7zRMcKH" target="_blank">交邮费</a>
                        ，等待发货
                    </li>
                </v-list-item>
                <v-list-item class="card-ol">
                    3.
                    <a href="https://docs.qq.com/form/page/DRmJ0aFNmVVJjcG1J?u=8f6fa7d1d4264bf3a110fcbf5738d417#/fill"
                        target="_blank">
                        <v-btn color="indigo-darken-1">转单登记</v-btn>
                    </a>
                    <a href="https://docs.qq.com/sheet/DRlNCWGN2YkJ2cGZU?tab=azfe1s" target="_blank">
                        <v-btn color="purple-lighten-1">审核结果查询</v-btn>
                    </a>
                </v-list-item>
            </v-list>
        </v-card>
    </div>
    <router-view v-if="cardTableVisible"></router-view>
</template>
<style scoped>
.v-btn {
    margin: 0 10px;

}
</style>

