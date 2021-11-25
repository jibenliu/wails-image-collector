<template>
  <el-container>
    <el-header>
      <el-row>
        <el-col :span="24">
          <div id="login-title">APP登录</div>
        </el-col>
        <el-col :span="24"><img id="login-logo" alt="Wails logo" src="../assets/images/logo.png"></el-col>
      </el-row>
    </el-header>
    <el-main>
      <h1>{{ message }}</h1>
      <el-form ref="ruleForm" :model="ruleForm" :rules="rules" class="demo-ruleForm" label-width="100px" status-icon>
        <el-form-item label="邮箱地址：" prop="email">
          <el-input v-model="ruleForm.email" autocomplete="off" type="email"></el-input>
        </el-form-item>
        <el-form-item label="密码：" prop="pass">
          <el-input v-model="ruleForm.pass" autocomplete="off" type="password"></el-input>
        </el-form-item>
        <el-form-item label="确认密码：" prop="checkPass">
          <el-input v-model="ruleForm.checkPass" autocomplete="off" type="password"></el-input>
        </el-form-item>
        <el-row :gutter="36">
          <el-form-item>
            <el-col :offset="6" :span="12">
              <el-button type="primary" @click="submitForm('ruleForm')">登录</el-button>
              <el-button type="primary" @click="getMessage">注册</el-button>
              <el-button type="primary" @click="getMyStruct">绑定数据</el-button>
              <el-button type="primary" @click="openFile">打开文件</el-button>
            </el-col>
          </el-form-item>
        </el-row>
      </el-form>
    </el-main>
    <el-footer>页尾</el-footer>
  </el-container>
</template>
<style scoped>
@import "../assets/css/base.css";
</style>
<script>
import {defineComponent} from "vue";

export default defineComponent({
  name: "RegisterPage",
  data() {
    let validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'));
      } else {
        if (this.ruleForm.checkPass !== '') {
          this.$refs.ruleForm.validateField('checkPass');
        }
        callback();
      }
    };
    let validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'));
      } else if (value !== this.ruleForm.pass) {
        callback(new Error('两次输入密码不一致!'));
      } else {
        callback();
      }
    };
    return {
      message: " ",
      myStruct: {},
      ruleForm: {
        pass: '',
        checkPass: '',
      },
      rules: {
        pass: [
          {validator: validatePass, trigger: 'blur'}
        ],
        checkPass: [
          {validator: validatePass2, trigger: 'blur'}
        ]
      }
    };
  },
  methods: {
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          alert('submit!');
        } else {
          console.log('error submit!!');
          return false;
        }
      });
    },
    getMessage: function () {
      window.backend.basic().then(result => {
        this.message = result;
      });
    },
    getMyStruct: function () {
      window.backend.MyStruct.Hello("aaa").then(result => {
        this.myStruct = result;
        console.log(result)
      });
    },
    openFile: function () {
      window.backend.MyStruct.OpenFile().then(result => {
        alert(result)
      })
    }
  }
})
</script>
