import './assets/main.css'

import { createApp } from 'vue'
import App from './App_Square.vue'
import router from './router'

const app = createApp(App)

app.use(router)

app.mount('#app')
