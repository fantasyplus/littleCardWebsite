<script setup>
import { ref, toRaw } from "vue"
import request from "./utils/request.js"
// data
let searchInput = $ref("")
let tableData = $ref([{
  id: "1",
  name: 'tom1',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
},
{
  id: "2",
  name: 'Tom2',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
},
{
  id: "3",
  name: 'Tom3',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
},
{
  id: "4",
  name: 'Tom4',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
},
{
  id: "5",
  name: 'Tom5',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
}
])
let tableDataCopy = Object.assign(tableData)

let multipleTableSelection = $ref([])
let addDialogLabelPosition = $ref("left")
let addDialogVisible = $ref(false)
let dialogType = $ref("add")
let tableRowData = $ref({
  id: "",
  name: "",
  state: "",
  city: "",
  address: "",
})

let totalPage = $ref(5)
let curPage=$ref(1)
// methods


//request page data from server
const getTableData = async (cur = 1) => {
  let res = await request.get("/list", {
    pageSize: 5,
    pageNum: cur
  })
  tableData = res.data.data
  totalPage = Math.ceil(res.data.totalitems / res.data.pageSize)
  console.log("init table:",res)
}
//init table data
getTableData()

const handleChangePage = (val) => {
  getTableData(val)
}

// delete row from server
const handleDelRow = async (row) => {
  // let id = row.id
  // let index = tableData.findIndex((item) => item.id === id)
  // tableData.splice(index, 1)

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

const handleSelectionChange = (val) => {
  multipleTableSelection = val
}

// add
const handleAdd = () => {
  addDialogVisible = true
  dialogType = "add"
  tableRowData = {}
}

//edit
const handleEditRow = (row) => {
  addDialogVisible = true
  dialogType = "edit"
  tableRowData = { ...row }
}

// confirm add or edit data
const handleDialogConfirm = async () => {
  addDialogVisible = false
  if (dialogType === "add") {
    let res = await request.post(`/add`, {
      ...tableRowData
    })
    if (res.code === 200) {
      await getTableData()
    }
  }
  else if (dialogType === "edit") {
    let res = await request.put(`/update/${tableRowData.ID}`, {
      ...tableRowData
    })
    if (res.code === 200) {
      await getTableData()
    }
  }
}

//search
const handleSearch = async () => {
  // pure front search
  // if (searchInput === "") {
  //   tableData = tableDataCopy
  //   return
  // }
  // let searchResult = tableDataCopy.filter((item) => {
  //   return item.name.match(searchInput)
  // })
  // tableData = searchResult

  // search from server
  if (searchInput.length > 0) {
    let res = await request.get(`/list/${searchInput}`)
    tableData = res.data
  }
  else{
    await getTableData(curPage)
  }
}

</script>

<template>
  <div class="table-box">
    <!-- Title -->
    <div class="title">
      <h2>the most simple crud demo</h2>
    </div>
    <!-- query -->
    <div class="query-box">
      <el-input class="search-input" @input="handleSearch" v-model="searchInput"
        placeholder="Please input name to search" />
      <div class="btn-list">
        <el-button @click="handleAdd" type="primary">Add</el-button>
        <el-button v-if="multipleTableSelection.length > 0" @click="handleDelMultiRows" type="danger">Delete</el-button>
      </div>
    </div>
    <!-- table -->
    <el-table border ref="multipleTableRef" :data="tableData" style="width: 100%"
      @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="Name" width="120" />
      <el-table-column prop="state" label="State" width="120" />
      <el-table-column prop="city" label="City" width="120" />
      <el-table-column prop="address" label="Address" width="300" />
      <el-table-column fixed="right" label="Operations" width="120">
        <template #default="scope">
          <el-button link type="primary" size="small" style="color: #F56C6C
;" @click="handleDelRow(scope.row)">Delete</el-button>
          <el-button link type="primary" size="small" @click="handleEditRow(scope.row)">Edit</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination @current-change="handleChangePage" class="table-pagination" layout="prev, pager, next"
      :page-count="totalPage" :hide-on-single-page="true" v-model:current-page="curPage"/>

  </div>
  <!-- Form -->
  <el-dialog @keyup.enter="handleDialogConfirm" v-model="addDialogVisible" :title="dialogType === 'add' ? 'add' : 'edit'"
    draggable>
    <el-form v-model="tableRowData" :label-position="addDialogLabelPosition" :label-width="120">
      <el-form-item label="name">
        <el-input v-model="tableRowData.name" autocomplete="off" />
      </el-form-item>
      <el-form-item label="state">
        <el-input v-model="tableRowData.state" autocomplete="off" />
      </el-form-item>
      <el-form-item label="city">
        <el-input v-model="tableRowData.city" autocomplete="off" />
      </el-form-item>
      <el-form-item label="address">
        <el-input v-model="tableRowData.address" autocomplete="off" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="addDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleDialogConfirm">
          Confirm
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style scoped>
.table-box {
  width: 800px;
  /* margin: 200px; */
  position: absolute;
  top: 50%;
  left: 50%;
  /* 通过使用transform属性和translate函数，
  将元素向左和向上移动自身宽度和高度的50%。
  这样，元素的中心点将准确地位于容器的中心位置。 */
  transform: translate(-50%, -50%);
}

.title {
  text-align: center;
}

.query-box {
  display: flex;
  /* 这个属性用于设置子元素在主轴（水平方向）上的对齐方式。
  space-between表示子元素会在主轴上均匀分布，让第一个子元素在最左边，最后一个子元素在最右边，
  中间的子元素则在它们之间均匀分布，形成空间间隔。 */
  justify-content: space-between;
  margin-bottom: 20px;
}

.table-pagination {
  display: flex;
  justify-content: center;
  margin-top: 10px;
}

.search-input {
  width: 300px;
}

.el-form-item {
  text-align: center;
}
</style>
