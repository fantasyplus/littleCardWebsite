<script setup>
import { ref, toRaw } from "vue"
import request from "../utils/request.js"
import axios from "axios"

//vuetify labs
import { VDataTableVirtual } from 'vuetify/labs/VDataTable'

// data
let searchInput = $ref("")
//卡片的基础信息
let CardInfoList = $ref([{
    CardID: "",
    card_name: "",
    card_character: "",
    card_type: "",
    card_condition: "",
    card_num: "",
    other_info: "",
}])
let CardInfo = $ref({
    CardID: "",
    card_name: "",
    card_character: "",
    card_type: "",
    card_condition: "",
    card_num: "",
    other_info: "",
})

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
    { title: 'OtherInfo', align: 'end', key: 'other_info' },
])


let cardExpandList = $ref([])
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
    if (searchInput.length > 0) {
        let params = {
            cn: searchInput,
            qq: searchInput
        };
        let res = await requestWithParams(`/data/search`, params);
        console.log("search res:", res.data)
        CardInfoList = res.data
        console.log("search CardInfoList:", CardInfoList)
    }
    else {
        await getCardInfoList()
    }
}


const handleCardListClick = (index) => {
    cardExpandList[index] = !cardExpandList[index]
}


//初始化所有谷子信息
getCardInfoList()
</script>

<template>
    <!-- test -->
    <div class="table-box">
        <!-- Title -->
        <div class="title">
            <h2>到货List</h2>
        </div>
        <!-- query -->
        <h3>请再次输入CN或QQ（注：只显示已到货且能发货的！！）</h3>
        <div class="query-box">
            <v-text-field label="CN或QQ" class="search-input" @input="handleSearchCNorQQ"
            v-model="searchInput"></v-text-field>
            <div class="btn-list">
                <v-btn @click="hanleDownload" color="green-darken-2">Download</v-btn>
            </div>
        </div>

        <!-- table -->
        <div class="table-container">
            <v-data-table-virtual class="card-table" :headers="tableHeaders" :items="CardInfoList"></v-data-table-virtual>
        </div>

        <v-container class="card-list-container">
            <v-card>
                <v-list>
                    <v-list-subheader>Card List</v-list-subheader>

                    <v-list-item v-for="(item, i) in CardInfoList" :key="i" :value="item" color="blue-lighten-3"
                        @click="handleCardListClick(i)">
                        <v-card>
                            <v-card-title>
                                <v-list-item-title>{{ item.card_name }}</v-list-item-title>
                                <v-list-item-subtitle>{{ item.card_character }}</v-list-item-subtitle>
                            </v-card-title>
                            <v-card-text>
                                <v-list-item-title>CardID: {{ item.CardID }}</v-list-item-title>
                                <v-list-item-subtitle>CardType: {{ item.card_type }}</v-list-item-subtitle>
                                <v-list-item-subtitle>CardCondition: {{ item.card_condition }}</v-list-item-subtitle>
                                <v-list-item-subtitle v-show="item.card_num">CardNum:{{ item.card_num
                                }}</v-list-item-subtitle>
                                <v-list-item-subtitle>OtherInfo: {{ item.other_info }}</v-list-item-subtitle>
                            </v-card-text>
                            <v-card-actions>
                            </v-card-actions>

                        </v-card>

                        <v-expand-transition>
                            <v-img v-show="cardExpandList[i]" src="https://cdn.vuetifyjs.com/images/cards/foster.jpg"
                                width="100%" height="auto"></v-img>
                        </v-expand-transition>
                    </v-list-item>

                </v-list>
            </v-card>
        </v-container>
    </div>

    <v-container class="bg-surface-variant">
        <v-row no-gutters>
            <v-col>
                <v-sheet class="pa-2 ma-2">
                    .v-col-auto
                </v-sheet>
            </v-col>
            <v-col>
                <v-sheet class="pa-2 ma-2">
                    .v-col-auto
                </v-sheet>
            </v-col>
        </v-row>

        <v-row no-gutters>
            <v-col>
                <v-sheet class="pa-2 ma-2">
                    .v-col-auto
                </v-sheet>
            </v-col>
            <v-col>
                <v-sheet class="pa-2 ma-2">
                    .v-col-auto
                </v-sheet>
            </v-col>
            <v-col>
                <v-sheet class="pa-2 ma-2">
                    .v-col-auto
                </v-sheet>
            </v-col>
        </v-row>

        <v-row no-gutters>
            <v-col cols="2">
                <v-sheet class="pa-2 ma-2">
                    .v-col-2
                </v-sheet>
            </v-col>
            <v-col>
                <v-sheet class="pa-2 ma-2">
                    .v-col-auto
                </v-sheet>
            </v-col>
        </v-row>
    </v-container>
    <router-link to="/">
        <v-btn>
            返回
        </v-btn>
    </router-link>
</template>

<style scoped>
.title {
    text-align: center;
}

.query-box {
    display: flex;
    justify-content: space-between;
    margin-bottom: 20px;
}
</style>
