<script setup>
import { ref, toRaw } from "vue"
import request from "../utils/request.js"
import axios from "axios"

//vuetify labs
import { VDataTableVirtual } from 'vuetify/labs/VDataTable'

// dataww 
let currentView = $ref("list")
let CNInput = $ref("")
let QQInput = $ref("")
let showSearchData = $ref(false)
let selectedCardsInfo = $ref([])
let cardExpandList = $ref([])
let cardSelectedInfo = $ref("全部发货")
let cardSelectedList = $ref([])
//卡片的基础信息
let CardInfoList = $ref([{
    CardID: "",
    card_name: "",
    card_character: "",
    card_type: "",
    card_condition: "",
    card_num: "",
    card_deliver: "",
    other_info: "",
}])

const tableHeaders = $ref([
    {
        title: 'CardID',
        align: 'start',
        sortable: true,
        key: 'CardID',
    },
    { title: 'CardName', align: 'end', key: 'card_name' },
    { title: 'CardCharacter', align: 'end', key: 'card_character' },
    { title: 'CardType', align: 'end', key: 'card_type' },
    { title: 'CardCondition', align: 'end', key: 'card_condition' },
    { title: 'CardNum', align: 'end', key: 'card_num' },
    { title: 'CardDeliver', align: 'end', key: 'card_deliver' },
    { title: 'OtherInfo', align: 'end', key: 'other_info' },
])

// methods

const requestWithParams = async (path, params) => {
    let encodedParams = {};
    for (let key in params) {
        encodedParams[key] = encodeURIComponent(params[key]);
    }

    let res = await request.get(path, encodedParams);
    return res;
}

//获取所有卡片信息
const getCardInfoList = async (page_num = 1) => {
    let params = {
        pageSize: 100,
        pageNum: page_num
    };
    let res = await requestWithParams(`/data/list`, params);

    // Filter out items with card_name as "none"
    const filteredData = res.data.data.filter(item => item.card_name !== "None");

    // Update CardInfoList with filteredData
    CardInfoList = filteredData;
    console.log("get CardInfoList:", res)

    for (let i = 0; i < CardInfoList.length; i++) {
        cardExpandList.push(false)
        cardSelectedList.push(false)
    }
    // console.log("init cardExpandList:", cardExpandList)
}


//download
const hanleDownload = async () => {
    await axios.get("/test/download", {
        responseType: "blob"
    }).then(res => {
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'filename.xlsx');
        document.body.appendChild(link);
        link.click();
    })

}

//搜索CN和QQ对应的谷子
const handleSearchCNorQQ = async () => {
    // search from server
    if (CNInput.length > 0 || QQInput.length > 0) {
        let params = {
            cn: CNInput,
            qq: QQInput
        };
        let res = await requestWithParams(`/data/search`, params);
        // console.log("search res:", res.data)
        CardInfoList = res.data

        CardInfoList.sort((a, b) => a.CardID - b.CardID);

        // console.log("search CardInfoList:", CardInfoList)
        if (res.code === 400) {
            showSearchData = false

        }
        else if (res.code === 200) {
            showSearchData = true
            //清空选择结果
            cardSelectedList.splice(false, cardSelectedList.length);
            console.log(cardSelectedList)
            // alert("找到对应的谷子")
        }
    }
    else {
        await getCardInfoList()
    }
}


const handleCardListClick = (index) => {
    cardExpandList[index] = !cardExpandList[index]
}


const handleCardCheckBoxClick = (item, index) => {
    cardSelectedList[index] = !cardSelectedList[index]

    if (cardSelectedList[index]) {
        selectedCardsInfo.push(item)
        // selectedCardsInfo.push()
    }
    else {
        const filteredData = selectedCardsInfo.filter((value) => {
            return value.CardID !== item.CardID
        })
        selectedCardsInfo = filteredData
    }
    // console.log(selectedCardsInfo)
}

//把用户选择的谷子信息发送到后端
//后端根据这些谷子修改状态（发货或其他）
const handleSaveSelectedCards = () => {
    console.log(selectedCardsInfo)
}

const handleGenerateDeliverList = () => {
    console.log("handleGenerateDeliverList")
}

//初始化所有谷子信息
getCardInfoList()
</script>

<template>
    <div class="search-box">

        <div class="query-box">
            <div class="title">
                <h2>到货List</h2>
            </div>
            <h3>请再次输入CN或QQ（注：只显示已到货且能发货的！！）</h3>
            <div class="input-box">
                <v-text-field label="CN" class="search-input" v-model="CNInput" color="light-blue-lighten-1"></v-text-field>
                <v-text-field label="QQ" class="search-input" v-model="QQInput" color="light-blue-lighten-1"></v-text-field>
            </div>

            <div class="btn-list">
                <v-btn @click="handleSearchCNorQQ" color="blue-lighten-2" style="width: 100px">搜索</v-btn>
                <!-- <v-btn @click="hanleDownload" color="green-darken-2">下载</v-btn> -->
            </div>
        </div>


        <div v-if="showSearchData">
            <div class="toggle-icons" style="display: flex; justify-content: flex-end; align-items: flex-start;">
                <v-icon @click="currentView = 'list'">mdi-menu</v-icon>
                <v-icon @click="currentView = 'table'">mdi-apps
                    <v-tooltip activator="parent" location="top">Tooltip</v-tooltip>
                </v-icon>
            </div>

            <div class="table-container" v-if="currentView === 'table'">
            <v-data-table-virtual class="card-table" :headers="tableHeaders"
                    :items="CardInfoList"></v-data-table-virtual>
            </div>

            <div class="card-list-container" v-if="currentView === 'list'">
                <v-list-subheader>谷子清单</v-list-subheader>
                <v-card v-for=" (item, i) in CardInfoList" :key="i" :value="item" variant="flat">
                    <v-row>
                        <v-col :cols="cardSelectedList[i] ? 6 : 9">
                            <v-card @click="handleCardListClick(i)" @click.stop>
                                <v-row>
                                    <v-col cols="6" class="d-flex align-center">
                                        <v-card-title class="card-fancy-title">
                                            <v-list-item-title>{{ item.card_name }}</v-list-item-title>
                                            <v-list-item-subtitle>{{ item.card_character }}</v-list-item-subtitle>
                                            <v-list-item-subtitle v-if="cardSelectedList[i]" class="d-flex align-center">
                                                数量:{{ item.card_num }}</v-list-item-subtitle>
                                        </v-card-title>
                                    </v-col>
                                    <v-col cols="3" v-if="!cardSelectedList[i]" class="d-flex align-center">
                                        <v-card-title class="card-fancy-title">
                                            <v-list-item-title>数量:{{ item.card_num }}</v-list-item-title>
                                        </v-card-title>
                                    </v-col>
                                    <v-col v-if="!cardSelectedList[i]" cols="2" class="d-flex align-center">
                                        <v-img :src="`/src/assets/${item.CardID}`" width="auto" height="auto"></v-img>
                                    </v-col>
                                </v-row>
                            </v-card>
                        </v-col>

                        <v-col :cols="cardSelectedList[i] ? 6 : 3" class="d-flex align-center">
                            <v-row>
                                <v-col cols="3">
                                    <v-checkbox @click.stop v-model="cardSelectedList[i]"
                                        @click="handleCardCheckBoxClick(item, i)"></v-checkbox>
                                </v-col>
                                <v-col cols="9">
                                    <v-text-field @click.stop v-if="cardSelectedList[i]" v-model="cardSelectedInfo"
                                        label="默认为全发" color="light-blue-lighten-1"></v-text-field>
                                </v-col>
                            </v-row>
                        </v-col>
                    </v-row>

                    <div class="card-overlay">
                        <v-overlay v-model="cardExpandList[i]" activator="parent" location-strategy="connected"
                            location="top center" origin="auto" scroll-strategy="none">
                            <v-card>
                                <!-- <v-col no-gutters> -->
                                <v-col>
                                    <v-sheet class="pa-2 ma-2">
                                        <v-card-title>
                                            <v-list-item-title>{{ item.card_name }}</v-list-item-title>
                                            <v-list-item-subtitle>{{ item.card_character
                                            }}</v-list-item-subtitle>
                                        </v-card-title>
                                        <v-card-text>
                                            <v-list-item-title>CardID: {{ item.CardID }}</v-list-item-title>
                                            <v-list-item-subtitle>CardType: {{ item.card_type
                                            }}</v-list-item-subtitle>
                                            <v-list-item-subtitle>CardCondition: {{ item.card_condition
                                            }}</v-list-item-subtitle>
                                            <v-list-item-subtitle v-show="item.card_num">CardNum:{{
                                                item.card_num
                                            }}</v-list-item-subtitle>
                                            <v-list-item-subtitle>CardDeliver: {{ item.card_deliver
                                            }}</v-list-item-subtitle>
                                            <v-list-item-subtitle>OtherInfo: {{ item.other_info
                                            }}</v-list-item-subtitle>
                                        </v-card-text>
                                    </v-sheet>
                                </v-col>
                                <v-col>
                                    <v-sheet class="pa-2 ma-2">
                                        <v-img src="https://cdn.vuetifyjs.com/images/cards/foster.jpg" width="100%"
                                            height="auto"></v-img>
                                    </v-sheet>
                                </v-col>
                            </v-card>
                        </v-overlay>
                    </div>
                </v-card>
            </div>
            <div style="display: flex; align-items: center;gap: 10px;width: 100%; justify-content: center;">
                <v-btn @click="handleSaveSelectedCards" color="blue-lighten-1">确认</v-btn>
                <v-btn @click="handleGenerateDeliverList" color="blue-lighten-1">生成发货list</v-btn>
            </div>


        </div>
    </div>
</template>

<style scoped>
.search-box {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    /* background-color: #f7f7f7; */
}

.title {
    text-align: center;
}

.query-box {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    padding: 20px;
}

.input-box {
    display: flex;
    gap: 10px;
    width: 100%;
    justify-content: center;
}

.search-input {
    /* width: 200px; */
    height: 100%;
    box-sizing: border-box;
}

.btn-list {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    justify-content: center;
}

.card-list-container {
    padding: 5px;
}

.card-fancy-title {
    font-size: calc(1em + 1vw);
    /* 使用基本大小和基于视口宽度的大小组合 */
    color: #506eae;
    /* 设置字体颜色 */
}
</style>
