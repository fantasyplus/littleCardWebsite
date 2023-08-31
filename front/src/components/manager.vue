<script setup>
import { computed, nextTick, onMounted, ref, toRaw } from 'vue'
import request from '../utils/request';

//vuetify labs
import { VDataTable } from 'vuetify/labs/VDataTable'



const requestWithParams = async (path, params) => {
    let encodedParams = {};
    for (let key in params) {
        encodedParams[key] = encodeURIComponent(params[key]);
    }

    let res = await request.get(path, encodedParams);
    return res;
}

let cardsPerPage = $ref(20)
let totalCards = $ref(0)

let CardTableData = $ref([])
let CardRowData = $ref({
    CardID: "123123",
    card_name: "",
    card_type: "",
    card_character: "",
    card_condition: "",
    other: "",
})
//这里的key是后端返回的字段名
const CardTableHeaders = $ref([
    {
        title: 'card_id',
        align: 'start',
        sortable: true,
        key: 'CardID',
    },
    { title: 'card_name', key: 'card_name', sortable: false },
    { title: 'card_type', key: 'card_type', sortable: false },
    { title: 'card_character', key: 'card_character', sortable: false },
    { title: 'card_condition', key: 'card_condition', sortable: true },
    { title: 'other', key: 'other', sortable: false },
    { title: 'Actions', key: 'actions', sortable: false },
])


//获取所有卡片信息
const getCardTableData = async (page_num = 1) => {
    let params = {
        pageSize: 1000,
        pageNum: page_num
    };
    let res = await requestWithParams(`/data/list`, params);

    const filteredData = res.data.data.filter(item => item.card_name !== "None")
    CardTableData = filteredData
    totalCards = CardTableData.length
    console.log("totalCards", totalCards)
    console.log("get CardTableData:", CardTableData)
}


let DialogVisible = $ref(false)
let deleteAlertDialog = $ref(false)

let editedIndex = $ref(-1)
const dialogTitle = computed(() => {
    return editedIndex === -1 ? '新增谷子' : '编辑谷子'
})

// 关闭添加对话框
const closeAddDialog = () => {
    //关闭对话框
    DialogVisible = false
    nextTick(() => {
        editedIndex = -1
    })
}

// 保存添加的数据
const saveAddDialog = () => {
    console.log("saveAddDialog")
    closeAddDialog()
}

const editCard = (item) => {
    console.log("editCard", item)
    //更改对话框标题
    editedIndex = item.index
    //显示对话框
    DialogVisible = true
    //设置对话框数据
    CardRowData = Object.assign({}, item)
}

const deleteCard = (item) => {
    console.log("deleteCard", item)

    //显示对话框
    deleteAlertDialog = true
}

const closeDeleteAlertDialog = () => {
    //关闭对话框
    deleteAlertDialog = false
}

const deleteCardConfirm = () => {
    console.log("deleteCardConfirm")
    closeDeleteAlertDialog()
}

onMounted(() => {
    getCardTableData();
})

</script>

<template>
    <div class="card-table-container">
        <v-data-table class="elevation-1" v-model:items-per-page="cardsPerPage" :headers="CardTableHeaders"
            :items="CardTableData">

            <template v-slot:top>
                <v-toolbar color="light-blue-lighten-3" dark="true">
                    <v-toolbar-title>时尚小垃圾管理器</v-toolbar-title>
                    <v-divider class="mx-4" inset vertical></v-divider>
                    <v-spacer></v-spacer>

                    <!-- 添加或编辑数据对话框 -->
                    <v-dialog v-model="DialogVisible" max-width="500px">
                        <template v-slot:activator="{ props }">
                            <v-btn dark class="mb-2" v-bind="props">
                                添加谷子
                            </v-btn>
                        </template>
                        <v-card>
                            <v-card-title>
                                <span class="text-h5">{{ dialogTitle }}</span>
                            </v-card-title>

                            <v-card-text>
                                <v-container>
                                    <v-row>
                                        <v-col cols="12" sm="6" md="4">
                                            <v-text-field v-model="CardRowData.CardID" label="CardID"></v-text-field>
                                        </v-col>

                                        <v-col cols="12" sm="6" md="4">
                                            <v-text-field v-model="CardRowData.card_name"
                                                label="card_name"></v-text-field>
                                        </v-col>
                                        <v-col cols="12" sm="6" md="4">
                                            <v-text-field v-model="CardRowData.card_type"
                                                label="card_type"></v-text-field>
                                        </v-col>
                                        <v-col cols="12" sm="6" md="4">
                                            <v-text-field v-model="CardRowData.card_character"
                                                label="card_character"></v-text-field>
                                        </v-col>
                                        <v-col cols="12" sm="6" md="4">
                                            <v-text-field v-model="CardRowData.card_condition"
                                                label="card_condition"></v-text-field>
                                        </v-col>
                                        <v-col cols="12" sm="6" md="4">
                                            <v-text-field v-model="CardRowData.other" label="other"></v-text-field>
                                        </v-col>
                                    </v-row>
                                </v-container>
                            </v-card-text>

                            <v-card-actions>
                                <v-spacer></v-spacer>
                                <v-btn color="blue-darken-1" variant="text" @click="closeAddDialog">
                                    取消
                                </v-btn>
                                <v-btn color="blue-darken-1" variant="text" @click="saveAddDialog">
                                    保存
                                </v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>

                    <!-- 删除数据询问框 -->
                    <v-dialog v-model="deleteAlertDialog" max-width="500px">
                        <v-card>
                            <v-card-title class="text-h5">你不准删</v-card-title>
                            <v-card-actions>
                                <v-spacer></v-spacer>
                                <v-btn color="blue-darken-1" variant="text" @click="closeDeleteAlertDialog">不删</v-btn>
                                <v-btn color="blue-darken-1" variant="text" @click="deleteCardConfirm">删！</v-btn>
                                <v-spacer></v-spacer>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>
                </v-toolbar>
            </template>

            <!-- 设置修改和删除按钮 -->
            <template v-slot:item.actions="{ item }">
                <v-icon size="small" class="me-2" @click="editCard(item.raw)">
                    mdi-pencil
                </v-icon>
                <v-icon size="small" @click="deleteCard(item.raw)">
                    mdi-delete
                </v-icon>
            </template>

            <!-- 设置已到货为绿色 -->
            <template v-slot:item.card_condition="{ item }">
                <v-chip :color="item.columns.card_condition === '已到货' ? 'green' : 'default'">
                    {{ item.columns.card_condition }}
                </v-chip>
            </template>
        </v-data-table>
    </div>
</template>