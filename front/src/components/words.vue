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
        <h1>排谷 微店APP搜索苹果喵藏宝阁</h1>
        <p>默认囤货，并已阅读定制周边除流团以外不退不换。如果真的退坑需要强行全部掉落，只退80%的款（或者您可以自行转单）！</p>
        <p>补邮步骤：</p>
        <ol>
            <li>首先核对苹果喵藏宝阁的<router-link to="/cardTable" @click="handleLinkCardTable">到货表</router-link></li>
            <li>
                <ul>
                    <li>快速发货通道（Beta）</li>
                    <li>查看自己微店的下单记录，罗列自己的list—&gt;到货表的对应序号+小料名字+数量</li>
                    <li>等待审核，查看补邮收集进度 <a class="link" href="链接补邮收集进度">链接</a></li>
                    <li>付款交邮费 <a class="link" href="链接付款交邮费">链接</a>，等待发货</li>
                </ul>
            </li>
            <li>转单</li>
        </ol>
    </div>
    <router-view v-if="cardTableVisible"></router-view>
</template>

<style scoped>
.container {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
}

.router-link {
    color: #409EFF;
}
</style>

