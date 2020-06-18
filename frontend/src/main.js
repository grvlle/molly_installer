import 'core-js/stable';
import 'regenerator-runtime/runtime';
import Vue from 'vue';
import App from './App.vue';
import router from './router'
import {store} from './store/store'
import '@fortawesome/fontawesome-free/css/all.css'
import '@fortawesome/fontawesome-free/js/all.js'


Vue.config.productionTip = false;
Vue.config.devtools = true;

import * as Wails from '@wailsapp/runtime';

Wails.Init(() => {
	new Vue({
		router,
		store: store,
		render: h => h(App)
	}).$mount('#app');
});

router.replace('/')
