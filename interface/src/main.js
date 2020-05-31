import Vue from 'vue'
import App from './App.vue'
import Buefy from 'buefy'
import axios from 'axios'
import VueAxios from 'vue-axios'
import Vuelidate from 'vuelidate'
import 'buefy/dist/buefy.css'

let api = axios.create({
  baseURL: process.env.VUE_APP_BACKEND_URL
})

Vue.config.productionTip = false
Vue.use(Buefy)
Vue.use(VueAxios, api)
Vue.use(Vuelidate)

new Vue({
  render: h => h(App)
}).$mount('#app')
