
let grade = {
};

/**
 * calType: 成绩计算类型，0：表示使用原有方式计算，1:表示采用新方式
 */
let calType = 0;


//适用于2018年查询的成绩统计
const GRADE_GPA_OLD = [
  [60, 0], //0-60[不包括60，绩点为0]
  [65, 1], //60-65[不包括61]
  [70, 1.7],
  [75, 2.2],
  [80, 2.7],
  [85, 3.2],
  [90, 3.6],
  [95, 3.8],
  [101, 4]
]

//分数绩点对照表
const GRADE_GPA = [
  [60, 0], //0-60[不包括60]
  [61, 1], //60-61[不包括61]
  [63, 1.3],
  [66, 1.7],
  [70, 2],
  [73, 2.3],
  [76, 2.7],
  [80, 3],
  [85, 3.3],
  [90, 3.7],
  [101, 4]
]

//等级成绩转换
const LEVLE_GRADE_OLD = {
  "优秀": 95,
  "良好": 86,
  "中等": 75,
  "通过": 60,
  "及格": 60,
  "未及格": 60,
  "不及格": 60
}

//等级成绩转换
const LEVLE_GRADE = {
  "优秀": 100,
  "良好": 84,
  "中等": 75,
  "合格": 69,
  "通过": 69,
  "及格": 69,
  "未及格": 60,
  "不及格": 60
}

/**
 * 绩点计算器
 * @param {Object} data 带计算的数据
 */
grade.cal = function (data) {
  let grade = [];
  for (let k = 0; k < data.length; k++) {
    if (!data[k]) {
      continue;
    }
    grade[k] = calTermGrade(data[k])
  }
  return grade;
}

/**
 * 设置对照类型
 * @param {String} term_name 
 */
let setType = (term_name) => {
  if (term_name == "") {
    calType = ((new Date()).getFullYear() > 2017) ? 1 : 0;
  } else {
    let y = term_name.match(/\d+\-(\d+)/);
    calType = (y.length > 1 && y[1] > 2017) ? 1 : 0;
  }
}

/**
 * 计算一个学期的平均成绩与绩点
 * @param {Array} grades 
 */
function calTermGrade(grades) {
  let sum = {
    all: {
      grade: 0,
      credit: 0,
      gpa: 0
    },
    required: {
      grade: 0,
      credit: 0,
      gpa: 0
    }
  };
  let avg = {
    all: {
      grade: 0,
      gpa: 0
    },
    required: {
      grade: 0,
      gpa: 0
    }

  };
  let results = grades.map((obj, i) => {
    if (obj.course_id) {
      i == 0 && setType(obj.term_name)
      obj.selected = false;
      obj.gradeCal = lv2grade(obj.grade);
      obj.gpa = grade2gpa(obj.gradeCal);
      obj.credit = obj.credit - 0;
      if (obj.grade==""){
        return obj
      }
      if (obj.course_type == "必修") {
        sum.required.grade += obj.gradeCal * obj.credit;
        sum.required.credit += obj.credit;
        sum.required.gpa += obj.gpa * obj.credit;
      }
      sum.all.grade += obj.gradeCal * obj.credit;
      sum.all.credit += obj.credit;
      sum.all.gpa += obj.gpa * obj.credit;
    }
    return obj
  })
  avg.all.gpa = (sum.all.gpa / sum.all.credit).toFixed(3);
  avg.all.grade = (sum.all.grade / sum.all.credit).toFixed(3);

  avg.required.gpa = (sum.required.gpa / sum.required.credit).toFixed(3);
  avg.required.grade = (sum.required.grade / sum.required.credit).toFixed(3);
  return {
    avg: avg,
    grades: results
  }
}

/**
 * 分数转绩点
 * @param {number} grade 分数
 */
function grade2gpa(grade) {
  let gpa = 0;
  let prev = 0;
  let gradeGPA = GRADE_GPA_OLD;
  if (calType == 1) {
    gradeGPA = GRADE_GPA
  }
  for (let i = 0; i < gradeGPA.length; i++) {
    if (grade >= prev && grade < gradeGPA[i][0]) {
      gpa = gradeGPA[i][1]
    }
    prev = gradeGPA[i][0]
  }
  return gpa;
}

/**
 * 去除所有的空格
 * @param {String} str 字符串
 */
function trimStr(str) {
  return str.replace(/(^\s*)|(\s*$)/g, "");
}

/**
 * 等级转换为分数
 * @param {String} g 等级
 */
function lv2grade(g) {
  g = trimStr(g);
  if (!isNaN(g)) {
    return g;
  }
  if (calType == 0) {
    return LEVLE_GRADE_OLD[g]
  } else {
    return LEVLE_GRADE[g]
  }
}

module.exports = grade;
