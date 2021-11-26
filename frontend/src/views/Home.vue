<template>
  <div class="home">
    <img @click="getMessage" alt="Vue logo" src="../assets/appicon.png" :style="{ height: '400px' }"/>
    <h2 @click="getNetStatus">测试network</h2>
    <HelloWorld :msg="message" :netStatus="netStatus"/>
    <h2 @click="storeCount">测试服务端存储数据</h2>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref} from "vue";
import HelloWorld from "@/components/HelloWorld.vue"; // @ is an alias to /src
import StoreFrontend from "@/utils/store"; // @ is an alias to /src

export default defineComponent({
  name: "HomePage",
  components: {
    HelloWorld,
  },
  setup() {
    const message = ref("Click the Icon");
    const netStatus = ref(false);

    const getMessage = () => {
      window.backend.basic().then(result => {
        message.value = result;
      });
    }

    const getNetStatus = () => {
      window.backend.NetWorkStatus().then(result => {
        netStatus.value = result
      });
    }

    const storeCount = () => {
      console.log(StoreFrontend(1))
    }
    return {message: message, netStatus: netStatus, getMessage: getMessage, getNetStatus: getNetStatus,storeCount:storeCount};
  },
});
</script>
