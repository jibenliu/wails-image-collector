<template>
  <div>
    <div id="container"></div>
  </div>
</template>

<script>
import * as THREE from "three";
import {OrbitControls} from "three/examples/jsm/controls/OrbitControls.js";
import {Water} from 'three/examples/jsm/objects/Water.js';
import {Sky} from "three/examples/jsm/objects/Sky.js";
import {defineComponent} from "vue";

export default defineComponent({
  name: "SkyPage",
  data() {
    return {
      renderer: "",
      scene: "",
      camera: "",
      sphere: "",
      waterMesh: "",
    }
  },
  methods: {
    onload: function () {
      let container = document.getElementById("container");//获取container
      this.scene = new THREE.Scene();//创建场景
      this.camera = new THREE.PerspectiveCamera(45, window.innerWidth / window.innerHeight, 1, 10000);
      this.camera.position.set(30, 20, 100);

      this.renderer = new THREE.WebGLRenderer({antialias: true});
      this.renderer.setPixelRatio(window.devicePixelRatio);
      this.renderer.setSize(window.innerWidth, window.innerHeight);
      container.appendChild(this.renderer.domElement);

      let directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
      this.scene.add(directionalLight);

      //创建水平面
      let waterPlane = new THREE.PlaneBufferGeometry(10000, 10000);
      this.waterMesh = new Water(waterPlane, {
        textureWidth: 512,//画布宽度
        textureHeight: 512,//画布高度
        waterNormals: new THREE.TextureLoader().load("./img/waternormals.jpg", function (texture) {
          texture.wrapS = texture.wrapT = THREE.RepeatWrapping;//法向量贴图
        }),
        alpha: 1.0,//透明度
        sunDirection: directionalLight.position.clone().normalize(),
        sunColor: 0xffffff,//太阳的颜色
        waterColor: 0x001e0f,//水的颜色
        distortionScale: 3.7,//物体倒影的分散度
        fog: this.scene.fog !== undefined,
      });
      this.waterMesh.rotation.x = -Math.PI / 2;
      this.scene.add(this.waterMesh);

      //创建天空盒子

      // Skybox
      let sky = new Sky();
      let uniforms = sky.material.uniforms;
      uniforms['turbidity'].value = 10;//内置变量
      uniforms['rayleigh'].value = 2;//视觉效果就是傍晚晚霞的红光的深度
      uniforms['luminance'].value = 1;//视觉效果整体提亮或变暗0-1
      uniforms['mieCoefficient'].value = 0.005;
      uniforms['mieDirectionalG'].value = 0.8;

      let parameters = {
        distance: 400,
        inclination: 0.49,//倾向
        azimuth: 0.205//方位
      };

      let cubeCamera = new THREE.CubeCamera(0.1, 1, 512);//创建反光效果
      cubeCamera.renderTarget.texture.generateMipmaps = true;
      cubeCamera.renderTarget.texture.minFilter = THREE.LinearMipmapLinearFilter;
      this.scene.background = cubeCamera.renderTarget;

      function updateSun(waterMesh, renderer) {
        let theta = Math.PI * (parameters.inclination - 0.5);//-0.01*Math.PI -0.0314
        let phi = 2 * Math.PI * (parameters.azimuth - 0.5);//-0.59*Math.PI=-1.8535
        directionalLight.position.x = parameters.distance * Math.cos(phi);//399.79
        directionalLight.position.y = parameters.distance * Math.sin(phi) * Math.sin(theta);//-0.14
        directionalLight.position.z = parameters.distance * Math.sin(phi) * Math.cos(theta);//-0.323
        sky.material.uniforms['sunPosition'].value = directionalLight.position.copy(directionalLight.position);//设置太阳的位置
        waterMesh.material.uniforms['sunDirection'].value.copy(directionalLight.position).normalize();//设置太阳的光照方向，并进行归一化(化为单位值)
        cubeCamera.update(renderer, sky);
      }

      updateSun(this.waterMesh, this.renderer);

      //创建多边形物体
      let geometry = new THREE.IcosahedronBufferGeometry(20, 1);
      let count = geometry.attributes.position.count;
      let colors = [];
      let color = new THREE.Color();
      for (let i = 0; i < count; i += 3) {//设置颜色
        color.setHex(Math.random() * 0xffffff);//255
        colors.push(color.r, color.g, color.b);
        colors.push(color.r, color.g, color.b);
        colors.push(color.r, color.g, color.b);
      }
      //console.log(colors);
      geometry.setAttribute('color', new THREE.Float32BufferAttribute(colors, 3));//向顶点传入顶点颜色
      let material = new THREE.MeshStandardMaterial({
        vertexColors: THREE.VertexColors,//使用顶点颜色进行着色
        roughness: 0.0,
        flatShading: true,
        envMap: cubeCamera.renderTarget.texture,//设置环境贴图
        side: THREE.DoubleSide
      });
      this.sphere = new THREE.Mesh(geometry, material);
      this.scene.add(this.sphere);

      let sphereGeometry = new THREE.SphereBufferGeometry(8, 100, 100);
      let sphereMaterial = new THREE.MeshBasicMaterial({envMap: cubeCamera.renderTarget.texture});
      let sphereMesh = new THREE.Mesh(sphereGeometry, sphereMaterial);
      sphereMesh.position.set(30, 10, 0);
      this.scene.add(sphereMesh);

      let contorl = new OrbitControls(this.camera, this.renderer.domElement);//添加鼠标滚动缩放，旋转对象
      contorl.minDistance = 40;//最大最小相机移动距离(景深相机)
      contorl.maxDistance = 200;
      contorl.maxPolarAngle = Math.PI * 0.495; //最大仰视角和俯视角
      contorl.minPolarAngle = 0;
      contorl.update();
      window.addEventListener('resize', this.onResize, false);//浏览器大小改变监听
    },
    onResize: function () {
      this.camera.aspect = window.innerWidth / window.innerHeight;
      this.camera.updateProjectionMatrix();
      this.renderer.setSize(window.innerWidth, window.innerHeight);
    },
    run: function () {
      requestAnimationFrame(this.run);
      let time = performance.now() * 0.001;
      this.sphere.position.y = Math.sin(time) * 20 + 5;//物体上下移动
      this.sphere.rotation.x = time * 0.5;
      this.sphere.rotation.z = time * 0.51;
      this.waterMesh.material.uniforms['time'].value += 1.0 / 60.0;//创建水面微波
      this.renderer.render(this.scene, this.camera);
    },
    mounted() {
      this.onload();
      this.run();
    },
  }
})
</script>

<style scoped>
</style>
