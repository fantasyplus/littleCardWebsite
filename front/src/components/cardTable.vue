<script setup>
import { ref, toRaw } from "vue"
import request from "../utils/request.js"
import axios from "axios"

//vuetify labs
import { VDataTableVirtual } from 'vuetify/labs/VDataTable'

// data
let searchInput = $ref("")
let tableData = $ref([{
    CardID: "",
    card_name: "",
    card_character: "",
    card_type: "",
    card_condition: "",
    other_info: "",
}])

let multipleTableSelection = $ref([])
let addDialogVisible = $ref(false)
let dialogType = $ref("add")
let tableRowData = $ref({
    CardID: "",
    card_name: "",
    card_character: "",
    card_type: "",
    card_condition: "",
    other_info: "",
})

let cardExpandList = $ref([])
// methods


//request page data from server
const getTableData = async (page_num = 1) => {
    let res = await request.get("/list", {
        pageSize: 1000,
        pageNum: page_num
    })
    // Filter out items with card_name as "none"
    const filteredData = res.data.data.filter(item => item.card_name !== "None");

    // Update tableData with filteredData
    tableData = filteredData;
    console.log("init tableData:", res)

    for (let i = 0; i < tableData.length; i++) {
        cardExpandList.push(false)
    }
    console.log("init cardExpandList:", cardExpandList)
}
//init table data
getTableData()


// delete row from server
const handleDelRow = async (row) => {
    await request.delete(`/delete/${row.ID}`)
    await getTableData(curPage)
}

// delete mutiple rows from server
const handleDelMultiRows = () => {
    multipleTableSelection.forEach(row => {
        handleDelRow(row)
    })
    multipleTableSelection = []
}


// add
const handleAdd = () => {
    addDialogVisible = true
    dialogType = "add"
    tableRowData = {}
}

//download
const hanleDownload = async () => {
    await axios.get("/user/download", {
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

//edit

// confirm add or edit data

//search
const handleSearch = async () => {
    // search from server
    if (searchInput.length > 0) {
        let res = await request.get(`/list/${searchInput}`)
        tableData = res.data
    }
    else {
        await getTableData(curPage)
    }
}


const handleCardListClick = (index) => {
    cardExpandList[index] = !cardExpandList[index]
}

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
            <v-text-field label="CN或QQ" class="search-input" @input="handleSearch" v-model="searchInput"></v-text-field>
            <div class="btn-list">
                <v-btn @click="handleAdd" color="red-lighten-1">Add</v-btn>
                <v-btn @click="hanleDownload" color="green-darken-2">Download</v-btn>
                <v-btn v-if="multipleTableSelection.length > 0" @click="handleDelMultiRows">Delete</v-btn>
            </div>
        </div>

        <!-- table -->
        <!-- <div class="table-container">
            <v-data-table-virtual class="card-table" :headers="tableHeaders" :items="tableData"></v-data-table-virtual>
        </div> -->

        <v-container class="card-list-container">
            <v-card>
                <v-list>
                    <v-list-subheader>Card List</v-list-subheader>

                    <v-list-item v-for="(item, i) in tableData" :key="i" :value="item" color="grey lighten-2"
                        @click="handleCardListClick(i)">
                        <v-card>
                            <v-card-title>
                                <v-list-item-title>{{ item.card_name }}</v-list-item-title>
                                <v-list-item-subtitle>{{ item.card_character }}</v-list-item-subtitle>
                            </v-card-title>
                            <v-card-text>
                                <v-list-item-title>CardID: {{ item.CardID }}</v-list-item-title>
                                <v-list-item-subtitle>CardType: {{ item.card_type }}</v-list-item-subtitle>
                                <v-list-item-subtitle>CardCondition: {{ item.card_condition
                                }}</v-list-item-subtitle>
                                <v-list-item-subtitle>OtherInfo: {{ item.other_info }}</v-list-item-subtitle>
                            </v-card-text>
                            <v-card-actions>
                                <!-- <v-btn color="red" text @click="handleDelRow(item)">Delete</v-btn> -->
                                <!-- <v-btn color="green" text @click="handleEdit(item)">Edit</v-btn> -->
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
    <!-- Form -->
    <!-- <el-dialog @keyup.enter="handleDialogConfirm" v-model="addDialogVisible" :title="dialogType === 'add' ? 'add' : 'edit'"
        draggable>
        <el-form v-model="tableRowData" :label-position="addDialogLabelPosition" :label-width="120">
            <el-form-item label="CardID">
                <el-input v-model="tableRowData.card_id" autocomplete="off" />
            </el-form-item>
            <el-form-item label="card_name">
                <el-input v-model="tableRowData.card_name" autocomplete="off" />
            </el-form-item>2
            <el-form-item label="card_character">
                <el-input v-model="tableRowData.card_character" autocomplete="off" />
            </el-form-item>
            <el-form-item label="card_type">
                <el-input v-model="tableRowData.card_type" autocomplete="off" />
            </el-form-item>
            <el-form-item label="card_condition">
                <el-input v-model="tableRowData.card_condition" autocomplete="off" />
            </el-form-item>
            <el-form-item label="other_info">
                <el-input v-model="tableRowData.other_info" autocomplete="off" />
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <v-btn @click="addDialogVisible = false">Cancel</v-btn>
                <v-btn  @click="handleDialogConfirm">
                    Confirm
                </v-btn>
            </span>
        </template>
    </el-dialog> -->
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
