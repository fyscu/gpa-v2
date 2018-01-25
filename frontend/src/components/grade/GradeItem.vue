<template>
  <!-- 不及格成绩 -->
  <mu-card>
    <mu-card-header :title="title">
      <template slot="default">
        <!--<a v-if="isShowCSV == true" :download="`${title}.csv`" :href="exportCSV"> 导出CSV </a>-->
      </template>
    </mu-card-header>
    <mu-card-text>
      <div v-if="isHead" class="grade-head">
        <div>
          <div class="left">必修绩点：
            <span>{{grade.avg.required.gpa}}</span>
          </div>
          <div class="right">必修平均分：
            <span>{{grade.avg.required.grade}}</span>
          </div>
        </div>
        <div>
          <div class="left">全部绩点：
            <span>{{grade.avg.all.gpa}}</span>
          </div>
          <div class="right">全部平均分：
            <span>{{grade.avg.all.grade}}</span>
          </div>
        </div>
      </div>
      <mu-table class="t-1" @rowSelection="handleSelect" :enableSelectAll=true :multiSelectable=true>
        <colgroup>
          <col width="40">
          <col width="120">
          <col width="40">
          <col width="40">
          <col width="40">
          <col width="50">
        </colgroup>
        <mu-thead>
          <mu-tr class="table-head">
            <mu-th @click.native="sort('course_name')">课程名<i class="iconfont" :class="iconSort('course_name')" ></i></mu-th>
            <mu-th @click.native="sort('gradeCal')">分数<i class="iconfont" :class="iconSort('gradeCal')" ></i></mu-th>
            <mu-th @click.native="sort('gpa')">绩点<i class="iconfont" :class="iconSort('gpa')" ></i></mu-th>
            <mu-th @click.native="sort('credit')">学分<i class="iconfont" :class="iconSort('credit')" ></i></mu-th>
            <mu-th @click.native="sort('course_type')">属性<i class="iconfont" :class="iconSort('course_type')" ></i></mu-th>
          </mu-tr>
        </mu-thead>
        <mu-tbody>
          <mu-tr v-if="g.course_name" v-for="g in grade.grades" :selected="g.selected" :key="g.id">
            <input type="hidden" class="grade" :value="g.gradeCal">
            <input type="hidden" class="credit" :value="g.credit">
            <input type="hidden" class="gpa" :value="g.gpa">
            <mu-td>{{g.course_name}}</mu-td>
            <mu-td>{{g.grade}}</mu-td>
            <mu-td>{{g.gpa}}</mu-td>
            <mu-td>{{g.credit}}</mu-td>
            <mu-td>{{g.course_type}}</mu-td>
          </mu-tr>
        </mu-tbody>
      </mu-table>
    </mu-card-text>
  
  </mu-card>
</template>

<script>
import { card, cardHeader, cardTitle, cardText } from "muse-components/card"
import { table, tr, td, th, tbody, thead } from "muse-components/table"

// const json2csv = require('json2csv');

export default {
  props: {
    grade: {
      default: {
        avg: {
          required: {},
          all: {}
        },
        grades: {}
      }
    },
    isHead: {
      default: true
    },
    title: {
      required: true
    }
  },
  computed: {
    //导出csv
    // exportCSV() {
    //   let cvs = json2csv({
    //     data: this.grade.grades,
    //     fields: [
    //       {
    //         label: "课程号",
    //         value: "course_id",
    //         default: 'null'
    //       },
    //       {
    //         label: "课序号",
    //         value: "lesson_id",
    //         default: 'null'
    //       },
    //       {
    //         label: "课程名",
    //         value: "course_name",
    //         default: 'null'
    //       },
    //       {
    //         label: "属性",
    //         value: "course_type",
    //         default: 'null'
    //       },
    //       {
    //         label: "学分",
    //         value: "credit",
    //         default: 'null'
    //       },
    //       {
    //         label: "绩点",
    //         value: "gpa",
    //         default: 'null'
    //       },
    //       {
    //         label: "成绩",
    //         value: "grade",
    //         default: 'null'
    //       },
    //       {
    //         label: "学期",
    //         value: "term_name",
    //         default: 'null'
    //       },
    //     ]
    //   })
    //   let blob = new Blob([cvs], { type: 'text/csv' }); //new way  
    //   let csvUrl = URL.createObjectURL(blob);
    //   return csvUrl
    // },

    //是否展示csv
    // isShowCSV(){
    //   grade.grades.length > 0
    //   return false
    // }
  },

  data(){
    return {
      modes:{
        'gradeCal': 0,
        'gpa': 0,
        'credit': 0,
        'course_type': 0,
        'course_name': 0
      }
    }
  },
  
  methods: {
    //选择一行
    handleSelect(rowIndexs) {
      this.grade.grades = this.grade.grades.map((item, index) => {
        item.selected = rowIndexs.indexOf(index) !== -1
        return item
      })
    },

    iconSort(field){
      let icon = 'icon-sort'
      switch(this.modes[field]){
        case 0:
          icon = 'icon-sort';
          break;
        case 1:
          icon = 'icon-up';
          break;
        case -1:
          icon = 'icon-down';
          break;
      }
      return icon
    },

    sort(field){
      this.modes[field] = this.modes[field] == 0 ? 1 : -this.modes[field]
      this.sortGrades(field,this.modes[field])
      for(let x in this.modes){
        if(x != field)this.modes[x]=0;
      }
    },

    /**
     * @param field 字段名
     * @param mode 模式：1顺序/-1逆序 
     */
    sortGrades(field,mode){
      if (field == 'course_type' || field == 'course_name'){
        this.grade.grades.sort((a,b)=> mode == 1 ? a[field].localeCompare(b[field]) : b[field].localeCompare(a[field]));
      }else{
        this.grade.grades.sort((a,b)=> mode == 1 ? a[field] -b[field] :  b[field] -a[field]);
      }
    }
  },
  components: {
    "mu-table": table,
    "mu-thead": thead,
    "mu-th": th,
    "mu-tr": tr,
    "mu-td": td,
    "mu-tbody": tbody,
    "mu-card": card,
    "mu-card-header": cardHeader,
    "mu-card-title": cardTitle,
    "mu-card-text": cardText
  }
}

</script>
<style lang="less">
.table-head{
  th{
    cursor: pointer;
  }
  .iconfont{
    vertical-align: middle;
    font-size: 15px;
  }
}
</style>
