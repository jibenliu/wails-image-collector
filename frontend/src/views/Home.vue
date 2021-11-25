<template>
  <div class="home">
    <img @click="getMessage" alt="Vue logo" src="../assets/appicon.png" :style="{ height: '400px' }"/>
    <h2 @click="getNetStatus">测试network</h2>
    <HelloWorld :msg="message" :netStatus="netStatus"/>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref} from "vue";
import HelloWorld from "@/components/HelloWorld.vue"; // @ is an alias to /src

interface Backend {
  basic(): Promise<string>;

  NetWorkStatus(): Promise<boolean>;
}

declare global {
  interface Window {
    backend: Backend;
  }
}

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
    return {message: message, netStatus: netStatus, getMessage: getMessage, getNetStatus: getNetStatus};
  },
});
</script>
