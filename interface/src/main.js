import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import Vuelidate from 'vuelidate'
import vuetify from './plugins/vuetify'

let api = axios.create({
  baseURL: process.env.VUE_APP_BACKEND_URL
})

Vue.config.productionTip = false
Vue.use(VueAxios, api)
Vue.use(Vuelidate)

new Vue({
  vuetify,
  render: h => h(App)
}).$mount('#app')
