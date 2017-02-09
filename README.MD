飞扬俱乐部GPA绩点计算器API

> gpa_v2

## server

```
端口：6627
```

## 获取本学期成绩

```json
path: /gpa
method: post
param: {
    uid: //学号
    password: //密码
}
return: {
  "status": 1,
  "msg": "本学期成绩获取成功",
  "data": [
    {
      "course_id": "998013000",
      "lesson_id": "24",
      "course_name": "形势与政策-5",
      "credit": "0",
      "course_type": "必修",
      "grade": "94",
      "term": 1, //学期标识，只在所有成绩有效
      "term_name": "" 
    },
    ......
  ]
}
```

## 获取所有成绩
```json
path: /gpa/all
method: post
param: {
    uid: //学号
    password: //密码
}
return: {
  "status": 1,
  "msg": "所有成绩获取成功",
  "data": [
    [
      {
        "course_id": "105366020",
        "lesson_id": "273",
        "course_name": "大学英语（综合）-1",
        "credit": "2",
        "course_type": "必修",
        "grade": "65.0",
        "term": 0,
        "term_name": "2014-2015学年秋(两学期)"
      },
    ],
  ],
}
```

## 获取不及格成绩
```json
path: /gpa/not-pass
method: post
param: {
    uid: //学号
    password: //密码
}
return:{
  "status": 1,
  "msg": "不及格成绩获取成功",
  "data": [
    [ //第一个数组 尚不及格成绩
      {
        "course_id": "201051020",
        "lesson_id": "01",
        "course_name": "数学建模",
        "credit": "2",
        "course_type": "选修",
        "grade": "不及格",
        "term": 0,
        "term_name": ""
      }
    ],
    [ 
        //第二个数组 曾不及格成绩
    ]
  ]
}
```