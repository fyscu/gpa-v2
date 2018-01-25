<template>
  <div>
    <div v-if="!isLogin" id="login">
      <div class="main">
        <h2>四川大学绩点/平均分一键计算</h2>
        <mu-text-field label="学号" hintText="请输入学号" :errorText="valid.uid" fullWidth labelFloat @input="getUid"></mu-text-field>
        <br/>
        <mu-text-field label="密码" hintText="请输入教务处密码" :errorText="valid.password" fullWidth type="password" labelFloat @input="getPassword"></mu-text-field>
        <br/>
        <mu-raised-button label="登录" @click="login" class="demo-raised-button" fullWidth primary></mu-raised-button>
      </div>

      <div id="footer">
        <a target="_blank" href="http://lailin.xyz">&#xa9; Mohuishou</a>
        <a target="_blank" href="http://fyscu.com"> &#x261b; 飞扬研发实验室</a>
      </div>
    </div>

    <div v-if="isLogin" id="grade">
      <div>
        <mu-tabs class="grade-tabs" :value="activeTab" @change="switchTab">
          <mu-tab value="0" title="本学期成绩"></mu-tab>
          <mu-tab value="1" title="所有成绩" @click="getGradeAll"></mu-tab>
          <mu-tab value="2" title="不及格成绩" @click="getNotPass"></mu-tab>
        </mu-tabs>
      </div>

      <my-grade v-show="activeTab==0" :isHead=true :grade="grade" title="本学期成绩"></my-grade>

      <my-grade v-show="activeTab==1" :isHead=true :grade="grade" :title="grade.grades[0].term_name" v-for="(grade,index) in gradeAll" :key="index"></my-grade>

      <my-grade v-show="activeTab!=0" :isHead=false :grade="notPass[0]" title="尚不及格"></my-grade>

      <my-grade v-show="activeTab!=0" :isHead=false :grade="notPass[1]" title="曾不及格"></my-grade>

      <mu-bottom-nav class="bottom" :value="bottomData" @change="handleChange">
        <mu-bottom-nav-item value="cal" title="计算" @click.native="calculation" icon=":icon-CombinedShape" iconClass="iconfont"></mu-bottom-nav-item>
        <mu-bottom-nav-item value="required" title="必修" @click.native="chooseRequire" icon=":icon-zhuanyebixiuke" iconClass="iconfont"></mu-bottom-nav-item>
        <mu-bottom-nav-item value="all" title="全选" @click.native="chooseAll" icon=":icon-quanxuan" iconClass="iconfont"></mu-bottom-nav-item>
        <mu-bottom-nav-item value="clear" title="清空" @click.native="clear" icon=":icon-qingkong" iconClass="iconfont"></mu-bottom-nav-item>
        <!--<mu-bottom-nav-item value="clear" title="导出" @click.native="clear" icon=":icon-qingkong" iconClass="iconfont" />-->
      </mu-bottom-nav>
    </div>

    <mu-dialog :open="dialogCal" title="计算结果" @close="closeCal">
      <div id="cal-result">
        <div class="result">
          <p>您一共选择了 {{avg.count}} 门课程，总计 {{avg.credit}} 学分</p>
          <p>平均分: {{avg.grade}}</p>
          <p>绩点: {{avg.gpa}}</p>
        </div>
        <div class="attention">
          <p>注：</p>
          <p>学分绩点=∑(课程绩点*课程学分) / 课程总学分</p>
          <p>平均成绩=∑(课程成绩*课程学分) / 课程总学分</p>
        </div>
      </div>
      <mu-flat-button slot="actions" primary @click="closeCal" label="确定"></mu-flat-button>
    </mu-dialog>

    <mu-dialog :open="errorLog" title="提示" @close="closeCal">
      {{errorText}}
      <br />
      <p style="background: #eee;padding: 15px;">
        注：如果遇到网络问题，请先测试能否直接打开教务处网站->方案成绩->本学期成绩 如果没有问题请发送邮件到
        <a href="mailto:1@lailin.xyz">1@lailin.xyz</a>
      </p>
      <mu-flat-button slot="actions" primary @click="closeCal" label="确定"></mu-flat-button>
    </mu-dialog>

    <mu-dialog :open="helpShow" title="绩点计算规则更新说明" @close="closeCal">
      <p>2018年学期之后的成绩将采用新的绩点换算方式，对于之前取得的成绩采用原有的计算方式，绩点计算器将会自动判断年份来进行换算：<br />
        <a href="https://mp.weixin.qq.com/s/IYRRRsQ6NjyXtykSMs71ZQ">点击查看新换算方式的通知详情</a><br />
        <a href="http://jwc.scu.edu.cn/jwc/newsShow.action?news.id=1896">点击查看旧换算方式的通知详情</a>
      </p>
      <p>最新成绩绩点对照表：</p>
      <img style="max-width:100%;" src="../../assets/1.png" alt="成绩对照表">
      <p>说明：中文等级暂时按照对应分数段最高的计算，如果有最新的更加详细的对照表请发送邮件
        <a href="mailto:1@lailin.xyz">1@lailin.xyz</a>
      </p>

      <mu-flat-button slot="actions" primary @click="closeCal" label="确定"></mu-flat-button>
    </mu-dialog>

    <mu-dialog :open="isLoading" dialogClass="loading">
      <mu-circal :size="60" color="#fff"></mu-circal>
    </mu-dialog>

    <div class="loading" v-show="isLoading">

    </div>

  </div>
</template>
<script>
// 自定义组件
import "less/grade.less"
import gradeItem from "./GradeItem.vue"

//框架组件
import { tabs, tab } from "muse-components/tabs"
import { bottomNav, bottomNavItem } from "muse-components/bottomNav"
import dialog from "muse-components/dialog"
import flatButton from "muse-components/flatButton"
import raisedButton from "muse-components/raisedButton"
import textField from "muse-components/textField"
import circularProgress from "muse-components/circularProgress"

//配置文件
import config from 'config/grade'

import axios from 'axios'
import gradeObj from 'js/grade.js'
const qs = require('querystring')

//加密
const skey = "xrcytjhndsgfysavasdvahfvwj,ehbfkjb"
const encryptor = require('simple-encryptor')(skey);

let gradeTest = {
  avg: {
    required: {},
    all: {}
  },
  grades: {}
}
export default {
  components: {
    "my-grade": gradeItem,

    "mu-bottom-nav": bottomNav,
    "mu-bottom-nav-item": bottomNavItem,
    "mu-tabs": tabs,
    "mu-tab": tab,
    "mu-dialog": dialog,
    "mu-flat-button": flatButton,
    "mu-raised-button": raisedButton,
    "mu-text-field": textField,
    "mu-circal": circularProgress
  },
  data() {
    let http = axios.create({
      baseURL: config.domain,
      timeout: 5000,
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    });

    //请求拦截器：显示loding
    http.interceptors.request.use(conf => {
      this.isLoading = true
      conf.data = qs.stringify(conf.data)
      return conf
    }, err => {
      // 对请求错误做些什么
      return Promise.reject(error);
    });

    //响应拦截器
    http.interceptors.response.use(
      res => {
        this.isLoading = false
        return res
      },
      // err => {
      // this.isLoading = false
      // this.error(err.message)
      // }
    );

    return {
      http: http,
      valid: {
        uid: '',
        password: ''
      },
      helpShow: false,
      exportCSVShow: false,
      isLoading: false,
      isLogin: false,
      bottomData: 'cal',
      activeTab: "0",
      gradeAll: {},
      grade: gradeTest,
      dialogCal: false,
      errorText: "",
      errorLog: false,
      notPass: [
        gradeTest,
        gradeTest
      ],
      check: {
        gpa: 0,
        gpaAll: 0,
        notPass: 0
      },
      params: {
      },
      avg: {
        gpa: 0,
        credit: 0,
        grade: 0,
        count: 0
      }
    }
  },

  computed: {
    isWechat() {
      let ua = navigator.userAgent.toLowerCase();
      if (ua.match(/MicroMessenger/i) == "micromessenger") return true
      return false
    }
  },
  methods: {
    error(msg) {
      this.isLoading = false
      this.errorText = msg
      this.errorLog = true
    },
    getUid(val) {
      this.valid.uid = '';
      this.params.uid = val;
    },
    getPassword(val) {
      this.valid.password = '';
      this.params.password = val;
    },
    login() {
      if (!this.params.uid || isNaN(this.params.uid)) {
        this.valid.uid = "学号必填且必须是数字"
        return
      }
      if (!this.params.password) {
        this.valid.password = "密码必填"
        return
      }

      this.getGrade();

    },
    /**
     * 关闭计算窗口
     */
    closeCal() {
      this.dialogCal = false;
      this.errorLog = false;
      this.helpShow = false;
    },
    handleChange(val) {
      this.bottomData = val
    },
    switchTab(val) {
      this.activeTab = val;
    },
    /**
     * 获取所有成绩
     */
    getGradeAll() {
      if (this.check.gpaAll != 0) {
        return
      }
      this.http.post("/gpa/all", this.params).then(
        resp => {
          if (resp.data.status == 1) {
            this.isLogin = true
            let g = gradeObj.cal(resp.data.data);
            this.gradeAll = g.reverse();
            this.check.gpaAll = 1;
            this.getNotPass()
          }
        }
      )
    },
    getNotPass() {
      if (this.check.notPass != 0) {
        return
      }
      this.http.post("/gpa/not-pass", this.params).then(
        resp => {
          if (resp.data.status == 1) {
            this.notPass = gradeObj.cal(resp.data.data);
            this.check.notPass = 1;
          }
        }
      )

    },
    getGrade() {
      if (this.check.gpa != 0) {
        return
      }
      this.http.post("/gpa", this.params).then(
        resp => {
          if (resp.data.status == 1) {
            this.isLogin = true
            this.grade = gradeObj.cal([resp.data.data])[0];
            this.check.gpa = 1;
            if (this.isWechat && !localStorage.uid) {
              //在微信中打开，保存用户名和密码
              localStorage.uid = this.encrypt(this.params.uid)
              localStorage.password = this.encrypt(this.params.password)
            }
          } else {
            this.error("本学期成绩获取失败，尝试获取所有成绩，请稍后：" + resp.data.msg)
            setTimeout(() => {
              this.closeCal()
              this.switchTab(1)
              this.getGradeAll()
            }, 2000)
          }
        },
        resp => {
          this.error("本学期成绩获取失败，尝试获取所有成绩：" + resp.data.msg)
          setTimeout(() => {
            this.getGradeAll()
          }, 500)

        }
      )
    },

    //计算成绩/绩点
    calculation() {
      let all = {
        gpa: 0,
        credit: 0,
        grade: 0,
        count: 0
      }
      if (this.activeTab == 0) {
        all = this.calCount(this.grade, all)
      } else if (this.activeTab == 1) {
        for (let i = 0; i < this.gradeAll.length; i++) {
          all = this.calCount(this.gradeAll[i], all);
        }
        for (let i = 0; i < this.notPass.length; i++) {
          all = this.calCount(this.notPass[i], all);
        }
      } else if (this.activeTab == 2) {
        for (let i = 0; i < this.notPass.length; i++) {
          all = this.calCount(this.notPass[i], all);
        }
      }

      this.avg = all;
      this.avg.gpa = (all.gpa / all.credit).toFixed(3);
      this.avg.grade = (all.grade / all.credit).toFixed(3);

      this.dialogCal = true;
      console.log(this.avg);
    },

    //计算统计（计算前统计总分）
    calCount(arr, all) {
      arr.grades.map((item, index) => {
        if (item && item.selected) {
          all.gpa += item.gpa * item.credit;
          all.credit += item.credit;
          all.grade += item.gradeCal * item.credit;;
          all.count++;
        }
      })
      return all;
    },

    //选择所有课程（不包含不及格）
    chooseAll() {
      if (this.activeTab == 0) {
        this.grade.grades = this.grade.grades.map((item, index) => {
          item.selected = true
          return item
        })
      } else if (this.activeTab == 1) {
        for (let i = 0; i < this.gradeAll.length; i++) {
          this.gradeAll[i].grades = this.gradeAll[i].grades.map((item, index) => {
            item.selected = true
            return item
          })
        }
      }
    },

    //选择所有的必修课程（不包含不及格课程）
    chooseRequire() {
      if (this.activeTab == 0) {
        this.grade.grades = this.grade.grades.map((item, index) => {
          if (item && item.course_type && item.course_type == "必修") {
            item.selected = true
          } else {
            item.selected = false
          }
          return item
        })
      } else if (this.activeTab == 1) {
        for (let i = 0; i < this.gradeAll.length; i++) {
          this.gradeAll[i].grades = this.gradeAll[i].grades.map((item, index) => {
            if (item && item.course_type && item.course_type == "必修") {
              item.selected = true
            } else {
              item.selected = false
            }
            return item
          })
        }
      }

    },

    //清除所有已选项
    clear() {
      this.grade.grades = this.grade.grades.map((item, index) => {
        item.selected = false
        return item
      })

      if (this.check.notPass) {
        for (let i = 0; i < this.gradeAll.length; i++) {
          this.gradeAll[i].grades = this.gradeAll[i].grades.map((item, index) => {
            item.selected = false
            return item
          })
        }
        for (let i = 0; i < this.notPass.length; i++) {
          this.notPass[i].grades = this.notPass[i].grades.map((item, index) => {
            item.selected = false
            return item
          })
        }
      }

    },
    encrypt(text) {
      return encryptor.encrypt(text);
    },
    decrypt(hex) {
      return encryptor.decrypt(hex);
    }
  },
  mounted() {
    let ua = navigator.userAgent.toLowerCase();
    if (ua.match(/MicroMessenger/i) == "micromessenger") {
      //在微信中打开
      if (localStorage.uid) {
        this.params.uid = this.decrypt(localStorage.uid)
        if (localStorage.password) {
          this.params.password = this.decrypt(localStorage.password)
          this.getGrade()
        }
      }
    }
    if (localStorage.isHelp2) {
      this.helpShow = false
    } else {
      this.helpShow = true
      localStorage.isHelp2 = '0'
    }
  }
}

</script>

<style lang="less">
.iconfont {
  font-size: 17px;
}
</style>