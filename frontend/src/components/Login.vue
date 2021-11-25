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
        <el-form-item label="账号：" prop="email">
          <el-input v-model="ruleForm.email" autocomplete="off" type="email"></el-input>
        </el-form-item>
          <el-form-item label="密码：" prop="pass">
          <el-input v-model="ruleForm.pass" autocomplete="off" type="password"></el-input>
        </el-form-item>
        <el-row :gutter="36">
          <el-form-item>
            <el-col :offset="6" :span="12">
              <el-button type="primary" @click="submitForm('ruleForm')">登录</el-button>
            </el-col>
            <el-link @click="register" type="primary" :underline="false">还没有账号，注册>></el-link>
          </el-form-item>
        </el-row>
      </el-form>
    </el-main>
    <el-footer>其他方式登录</el-footer>
  </el-container>
</template>
<style scoped>
@import "../assets/css/base.css";
</style>
<script>
import {defineComponent} from "vue";

export default defineComponent({
  name: "LoginPage",
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
    register: function () {
      this.$router.push({path:'/register'})
    }
  }
});
</script>
